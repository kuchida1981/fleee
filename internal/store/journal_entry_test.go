package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kosuke/fleee/internal/model"
	"github.com/kosuke/fleee/internal/store"
	"github.com/kosuke/fleee/internal/testutil"
)

func createTestAccounts(t *testing.T, s *store.AccountStore, ctx context.Context) (debitAcc, creditAcc *model.Account) {
	debitAcc = &model.Account{
		Name:         "Expense Account",
		AccountType:  model.AccountTypeExpense,
		DisplayOrder: 10,
	}
	creditAcc = &model.Account{
		Name:         "Asset Account",
		AccountType:  model.AccountTypeAsset,
		DisplayOrder: 20,
	}
	if err := s.Create(ctx, debitAcc); err != nil {
		t.Fatalf("failed to create debit account: %v", err)
	}
	if err := s.Create(ctx, creditAcc); err != nil {
		t.Fatalf("failed to create credit account: %v", err)
	}
	return debitAcc, creditAcc
}

func TestJournalEntryStore_Create(t *testing.T) {
	db := testutil.NewTestDB(t)
	accStore := store.NewAccountStore(db)
	s := store.NewJournalEntryStore(db)
	ctx := context.Background()

	debitAcc, creditAcc := createTestAccounts(t, accStore, ctx)

	// 1. Success case
	entry := &model.JournalEntry{
		Date:            "2026-06-18",
		Description:     "テスト仕訳",
		ReceiptRequired: false,
		Memo:            "",
		Lines: []model.JournalLine{
			{AccountID: debitAcc.ID, DebitAmount: 10000, CreditAmount: 0},
			{AccountID: creditAcc.ID, DebitAmount: 0, CreditAmount: 10000},
		},
	}
	err := s.Create(ctx, entry)
	if err != nil {
		t.Fatalf("unexpected error creating journal entry: %v", err)
	}
	if entry.ID == 0 {
		t.Errorf("expected generated ID, got 0")
	}
	for i, line := range entry.Lines {
		if line.ID == 0 {
			t.Errorf("expected generated line ID for line %d, got 0", i)
		}
		if line.JournalEntryID != entry.ID {
			t.Errorf("expected line JournalEntryID %d, got %d", entry.ID, line.JournalEntryID)
		}
	}

	// 2. Unbalanced error case
	unbalancedEntry := &model.JournalEntry{
		Date:            "2026-06-18",
		Description:     "アンバランス仕訳",
		ReceiptRequired: false,
		Memo:            "",
		Lines: []model.JournalLine{
			{AccountID: debitAcc.ID, DebitAmount: 10000, CreditAmount: 0},
			{AccountID: creditAcc.ID, DebitAmount: 0, CreditAmount: 9000},
		},
	}
	err = s.Create(ctx, unbalancedEntry)
	if !errors.Is(err, store.ErrUnbalanced) {
		t.Errorf("expected ErrUnbalanced, got %v", err)
	}

	// 3. Insufficient lines error case
	insufficientEntry := &model.JournalEntry{
		Date:            "2026-06-18",
		Description:     "行数不足仕訳",
		ReceiptRequired: false,
		Memo:            "",
		Lines: []model.JournalLine{
			{AccountID: debitAcc.ID, DebitAmount: 10000, CreditAmount: 0},
		},
	}
	err = s.Create(ctx, insufficientEntry)
	if !errors.Is(err, store.ErrInsufficientLines) {
		t.Errorf("expected ErrInsufficientLines, got %v", err)
	}
}

