package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kosuke/fleee/internal/handler"
	"github.com/kosuke/fleee/internal/importer"
	"github.com/kosuke/fleee/internal/store"
	"github.com/kosuke/fleee/internal/testutil"
)

func setupJournalEntryTestHandler(t *testing.T) *chi.Mux {
	t.Helper()
	db := testutil.NewTestDB(t)
	accountStore := store.NewAccountStore(db)
	accountImporter := importer.NewAccountImporter(accountStore)
	accountHandler := handler.NewAccountHandler(accountStore, accountImporter)
	journalEntryStore := store.NewJournalEntryStore(db)
	journalEntryHandler := handler.NewJournalEntryHandler(journalEntryStore)
	r := chi.NewRouter()
	r.Mount("/api/accounts", accountHandler.Routes())
	r.Mount("/api/journal-entries", journalEntryHandler.Routes())
	return r
}

func createTestAccountsViaAPI(t *testing.T, r *chi.Mux) (expenseID, assetID int) {
	t.Helper()
	for _, body := range []struct {
		json string
		id   *int
	}{
		{`{"name":"通信費","account_type":"expense","display_order":1}`, &expenseID},
		{`{"name":"普通預金","account_type":"asset","display_order":2}`, &assetID},
	} {
		req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBufferString(body.json))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("setup: failed to create account: %s", w.Body.String())
		}
		var resp map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("setup: failed to decode: %v", err)
		}
		*body.id = int(resp["id"].(float64))
	}
	return
}

func createJournalEntryJSON(expenseID, assetID int) string {
	return fmt.Sprintf(`{
		"date":"2026-06-18",
		"description":"テスト仕訳",
		"receipt_required":false,
		"memo":"",
		"lines":[
			{"account_id":%d,"debit_amount":10000,"credit_amount":0},
			{"account_id":%d,"debit_amount":0,"credit_amount":10000}
		]
	}`, expenseID, assetID)
}

func TestJournalEntryHandler_List(t *testing.T) {
	r := setupJournalEntryTestHandler(t)

	t.Run("empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/journal-entries", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var entries []json.RawMessage
		if err := json.NewDecoder(w.Body).Decode(&entries); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if len(entries) != 0 {
			t.Errorf("expected empty array, got %d items", len(entries))
		}
	})

	t.Run("with entries", func(t *testing.T) {
		expenseID, assetID := createTestAccountsViaAPI(t, r)
		body := createJournalEntryJSON(expenseID, assetID)
		req := httptest.NewRequest(http.MethodPost, "/api/journal-entries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("setup: expected 201, got %d: %s", w.Code, w.Body.String())
		}

		req = httptest.NewRequest(http.MethodGet, "/api/journal-entries", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var entries []json.RawMessage
		if err := json.NewDecoder(w.Body).Decode(&entries); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if len(entries) < 1 {
			t.Errorf("expected at least 1 entry, got %d", len(entries))
		}
	})
}

