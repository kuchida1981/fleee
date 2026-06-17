export type AccountType = 'asset' | 'liability' | 'equity' | 'revenue' | 'expense';

export interface Account {
  id: number;
  name: string;
  account_type: AccountType;
  display_order: number;
  normal_balance: 'debit' | 'credit';
  statement_type: 'balance_sheet' | 'income_statement';
  created_at: string;
  updated_at: string;
}

export interface ImportResult {
  total: number;
  success: number;
  skipped: number;
}
