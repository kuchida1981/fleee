package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kosuke/fleee/internal/handler"
	"github.com/kosuke/fleee/internal/importer"
	"github.com/kosuke/fleee/internal/store"
	"github.com/kosuke/fleee/internal/testutil"
)

func setupTestHandler(t *testing.T) (*handler.AccountHandler, *chi.Mux) {
	t.Helper()
	db := testutil.NewTestDB(t)
	accountStore := store.NewAccountStore(db)
	accountImporter := importer.NewAccountImporter(accountStore)
	h := handler.NewAccountHandler(accountStore, accountImporter)
	r := chi.NewRouter()
	r.Mount("/api/accounts", h.Routes())
	return h, r
}

func TestAccountHandler_List(t *testing.T) {
	_, r := setupTestHandler(t)

	t.Run("empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/accounts", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var accounts []json.RawMessage
		if err := json.NewDecoder(w.Body).Decode(&accounts); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(accounts) != 0 {
			t.Errorf("expected empty array, got %d items", len(accounts))
		}
	})

	t.Run("with accounts", func(t *testing.T) {
		body := `{"name":"Cash","account_type":"asset","display_order":1}`
		req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("setup: expected 201, got %d: %s", w.Code, w.Body.String())
		}

		req = httptest.NewRequest(http.MethodGet, "/api/accounts", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var accounts []json.RawMessage
		if err := json.NewDecoder(w.Body).Decode(&accounts); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(accounts) != 1 {
			t.Errorf("expected 1 account, got %d", len(accounts))
		}
	})
}

func TestAccountHandler_Create(t *testing.T) {
	_, r := setupTestHandler(t)

	t.Run("success", func(t *testing.T) {
		body := `{"name":"Revenue","account_type":"revenue","display_order":10}`
		req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
		var resp map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if resp["name"] != "Revenue" {
			t.Errorf("expected name Revenue, got %v", resp["name"])
		}
	})

	t.Run("validation error - empty name", func(t *testing.T) {
		body := `{"name":"","account_type":"asset","display_order":1}`
		req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("validation error - invalid type", func(t *testing.T) {
		body := `{"name":"Test","account_type":"invalid","display_order":1}`
		req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		body := `{"name":"Revenue","account_type":"revenue","display_order":20}`
		req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("expected 409, got %d", w.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBufferString("{bad"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestAccountHandler_Get(t *testing.T) {
	_, r := setupTestHandler(t)

	// Create an account first
	body := `{"name":"Bank","account_type":"asset","display_order":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("setup: failed to create account: %s", w.Body.String())
	}
	var created map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("setup: failed to decode response: %v", err)
	}
	id := int(created["id"].(float64))

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/accounts/%d", id), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var resp map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if resp["name"] != "Bank" {
			t.Errorf("expected name Bank, got %v", resp["name"])
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/accounts/99999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/accounts/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestAccountHandler_Update(t *testing.T) {
	_, r := setupTestHandler(t)

	// Create accounts and capture IDs
	var ids []int
	for _, b := range []string{
		`{"name":"Cash","account_type":"asset","display_order":1}`,
		`{"name":"Bank","account_type":"asset","display_order":2}`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBufferString(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("setup: failed to create account: %s", w.Body.String())
		}
		var created map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
			t.Fatalf("setup: failed to decode response: %v", err)
		}
		ids = append(ids, int(created["id"].(float64)))
	}

	t.Run("success", func(t *testing.T) {
		body := `{"name":"Petty Cash","account_type":"asset","display_order":1}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/accounts/%d", ids[0]), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("not found", func(t *testing.T) {
		body := `{"name":"Test","account_type":"asset","display_order":1}`
		req := httptest.NewRequest(http.MethodPut, "/api/accounts/99999", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		body := `{"name":"Bank","account_type":"asset","display_order":1}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/accounts/%d", ids[0]), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("expected 409, got %d", w.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/accounts/%d", ids[0]), bytes.NewBufferString("{bad"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		body := `{"name":"","account_type":"asset","display_order":1}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/accounts/%d", ids[0]), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		body := `{"name":"Test","account_type":"invalid","display_order":1}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/accounts/%d", ids[0]), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestAccountHandler_Delete(t *testing.T) {
	_, r := setupTestHandler(t)

	// Create an account
	body := `{"name":"Temp","account_type":"expense","display_order":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("setup: failed to create account: %s", w.Body.String())
	}
	var created map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("setup: failed to decode response: %v", err)
	}
	deleteID := int(created["id"].(float64))

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/accounts/%d", deleteID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/accounts/99999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/accounts/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestAccountHandler_Import(t *testing.T) {
	_, r := setupTestHandler(t)

	t.Run("success tsv", func(t *testing.T) {
		tsvContent := "科目名\t科目貸借タイプ\t出力順番\t精算種別\n普通預金\t借方\t0\t貸借対照表\n売上\t貸方\t1\t損益計算書\n"

		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		part, err := writer.CreateFormFile("file", "accounts.tsv")
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		if _, err := io.WriteString(part, tsvContent); err != nil {
			t.Fatalf("failed to write tsv content: %v", err)
		}
		if err := writer.Close(); err != nil {
			t.Fatalf("failed to close writer: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/accounts/import", &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var resp map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if resp["success"] != float64(2) {
			t.Errorf("expected 2 successes, got %v", resp["success"])
		}
	})

	t.Run("no file", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/accounts/import", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}
