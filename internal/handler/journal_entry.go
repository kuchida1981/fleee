package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kosuke/fleee/internal/model"
	"github.com/kosuke/fleee/internal/store"
)

// JournalEntryHandler coordinates endpoints for journal entry operations
type JournalEntryHandler struct {
	store *store.JournalEntryStore
}

// NewJournalEntryHandler creates a new JournalEntryHandler instance
func NewJournalEntryHandler(store *store.JournalEntryStore) *JournalEntryHandler {
	return &JournalEntryHandler{
		store: store,
	}
}

// Routes configures endpoints and returns a sub-router for journal entries
func (h *JournalEntryHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

type journalEntryRequest struct {
	Date            string               `json:"date"`
	Description     string               `json:"description"`
	ReceiptRequired bool                 `json:"receipt_required"`
	Memo            string               `json:"memo"`
	Lines           []journalLineRequest `json:"lines"`
}

type journalLineRequest struct {
	AccountID    int64 `json:"account_id"`
	DebitAmount  int64 `json:"debit_amount"`
	CreditAmount int64 `json:"credit_amount"`
}

func parseJournalEntryRequest(r *http.Request) (*journalEntryRequest, error) {
	var req journalEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.Date = strings.TrimSpace(req.Date)
	req.Description = strings.TrimSpace(req.Description)
	return &req, nil
}

func validateJournalEntryRequest(req *journalEntryRequest, w http.ResponseWriter) bool {
	if req.Date == "" {
		respondWithError(w, http.StatusBadRequest, "date is required")
		return false
	}
	if req.Description == "" {
		respondWithError(w, http.StatusBadRequest, "description is required")
		return false
	}
	return true
}

func (req *journalEntryRequest) toModel(id int64) *model.JournalEntry {
	lines := make([]model.JournalLine, len(req.Lines))
	for i, l := range req.Lines {
		lines[i] = model.JournalLine{
			AccountID:    l.AccountID,
			DebitAmount:  l.DebitAmount,
			CreditAmount: l.CreditAmount,
		}
	}
	return &model.JournalEntry{
		ID:              id,
		Date:            req.Date,
		Description:     req.Description,
		ReceiptRequired: req.ReceiptRequired,
		Memo:            req.Memo,
		Lines:           lines,
	}
}

func handleStoreError(w http.ResponseWriter, err error) {
	if errors.Is(err, store.ErrUnbalanced) {
		respondWithError(w, http.StatusBadRequest, "journal entry is not balanced")
		return
	}
	if errors.Is(err, store.ErrInsufficientLines) {
		respondWithError(w, http.StatusBadRequest, "journal entry must have at least 2 lines")
		return
	}
	if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		respondWithError(w, http.StatusBadRequest, "invalid account ID")
		return
	}
	respondWithError(w, http.StatusInternalServerError, err.Error())
}

// List handles GET /
func (h *JournalEntryHandler) List(w http.ResponseWriter, r *http.Request) {
	entries, err := h.store.ListAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if entries == nil {
		entries = []*model.JournalEntry{}
	}
	respondWithJSON(w, http.StatusOK, entries)
}

// Create handles POST /
func (h *JournalEntryHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := parseJournalEntryRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !validateJournalEntryRequest(req, w) {
		return
	}

	entry := req.toModel(0)
	if err := h.store.Create(r.Context(), entry); err != nil {
		handleStoreError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, entry)
}

// Get handles GET /{id}
func (h *JournalEntryHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid journal entry ID")
		return
	}

	entry, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "journal entry not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, entry)
}

// Update handles PUT /{id}
func (h *JournalEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid journal entry ID")
		return
	}

	req, err := parseJournalEntryRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !validateJournalEntryRequest(req, w) {
		return
	}

	entry := req.toModel(id)
	err = h.store.Update(r.Context(), entry)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "journal entry not found")
			return
		}
		handleStoreError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, entry)
}

// Delete handles DELETE /{id}
func (h *JournalEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid journal entry ID")
		return
	}

	err = h.store.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "journal entry not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
