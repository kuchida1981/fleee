package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kosuke/fleee/internal/model"
)

var (
	// ErrUnbalanced is returned when the sum of debits does not equal the sum of credits
	ErrUnbalanced = errors.New("journal entry is not balanced")
	// ErrInsufficientLines is returned when a journal entry has fewer than 2 lines
	ErrInsufficientLines = errors.New("journal entry must have at least 2 lines")
)

// JournalEntryStore manages journal entry persistent operations in SQLite
type JournalEntryStore struct {
	db *DB
}

// NewJournalEntryStore creates a new JournalEntryStore
func NewJournalEntryStore(db *DB) *JournalEntryStore {
	return &JournalEntryStore{db: db}
}

func validateJournalLines(lines []model.JournalLine) error {
	if len(lines) < 2 {
		return ErrInsufficientLines
	}
	var totalDebit, totalCredit int64
	for _, line := range lines {
		totalDebit += line.DebitAmount
		totalCredit += line.CreditAmount
	}
	if totalDebit != totalCredit {
		return ErrUnbalanced
	}
	return nil
}

// Create inserts a new journal entry and its lines in a transaction
func (s *JournalEntryStore) Create(ctx context.Context, entry *model.JournalEntry) (err error) {
	if err := validateJournalLines(entry.Lines); err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	now := time.Now().UTC().Format(time.RFC3339)
	query := `
		INSERT INTO journal_entries (date, description, receipt_required, memo, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := tx.ExecContext(ctx, query, entry.Date, entry.Description, entry.ReceiptRequired, entry.Memo, now, now)
	if err != nil {
		return fmt.Errorf("failed to insert journal entry: %w", err)
	}

	entryID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	entry.ID = entryID

	// Parse timestamps back to model
	entry.CreatedAt, _ = time.Parse(time.RFC3339, now)
	entry.UpdatedAt, _ = time.Parse(time.RFC3339, now)

	lineQuery := `
		INSERT INTO journal_lines (journal_entry_id, account_id, debit_amount, credit_amount, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	for i := range entry.Lines {
		var resLine sql.Result
		resLine, err = tx.ExecContext(ctx, lineQuery, entryID, entry.Lines[i].AccountID, entry.Lines[i].DebitAmount, entry.Lines[i].CreditAmount, now, now)
		if err != nil {
			return fmt.Errorf("failed to insert journal line: %w", err)
		}
		var lineID int64
		lineID, err = resLine.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get line last insert ID: %w", err)
		}
		entry.Lines[i].ID = lineID
		entry.Lines[i].JournalEntryID = entryID
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// GetByID retrieves a single journal entry by ID including its lines and account names
func (s *JournalEntryStore) GetByID(ctx context.Context, id int64) (*model.JournalEntry, error) {
	query := `
		SELECT id, date, description, receipt_required, memo, created_at, updated_at
		FROM journal_entries
		WHERE id = ?
	`
	var entry model.JournalEntry
	var createdAtStr, updatedAtStr string
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&entry.ID, &entry.Date, &entry.Description, &entry.ReceiptRequired, &entry.Memo, &createdAtStr, &updatedAtStr,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get journal entry: %w", err)
	}

	// Parse timestamps
	entry.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	}
	entry.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
	if entry.UpdatedAt.IsZero() {
		entry.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	}

	// Retrieve associated journal lines
	linesQuery := `
		SELECT jl.id, jl.journal_entry_id, jl.account_id, COALESCE(a.name, ''), jl.debit_amount, jl.credit_amount
		FROM journal_lines jl
		LEFT JOIN accounts a ON jl.account_id = a.id
		WHERE jl.journal_entry_id = ?
		ORDER BY jl.id ASC
	`
	rows, err := s.db.QueryContext(ctx, linesQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get journal lines: %w", err)
	}
	defer func() { _ = rows.Close() }()

	entry.Lines = []model.JournalLine{}
	for rows.Next() {
		var line model.JournalLine
		err := rows.Scan(&line.ID, &line.JournalEntryID, &line.AccountID, &line.AccountName, &line.DebitAmount, &line.CreditAmount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan journal line: %w", err)
		}
		entry.Lines = append(entry.Lines, line)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("journal lines rows error: %w", err)
	}

	return &entry, nil
}

