package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kosuke/fleee/internal/model"
)

var (
	// ErrNotFound is returned when an account is not found
	ErrNotFound = errors.New("account not found")
	// ErrDuplicateName is returned when an account name already exists
	ErrDuplicateName = errors.New("duplicate account name")
)

// AccountStore manages account persistent operations in SQLite
type AccountStore struct {
	db *DB
}

// NewAccountStore creates a new AccountStore
func NewAccountStore(db *DB) *AccountStore {
	return &AccountStore{db: db}
}

// Create inserts a new account record
func (s *AccountStore) Create(ctx context.Context, a *model.Account) error {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `
		INSERT INTO accounts (name, account_type, display_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := s.db.ExecContext(ctx, query, a.Name, a.AccountType, a.DisplayOrder, now, now)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrDuplicateName
		}
		return fmt.Errorf("failed to create account: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	a.ID = id
	return nil
}

// GetByID retrieves a single account by ID
func (s *AccountStore) GetByID(ctx context.Context, id int64) (*model.Account, error) {
	query := `
		SELECT id, name, account_type, display_order, created_at, updated_at
		FROM accounts
		WHERE id = ?
	`
	var a model.Account
	var createdAtStr, updatedAtStr string
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.Name, &a.AccountType, &a.DisplayOrder, &createdAtStr, &updatedAtStr,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get account by ID: %w", err)
	}

	a.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	if a.CreatedAt.IsZero() {
		a.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	}
	a.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
	if a.UpdatedAt.IsZero() {
		a.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	}

	return &a, nil
}

// ListAll retrieves all accounts ordered by display_order then id
func (s *AccountStore) ListAll(ctx context.Context) ([]*model.Account, error) {
	query := `
		SELECT id, name, account_type, display_order, created_at, updated_at
		FROM accounts
		ORDER BY display_order ASC, id ASC
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*model.Account
	for rows.Next() {
		var a model.Account
		var createdAtStr, updatedAtStr string
		err := rows.Scan(&a.ID, &a.Name, &a.AccountType, &a.DisplayOrder, &createdAtStr, &updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}

		a.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		if a.CreatedAt.IsZero() {
			a.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}
		a.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
		if a.UpdatedAt.IsZero() {
			a.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		}

		accounts = append(accounts, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return accounts, nil
}

// Update updates an existing account record
func (s *AccountStore) Update(ctx context.Context, a *model.Account) error {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `
		UPDATE accounts
		SET name = ?, account_type = ?, display_order = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := s.db.ExecContext(ctx, query, a.Name, a.AccountType, a.DisplayOrder, now, a.ID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrDuplicateName
		}
		return fmt.Errorf("failed to update account: %w", err)
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

// Delete removes an account by ID
func (s *AccountStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM accounts WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
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
