package model

import (
	"time"
)

// JournalEntry represents a journal entry in double-entry bookkeeping.
type JournalEntry struct {
	ID              int64         `json:"id"`
	Date            string        `json:"date"`
	Description     string        `json:"description"`
	ReceiptRequired bool          `json:"receipt_required"`
	Memo            string        `json:"memo"`
	Lines           []JournalLine `json:"lines"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// JournalLine represents a single line detail of a journal entry, representing either a debit or credit.
type JournalLine struct {
	ID             int64  `json:"id"`
	JournalEntryID int64  `json:"journal_entry_id"`
	AccountID      int64  `json:"account_id"`
	AccountName    string `json:"account_name"`
	DebitAmount    int64  `json:"debit_amount"`
	CreditAmount   int64  `json:"credit_amount"`
}
