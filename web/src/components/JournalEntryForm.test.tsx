import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { JournalEntryForm } from './JournalEntryForm';
import type { Account } from '@/types/account';
import type { JournalEntry } from '@/types/journalEntry';

const mockAccounts: Account[] = [
  {
    id: 1,
    name: '通信費',
    account_type: 'expense',
    display_order: 1,
    normal_balance: 'debit',
    statement_type: 'income_statement',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 2,
    name: '普通預金',
    account_type: 'asset',
    display_order: 2,
    normal_balance: 'debit',
    statement_type: 'balance_sheet',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
];

const mockEntry: JournalEntry = {
  id: 1,
  date: '2026-06-18',
  description: '携帯電話料金',
  receipt_required: false,
  memo: '',
  lines: [
    {
      id: 1,
      journal_entry_id: 1,
      account_id: 1,
      account_name: '通信費',
      debit_amount: 10000,
      credit_amount: 0,
    },
    {
      id: 2,
      journal_entry_id: 1,
      account_id: 2,
      account_name: '普通預金',
      debit_amount: 0,
      credit_amount: 10000,
    },
  ],
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
};

describe('JournalEntryForm', () => {
  it('renders create mode', () => {
    render(
      <JournalEntryForm
        open={true}
        onOpenChange={vi.fn()}
        accounts={mockAccounts}
        onSubmit={vi.fn()}
      />,
    );

    expect(screen.getByText('仕訳を作成')).toBeInTheDocument();
    expect(screen.getByLabelText('日付')).toBeInTheDocument();
    expect(screen.getByLabelText('摘要')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '作成' })).toBeInTheDocument();
  });

  it('renders edit mode', () => {
    render(
      <JournalEntryForm
        open={true}
        onOpenChange={vi.fn()}
        accounts={mockAccounts}
        entry={mockEntry}
        onSubmit={vi.fn()}
      />,
    );

    expect(screen.getByText('仕訳を編集')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '更新' })).toBeInTheDocument();
  });

  it('shows validation error for empty description', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();

    render(
      <JournalEntryForm
        open={true}
        onOpenChange={vi.fn()}
        accounts={mockAccounts}
        onSubmit={onSubmit}
      />,
    );

    // Click submit when description is empty
    await user.click(screen.getByRole('button', { name: '作成' }));

    expect(screen.getByText('摘要は必須です')).toBeInTheDocument();
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('switches to compound mode', async () => {
    const user = userEvent.setup();

    render(
      <JournalEntryForm
        open={true}
        onOpenChange={vi.fn()}
        accounts={mockAccounts}
        onSubmit={vi.fn()}
      />,
    );

    // Initial state is single mode, "+ 行追加" button should not be present
    expect(screen.queryByRole('button', { name: '+ 行追加' })).not.toBeInTheDocument();

    // Click "複合仕訳に切替"
    await user.click(screen.getByRole('button', { name: '複合仕訳に切替' }));

    // Now "+ 行追加" and "単一仕訳に切替" should be visible
    expect(screen.getByRole('button', { name: '+ 行追加' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '単一仕訳に切替' })).toBeInTheDocument();
  });
});