func TestJournalEntryStore_GetByID(t *testing.T) {
	db := testutil.NewTestDB(t)
	accStore := store.NewAccountStore(db)
	s := store.NewJournalEntryStore(db)
	ctx := context.Background()

	debitAcc, creditAcc := createTestAccounts(t, accStore, ctx)

	entry := &model.JournalEntry{
		Date:            "2026-06-18",
		Description:     "テスト仕訳",
		ReceiptRequired: false,
		Memo:            "",
		Lines: []model.JournalLine{
			{AccountID: debitAcc.ID, DebitAmount: 10000, CreditAmount: 0},
			{AccountID: creditAcc.ID, DebitAmount: 0, CreditAmount: 10000},
		},
	}
	if err := s.Create(ctx, entry); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// 1. Success case
	got, err := s.GetByID(ctx, entry.ID)
	if err != nil {
		t.Fatalf("unexpected error getting journal entry: %v", err)
	}
	if got.Description != entry.Description || got.Date != entry.Date {
		t.Errorf("got description %q, date %q; expected %q, %q", got.Description, got.Date, entry.Description, entry.Date)
	}
	if len(got.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got.Lines))
	}
	
	// Check AccountName is loaded
	if got.Lines[0].AccountName != debitAcc.Name {
		t.Errorf("expected AccountName %q, got %q", debitAcc.Name, got.Lines[0].AccountName)
	}
	if got.Lines[1].AccountName != creditAcc.Name {
		t.Errorf("expected AccountName %q, got %q", creditAcc.Name, got.Lines[1].AccountName)
	}

	// 2. Not found case
	_, err = s.GetByID(ctx, 99999)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestJournalEntryStore_ListAll(t *testing.T) {
	db := testutil.NewTestDB(t)
	accStore := store.NewAccountStore(db)
	s := store.NewJournalEntryStore(db)
	ctx := context.Background()

	// 1. Zero case
	list, err := s.ListAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing journal entries: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected 0 entries, got %d", len(list))
	}

	debitAcc, creditAcc := createTestAccounts(t, accStore, ctx)

	entry1 := &model.JournalEntry{
		Date:            "2026-06-18",
		Description:     "仕訳1",
		ReceiptRequired: false,
		Memo:            "",
		Lines: []model.JournalLine{
			{AccountID: debitAcc.ID, DebitAmount: 10000, CreditAmount: 0},
			{AccountID: creditAcc.ID, DebitAmount: 0, CreditAmount: 10000},
		},
	}
	entry2 := &model.JournalEntry{
		Date:            "2026-06-19",
		Description:     "仕訳2",
		ReceiptRequired: false,
		Memo:            "",
		Lines: []model.JournalLine{
			{AccountID: debitAcc.ID, DebitAmount: 5000, CreditAmount: 0},
			{AccountID: creditAcc.ID, DebitAmount: 0, CreditAmount: 5000},
		},
	}
	entry3 := &model.JournalEntry{
		Date:            "2026-06-17",
		Description:     "仕訳3",
		ReceiptRequired: false,
		Memo:            "",
		Lines: []model.JournalLine{
			{AccountID: debitAcc.ID, DebitAmount: 20000, CreditAmount: 0},
			{AccountID: creditAcc.ID, DebitAmount: 0, CreditAmount: 20000},
		},
	}

	if err := s.Create(ctx, entry1); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if err := s.Create(ctx, entry2); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if err := s.Create(ctx, entry3); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// 2. Success case: ordered by date DESC
	list, err = s.ListAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing journal entries: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(list))
	}

	// Order should be: entry2 (2026-06-19), entry1 (2026-06-18), entry3 (2026-06-17)
	if list[0].Date != "2026-06-19" || list[1].Date != "2026-06-18" || list[2].Date != "2026-06-17" {
		t.Errorf("list order is incorrect: expected dates in DESC order, got %s, %s, %s", list[0].Date, list[1].Date, list[2].Date)
	}
	if list[0].Description != "仕訳2" || list[1].Description != "仕訳1" || list[2].Description != "仕訳3" {
		t.Errorf("list order is incorrect: got descriptions %q, %q, %q", list[0].Description, list[1].Description, list[2].Description)
	}
}

