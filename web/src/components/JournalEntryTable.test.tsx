import { render, screen } from '@testing-library/react';
import { JournalEntryTable } from './JournalEntryTable';
import type { JournalEntry } from '@/types/journalEntry';

const mockEntries: JournalEntry[] = [
  {
    id: 1,
    date: '2026-06-18',
    description: '携帯電話料金',
    receipt_required: false,
    memo: '',
    lines: [
      { id: 1, journal_entry_id: 1, account_id: 1, account_name: '通信費', debit_amount: 10000, credit_amount: 0 },
      { id: 2, journal_entry_id: 1, account_id: 2, account_name: '普通預金', debit_amount: 0, credit_amount: 10000 },
    ],
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
];

describe('JournalEntryTable', () => {
  it('renders entry list', () => {
    render(<JournalEntryTable entries={mockEntries} onEdit={vi.fn()} onDelete={vi.fn()} />);

    expect(screen.getByText('2026-06-18')).toBeInTheDocument();
    expect(screen.getByText('携帯電話料金')).toBeInTheDocument();
    expect(screen.getByText('10,000')).toBeInTheDocument();
  });

  it('shows empty state message', () => {
    render(<JournalEntryTable entries={[]} onEdit={vi.fn()} onDelete={vi.fn()} />);

    expect(screen.getByText('仕訳が登録されていません')).toBeInTheDocument();
  });

  it('renders edit and delete buttons', () => {
    render(<JournalEntryTable entries={mockEntries} onEdit={vi.fn()} onDelete={vi.fn()} />);

    const editButtons = screen.getAllByRole('button', { name: '編集' });
    const deleteButtons = screen.getAllByRole('button', { name: '削除' });

    expect(editButtons).toHaveLength(1);
    expect(deleteButtons).toHaveLength(1);
  });
});
