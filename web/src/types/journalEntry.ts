export interface JournalLine {
  id: number;
  journal_entry_id: number;
  account_id: number;
  account_name: string;
  debit_amount: number;
  credit_amount: number;
}

export interface JournalEntry {
  id: number;
  date: string;
  description: string;
  receipt_required: boolean;
  memo: string;
  lines: JournalLine[];
  created_at: string;
  updated_at: string;
}

export interface CreateJournalEntryRequest {
  date: string;
  description: string;
  receipt_required: boolean;
  memo: string;
  lines: {
    account_id: number;
    debit_amount: number;
    credit_amount: number;
  }[];
}

export type UpdateJournalEntryRequest = CreateJournalEntryRequest;