func TestJournalEntryStore_Update(t *testing.T) {
	db := testutil.NewTestDB(t)
	accStore := store.NewAccountStore(db)
	s := store.NewJournalEntryStore(db)
	ctx := context.Background()

	debitAcc, creditAcc := createTestAccounts(t, accStore, ctx)

	entry := &model.JournalEntry{
		Date:            "2026-06-18",
		Description:     "テスト仕訳",
		ReceiptRequired: false,
		Memo:            "",
		Lines: []model.JournalLine{
			{AccountID: debitAcc.ID, DebitAmount: 10000, CreditAmount: 0},
			{AccountID: creditAcc.ID, DebitAmount: 0, CreditAmount: 10000},
		},
	}
	if err := s.Create(ctx, entry); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// 1. Success case: Description and Lines changed (including line count)
	extraAcc := &model.Account{
		Name:         "Extra Asset Account",
		AccountType:  model.AccountTypeAsset,
		DisplayOrder: 30,
	}
	if err := accStore.Create(ctx, extraAcc); err != nil {
		t.Fatalf("failed to create extra account: %v", err)
	}

	entry.Description = "更新後仕訳"
	entry.Lines = []model.JournalLine{
		{AccountID: debitAcc.ID, DebitAmount: 10000, CreditAmount: 0},
		{AccountID: creditAcc.ID, DebitAmount: 0, CreditAmount: 4000},
		{AccountID: extraAcc.ID, DebitAmount: 0, CreditAmount: 6000},
	}

	err := s.Update(ctx, entry)
	if err != nil {
		t.Fatalf("unexpected error updating journal entry: %v", err)
	}

	// Verify the update
	got, err := s.GetByID(ctx, entry.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated entry: %v", err)
	}
	if got.Description != "更新後仕訳" {
		t.Errorf("expected Description %q, got %q", "更新後仕訳", got.Description)
	}
	if len(got.Lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got.Lines))
	}
	if got.Lines[1].CreditAmount != 4000 || got.Lines[2].CreditAmount != 6000 {
		t.Errorf("unexpected credit amounts: line 1 = %d, line 2 = %d", got.Lines[1].CreditAmount, got.Lines[2].CreditAmount)
	}

	// 2. Not found case
	notFoundEntry := &model.JournalEntry{
		ID:              99999,
		Date:            "2026-06-18",
		Description:     "存在しない仕訳",
		ReceiptRequired: false,
		Memo:            "",
		Lines: []model.JournalLine{
			{AccountID: debitAcc.ID, DebitAmount: 10000, CreditAmount: 0},
			{AccountID: creditAcc.ID, DebitAmount: 0, CreditAmount: 10000},
		},
	}
	err = s.Update(ctx, notFoundEntry)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	// 3. Unbalanced error case
	entry.Lines = []model.JournalLine{
		{AccountID: debitAcc.ID, DebitAmount: 10000, CreditAmount: 0},
		{AccountID: creditAcc.ID, DebitAmount: 0, CreditAmount: 8000},
	}
	err = s.Update(ctx, entry)
	if !errors.Is(err, store.ErrUnbalanced) {
		t.Errorf("expected ErrUnbalanced, got %v", err)
	}
}

func TestJournalEntryStore_Delete(t *testing.T) {
	db := testutil.NewTestDB(t)
	accStore := store.NewAccountStore(db)
	s := store.NewJournalEntryStore(db)
	ctx := context.Background()

	debitAcc, creditAcc := createTestAccounts(t, accStore, ctx)

	entry := &model.JournalEntry{
		Date:            "2026-06-18",
		Description:     "テスト仕訳",
		ReceiptRequired: false,
		Memo:            "",
		Lines: []model.JournalLine{
			{AccountID: debitAcc.ID, DebitAmount: 10000, CreditAmount: 0},
			{AccountID: creditAcc.ID, DebitAmount: 0, CreditAmount: 10000},
		},
	}
	if err := s.Create(ctx, entry); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// 1. Success case
	err := s.Delete(ctx, entry.ID)
	if err != nil {
		t.Fatalf("unexpected error deleting journal entry: %v", err)
	}

	// Verify deletion
	_, err = s.GetByID(ctx, entry.ID)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound after deletion, got %v", err)
	}

	// 2. Not found case
	err = s.Delete(ctx, 99999)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
