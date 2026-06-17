package importer

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/kosuke/fleee/internal/model"
	"github.com/kosuke/fleee/internal/store"
)

// ImportResult contains statistics of the import operation
type ImportResult struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Skipped int `json:"skipped"`
}

// AccountImporter handles importing accounts from formatted CSV/TSV data
type AccountImporter struct {
	store *store.AccountStore
}

// NewAccountImporter creates a new AccountImporter
func NewAccountImporter(store *store.AccountStore) *AccountImporter {
	return &AccountImporter{store: store}
}

// Import parses the CSV or TSV data and imports accounts into the store
func (imp *AccountImporter) Import(ctx context.Context, r io.Reader, isTSV bool) (*ImportResult, error) {
	reader := csv.NewReader(r)
	if isTSV {
		reader.Comma = '\t'
	}
	// Allow flexible field counts
	reader.FieldsPerRecord = -1

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Identify column indexes for "科目名", "科目貸借タイプ", "出力順番", "精算種別"
	nameIdx, typeIdx, orderIdx, statementIdx := -1, -1, -1, -1
	for i, h := range headers {
		h = strings.TrimSpace(h)
		// Byte Order Mark (BOM) removal if any
		h = strings.TrimPrefix(h, "\ufeff")
		switch h {
		case "科目名":
			nameIdx = i
		case "科目貸借タイプ":
			typeIdx = i
		case "出力順番":
			orderIdx = i
		case "精算種別":
			statementIdx = i
		}
	}

	if nameIdx == -1 || typeIdx == -1 || orderIdx == -1 || statementIdx == -1 {
		return nil, errors.New("invalid file format: missing required columns (科目名, 科目貸借タイプ, 出力順番, 精算種別)")
	}

	// Pre-load existing accounts for fast duplicate checks
	existingList, err := imp.store.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load existing accounts: %w", err)
	}

	existingNames := make(map[string]bool)
	for _, acc := range existingList {
		existingNames[acc.Name] = true
	}

	var successCount, skippedCount, totalCount int
	lineNum := 1 // Header is line 1

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read data record at line %d: %w", lineNum+1, err)
		}
		lineNum++

		// Skip empty lines
		isEmpty := true
		for _, field := range record {
			if strings.TrimSpace(field) != "" {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			continue
		}

		totalCount++

		// Verify record layout
		maxIdx := nameIdx
		if typeIdx > maxIdx {
			maxIdx = typeIdx
		}
		if orderIdx > maxIdx {
			maxIdx = orderIdx
		}
		if statementIdx > maxIdx {
			maxIdx = statementIdx
		}

		if len(record) <= maxIdx {
			return nil, fmt.Errorf("invalid format: record has fewer fields than required at line %d", lineNum)
		}

		name := strings.TrimSpace(record[nameIdx])
		balanceType := strings.TrimSpace(record[typeIdx])
		orderStr := strings.TrimSpace(record[orderIdx])
		statementType := strings.TrimSpace(record[statementIdx])

		if name == "" {
			return nil, fmt.Errorf("invalid data: account name is empty at line %d", lineNum)
		}

		// Check for duplicates
		if existingNames[name] {
			skippedCount++
			continue
		}

		displayOrder, _ := strconv.Atoi(orderStr)

		accountType, err := resolveAccountType(name, balanceType, statementType)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve account type at line %d: %w", lineNum, err)
		}

		acc := &model.Account{
			Name:         name,
			AccountType:  accountType,
			DisplayOrder: displayOrder,
		}

		err = imp.store.Create(ctx, acc)
		if err != nil {
			if errors.Is(err, store.ErrDuplicateName) {
				skippedCount++
				continue
			}
			return nil, fmt.Errorf("failed to create account for '%s' at line %d: %w", name, lineNum, err)
		}

		existingNames[name] = true
		successCount++
	}

	return &ImportResult{
		Total:   totalCount,
		Success: successCount,
		Skipped: skippedCount,
	}, nil
}

func resolveAccountType(name, balanceType, statementType string) (model.AccountType, error) {
	balanceType = strings.TrimSpace(balanceType)
	statementType = strings.TrimSpace(statementType)

	if balanceType == "借方" && statementType == "貸借対照表" {
		return model.AccountTypeAsset, nil
	}
	if balanceType == "借方" && statementType == "損益計算書" {
		return model.AccountTypeExpense, nil
	}
	if balanceType == "貸方" && statementType == "損益計算書" {
		return model.AccountTypeRevenue, nil
	}
	if balanceType == "貸方" && statementType == "貸借対照表" {
		// Known equity accounts in Japan
		if name == "元入金" || name == "事業主借" || name == "事業主貸" {
			return model.AccountTypeEquity, nil
		}
		return model.AccountTypeLiability, nil
	}

	return "", fmt.Errorf("unsupported type mapping: balance_type=%s, statement_type=%s", balanceType, statementType)
}
