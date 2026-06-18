package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kosuke/fleee/internal/model"
	"github.com/kosuke/fleee/internal/store"
	"github.com/kosuke/fleee/internal/testutil"
)

func TestAccountStore_Create(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := store.NewAccountStore(db)
	ctx := context.Background()

	// 1. Success case
	acc := &model.Account{
		Name:         "Cash",
		AccountType:  model.AccountTypeAsset,
		DisplayOrder: 10,
	}
	err := s.Create(ctx, acc)
	if err != nil {
		t.Fatalf("unexpected error creating account: %v", err)
	}
	if acc.ID == 0 {
		t.Errorf("expected generated ID, got 0")
	}

	// 2. Duplicate name error case
	dupAcc := &model.Account{
		Name:         "Cash",
		AccountType:  model.AccountTypeAsset,
		DisplayOrder: 20,
	}
	err = s.Create(ctx, dupAcc)
	if !errors.Is(err, store.ErrDuplicateName) {
		t.Errorf("expected ErrDuplicateName, got %v", err)
	}
}

func TestAccountStore_GetByID(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := store.NewAccountStore(db)
	ctx := context.Background()

	// Insert test data
	acc := &model.Account{
		Name:         "Accounts Receivable",
		AccountType:  model.AccountTypeAsset,
		DisplayOrder: 10,
	}
	if err := s.Create(ctx, acc); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// 1. Success case
	got, err := s.GetByID(ctx, acc.ID)
	if err != nil {
		t.Fatalf("unexpected error getting account: %v", err)
	}
	if got.Name != acc.Name || got.AccountType != acc.AccountType {
		t.Errorf("got %+v, expected %+v", got, acc)
	}

	// 2. Not found case
	_, err = s.GetByID(ctx, 99999)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAccountStore_ListAll(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := store.NewAccountStore(db)
	ctx := context.Background()

	// 1. Zero accounts case
	list, err := s.ListAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing accounts: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected 0 accounts, got %d", len(list))
	}

	// Insert test data
	acc1 := &model.Account{
		Name:         "Revenue",
		AccountType:  model.AccountTypeRevenue,
		DisplayOrder: 20,
	}
	acc2 := &model.Account{
		Name:         "Cash",
		AccountType:  model.AccountTypeAsset,
		DisplayOrder: 10,
	}
	if err := s.Create(ctx, acc1); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if err := s.Create(ctx, acc2); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// 2. Success case: Multiple accounts, ordered by display_order
	list, err = s.ListAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing accounts: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(list))
	}
	// display_order: acc2 (10) < acc1 (20)
	if list[0].Name != "Cash" || list[1].Name != "Revenue" {
		t.Errorf("list order is incorrect: first is %s, second is %s", list[0].Name, list[1].Name)
	}
}

func TestAccountStore_Update(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := store.NewAccountStore(db)
	ctx := context.Background()

	// Insert test data
	acc1 := &model.Account{
		Name:         "Cash",
		AccountType:  model.AccountTypeAsset,
		DisplayOrder: 10,
	}
	acc2 := &model.Account{
		Name:         "Bank",
		AccountType:  model.AccountTypeAsset,
		DisplayOrder: 20,
	}
	if err := s.Create(ctx, acc1); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if err := s.Create(ctx, acc2); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// 1. Success case
	acc1.Name = "Petty Cash"
	acc1.DisplayOrder = 5
	err := s.Update(ctx, acc1)
	if err != nil {
		t.Fatalf("unexpected error updating account: %v", err)
	}

	got, err := s.GetByID(ctx, acc1.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated account: %v", err)
	}
	if got.Name != "Petty Cash" || got.DisplayOrder != 5 {
		t.Errorf("update failed to save attributes, got name: %s, display_order: %d", got.Name, got.DisplayOrder)
	}

	// 2. Not found case
	notFoundAcc := &model.Account{
		ID:           99999,
		Name:         "Sales",
		AccountType:  model.AccountTypeRevenue,
		DisplayOrder: 30,
	}
	err = s.Update(ctx, notFoundAcc)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	// 3. Duplicate name error case
	acc1.Name = "Bank" // Exists on acc2
	err = s.Update(ctx, acc1)
	if !errors.Is(err, store.ErrDuplicateName) {
		t.Errorf("expected ErrDuplicateName, got %v", err)
	}
}

func TestAccountStore_Delete(t *testing.T) {
	db := testutil.NewTestDB(t)
	s := store.NewAccountStore(db)
	ctx := context.Background()

	// Insert test data
	acc := &model.Account{
		Name:         "Accounts Payable",
		AccountType:  model.AccountTypeLiability,
		DisplayOrder: 10,
	}
	if err := s.Create(ctx, acc); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// 1. Success case
	err := s.Delete(ctx, acc.ID)
	if err != nil {
		t.Fatalf("unexpected error deleting account: %v", err)
	}

	// Verify deletion
	_, err = s.GetByID(ctx, acc.ID)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound after deletion, got %v", err)
	}

	// 2. Not found case
	err = s.Delete(ctx, 99999)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