// ListAll retrieves all journal entries ordered by date DESC then id DESC, including lines
func (s *JournalEntryStore) ListAll(ctx context.Context) ([]*model.JournalEntry, error) {
	query := `
		SELECT id, date, description, receipt_required, memo, created_at, updated_at
		FROM journal_entries
		ORDER BY date DESC, id DESC
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list journal entries: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var entries []*model.JournalEntry
	for rows.Next() {
		var entry model.JournalEntry
		var createdAtStr, updatedAtStr string
		err := rows.Scan(&entry.ID, &entry.Date, &entry.Description, &entry.ReceiptRequired, &entry.Memo, &createdAtStr, &updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan journal entry: %w", err)
		}

		entry.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		if entry.CreatedAt.IsZero() {
			entry.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}
		entry.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
		if entry.UpdatedAt.IsZero() {
			entry.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		}

		entries = append(entries, &entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("journal entries rows error: %w", err)
	}

	// Retrieve lines for each entry (N+1 query is fine for this scale)
	linesQuery := `
		SELECT jl.id, jl.journal_entry_id, jl.account_id, COALESCE(a.name, ''), jl.debit_amount, jl.credit_amount
		FROM journal_lines jl
		LEFT JOIN accounts a ON jl.account_id = a.id
		WHERE jl.journal_entry_id = ?
		ORDER BY jl.id ASC
	`
	for _, entry := range entries {
		lrows, err := s.db.QueryContext(ctx, linesQuery, entry.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get journal lines for entry %d: %w", entry.ID, err)
		}
		entry.Lines = []model.JournalLine{}
		for lrows.Next() {
			var line model.JournalLine
			err := lrows.Scan(&line.ID, &line.JournalEntryID, &line.AccountID, &line.AccountName, &line.DebitAmount, &line.CreditAmount)
			if err != nil {
				_ = lrows.Close()
				return nil, fmt.Errorf("failed to scan journal line for entry %d: %w", entry.ID, err)
			}
			entry.Lines = append(entry.Lines, line)
		}
		_ = lrows.Close()
		if err := lrows.Err(); err != nil {
			return nil, fmt.Errorf("journal lines rows error for entry %d: %w", entry.ID, err)
		}
	}

	return entries, nil
}

// Update updates an existing journal entry and its lines in a transaction
func (s *JournalEntryStore) Update(ctx context.Context, entry *model.JournalEntry) (err error) {
	if err := validateJournalLines(entry.Lines); err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	now := time.Now().UTC().Format(time.RFC3339)
	query := `
		UPDATE journal_entries
		SET date = ?, description = ?, receipt_required = ?, memo = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := tx.ExecContext(ctx, query, entry.Date, entry.Description, entry.ReceiptRequired, entry.Memo, now, entry.ID)
	if err != nil {
		return fmt.Errorf("failed to update journal entry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	// Update timestamp in model
	entry.UpdatedAt, _ = time.Parse(time.RFC3339, now)

	// Delete existing journal lines
	deleteQuery := `DELETE FROM journal_lines WHERE journal_entry_id = ?`
	_, err = tx.ExecContext(ctx, deleteQuery, entry.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing journal lines: %w", err)
	}

	// Insert new lines
	lineQuery := `
		INSERT INTO journal_lines (journal_entry_id, account_id, debit_amount, credit_amount, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	for i := range entry.Lines {
		var resLine sql.Result
		resLine, err = tx.ExecContext(ctx, lineQuery, entry.ID, entry.Lines[i].AccountID, entry.Lines[i].DebitAmount, entry.Lines[i].CreditAmount, now, now)
		if err != nil {
			return fmt.Errorf("failed to insert journal line: %w", err)
		}
		var lineID int64
		lineID, err = resLine.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get line last insert ID: %w", err)
		}
		entry.Lines[i].ID = lineID
		entry.Lines[i].JournalEntryID = entry.ID
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Delete removes a journal entry by ID (associated lines will be deleted by CASCADE)
func (s *JournalEntryStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM journal_entries WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete journal entry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