func TestJournalEntryHandler_Create(t *testing.T) {
	r := setupJournalEntryTestHandler(t)
	expenseID, assetID := createTestAccountsViaAPI(t, r)

	t.Run("success", func(t *testing.T) {
		body := createJournalEntryJSON(expenseID, assetID)
		req := httptest.NewRequest(http.MethodPost, "/api/journal-entries", bytes.NewBufferString(body))
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
		if resp["id"] == nil || resp["id"].(float64) == 0 {
			t.Error("expected non-zero id")
		}
		if resp["lines"] == nil {
			t.Error("expected lines in response")
		}
	})

	t.Run("empty date", func(t *testing.T) {
		body := fmt.Sprintf(`{"date":"","description":"test","lines":[{"account_id":%d,"debit_amount":10000,"credit_amount":0},{"account_id":%d,"debit_amount":0,"credit_amount":10000}]}`, expenseID, assetID)
		req := httptest.NewRequest(http.MethodPost, "/api/journal-entries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("empty description", func(t *testing.T) {
		body := fmt.Sprintf(`{"date":"2026-06-18","description":"","lines":[{"account_id":%d,"debit_amount":10000,"credit_amount":0},{"account_id":%d,"debit_amount":0,"credit_amount":10000}]}`, expenseID, assetID)
		req := httptest.NewRequest(http.MethodPost, "/api/journal-entries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("unbalanced", func(t *testing.T) {
		body := fmt.Sprintf(`{"date":"2026-06-18","description":"test","lines":[{"account_id":%d,"debit_amount":10000,"credit_amount":0},{"account_id":%d,"debit_amount":0,"credit_amount":8000}]}`, expenseID, assetID)
		req := httptest.NewRequest(http.MethodPost, "/api/journal-entries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("insufficient lines", func(t *testing.T) {
		body := fmt.Sprintf(`{"date":"2026-06-18","description":"test","lines":[{"account_id":%d,"debit_amount":10000,"credit_amount":0}]}`, expenseID)
		req := httptest.NewRequest(http.MethodPost, "/api/journal-entries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/journal-entries", bytes.NewBufferString("{bad"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestJournalEntryHandler_Get(t *testing.T) {
	r := setupJournalEntryTestHandler(t)
	expenseID, assetID := createTestAccountsViaAPI(t, r)

	// Create an entry first
	body := createJournalEntryJSON(expenseID, assetID)
	req := httptest.NewRequest(http.MethodPost, "/api/journal-entries", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("setup: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var created map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("setup: failed to decode: %v", err)
	}
	id := int(created["id"].(float64))

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/journal-entries/%d", id), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var resp map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if resp["description"] != "テスト仕訳" {
			t.Errorf("expected description テスト仕訳, got %v", resp["description"])
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/journal-entries/99999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/journal-entries/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestJournalEntryHandler_Update(t *testing.T) {
	r := setupJournalEntryTestHandler(t)
	expenseID, assetID := createTestAccountsViaAPI(t, r)

	// Create an entry
	body := createJournalEntryJSON(expenseID, assetID)
	req := httptest.NewRequest(http.MethodPost, "/api/journal-entries", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("setup: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var created map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("setup: failed to decode: %v", err)
	}
	id := int(created["id"].(float64))

	t.Run("success", func(t *testing.T) {
		updateBody := fmt.Sprintf(`{
			"date":"2026-06-19",
			"description":"更新後の摘要",
			"receipt_required":true,
			"memo":"メモ追加",
			"lines":[
				{"account_id":%d,"debit_amount":20000,"credit_amount":0},
				{"account_id":%d,"debit_amount":0,"credit_amount":20000}
			]
		}`, expenseID, assetID)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/journal-entries/%d", id), bytes.NewBufferString(updateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("not found", func(t *testing.T) {
		updateBody := fmt.Sprintf(`{"date":"2026-06-18","description":"test","lines":[{"account_id":%d,"debit_amount":10000,"credit_amount":0},{"account_id":%d,"debit_amount":0,"credit_amount":10000}]}`, expenseID, assetID)
		req := httptest.NewRequest(http.MethodPut, "/api/journal-entries/99999", bytes.NewBufferString(updateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("empty description", func(t *testing.T) {
		updateBody := fmt.Sprintf(`{"date":"2026-06-18","description":"","lines":[{"account_id":%d,"debit_amount":10000,"credit_amount":0},{"account_id":%d,"debit_amount":0,"credit_amount":10000}]}`, expenseID, assetID)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/journal-entries/%d", id), bytes.NewBufferString(updateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestJournalEntryHandler_Delete(t *testing.T) {
	r := setupJournalEntryTestHandler(t)
	expenseID, assetID := createTestAccountsViaAPI(t, r)

	// Create an entry
	body := createJournalEntryJSON(expenseID, assetID)
	req := httptest.NewRequest(http.MethodPost, "/api/journal-entries", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("setup: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var created map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("setup: failed to decode: %v", err)
	}
	deleteID := int(created["id"].(float64))

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/journal-entries/%d", deleteID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/journal-entries/99999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/journal-entries/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}
