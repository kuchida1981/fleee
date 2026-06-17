package model

import (
	"encoding/json"
	"time"
)

// AccountType represents the category of the account
type AccountType string

const (
	AccountTypeAsset     AccountType = "asset"
	AccountTypeLiability AccountType = "liability"
	AccountTypeEquity    AccountType = "equity"
	AccountTypeRevenue   AccountType = "revenue"
	AccountTypeExpense   AccountType = "expense"
)

// Account represents the domain model for chart of accounts
type Account struct {
	ID           int64       `json:"id"`
	Name         string      `json:"name"`
	AccountType  AccountType `json:"account_type"`
	DisplayOrder int         `json:"display_order"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// NormalBalance returns the normal balance side of the account ("debit" or "credit")
func (a *Account) NormalBalance() string {
	switch a.AccountType {
	case AccountTypeAsset, AccountTypeExpense:
		return "debit"
	case AccountTypeLiability, AccountTypeEquity, AccountTypeRevenue:
		return "credit"
	default:
		return ""
	}
}

// StatementType returns the financial statement category ("balance_sheet" or "income_statement")
func (a *Account) StatementType() string {
	switch a.AccountType {
	case AccountTypeAsset, AccountTypeLiability, AccountTypeEquity:
		return "balance_sheet"
	case AccountTypeRevenue, AccountTypeExpense:
		return "income_statement"
	default:
		return ""
	}
}

// MarshalJSON customizes JSON serialization to include derived attributes
func (a Account) MarshalJSON() ([]byte, error) {
	type Alias Account
	return json.Marshal(&struct {
		Alias
		NormalBalance string `json:"normal_balance"`
		StatementType string `json:"statement_type"`
	}{
		Alias:         (Alias)(a),
		NormalBalance: a.NormalBalance(),
		StatementType: a.StatementType(),
	})
}
