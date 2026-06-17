package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kosuke/fleee/internal/importer"
	"github.com/kosuke/fleee/internal/model"
	"github.com/kosuke/fleee/internal/store"
)

// AccountHandler coordinates endpoints for account operations
type AccountHandler struct {
	store    *store.AccountStore
	importer *importer.AccountImporter
}

// NewAccountHandler creates a new AccountHandler instance
func NewAccountHandler(store *store.AccountStore, importer *importer.AccountImporter) *AccountHandler {
	return &AccountHandler{
		store:    store,
		importer: importer,
	}
}

// Routes configures endpoints and returns a sub-router for accounts
func (h *AccountHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	r.Post("/import", h.Import)
	return r
}

type errorResponse struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func isValidAccountType(t model.AccountType) bool {
	switch t {
	case model.AccountTypeAsset, model.AccountTypeLiability, model.AccountTypeEquity, model.AccountTypeRevenue, model.AccountTypeExpense:
		return true
	}
	return false
}

// List handles GET /api/accounts
func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.store.ListAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if accounts == nil {
		accounts = []*model.Account{}
	}
	respondWithJSON(w, http.StatusOK, accounts)
}

type createRequest struct {
	Name         string            `json:"name"`
	AccountType  model.AccountType `json:"account_type"`
	DisplayOrder int               `json:"display_order"`
}

// Create handles POST /api/accounts
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "account name is required")
		return
	}
	if !isValidAccountType(req.AccountType) {
		respondWithError(w, http.StatusBadRequest, "invalid account type")
		return
	}

	acc := &model.Account{
		Name:         req.Name,
		AccountType:  req.AccountType,
		DisplayOrder: req.DisplayOrder,
	}

	err := h.store.Create(r.Context(), acc)
	if err != nil {
		if errors.Is(err, store.ErrDuplicateName) {
			respondWithError(w, http.StatusConflict, "account name already exists")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, acc)
}

// Get handles GET /api/accounts/:id
func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid account ID")
		return
	}

	acc, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "account not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, acc)
}

type updateRequest struct {
	Name         string            `json:"name"`
	AccountType  model.AccountType `json:"account_type"`
	DisplayOrder int               `json:"display_order"`
}

// Update handles PUT /api/accounts/:id
func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid account ID")
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "account name is required")
		return
	}
	if !isValidAccountType(req.AccountType) {
		respondWithError(w, http.StatusBadRequest, "invalid account type")
		return
	}

	acc := &model.Account{
		ID:           id,
		Name:         req.Name,
		AccountType:  req.AccountType,
		DisplayOrder: req.DisplayOrder,
	}

	err = h.store.Update(r.Context(), acc)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "account not found")
			return
		}
		if errors.Is(err, store.ErrDuplicateName) {
			respondWithError(w, http.StatusConflict, "account name already exists")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, acc)
}

// Delete handles DELETE /api/accounts/:id
func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid account ID")
		return
	}

	err = h.store.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "account not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Import handles POST /api/accounts/import
func (h *AccountHandler) Import(w http.ResponseWriter, r *http.Request) {
	// Restrict size to 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "file parameter is required")
		return
	}
	defer file.Close()

	isTSV := strings.HasSuffix(strings.ToLower(header.Filename), ".tsv")

	result, err := h.importer.Import(r.Context(), file, isTSV)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, result)
}
