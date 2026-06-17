import { render, screen } from '@testing-library/react';
import { AccountTable } from './AccountTable';
import type { Account } from '@/types/account';

const mockAccounts: Account[] = [
  {
    id: 1,
    name: '普通預金',
    account_type: 'asset',
    display_order: 0,
    normal_balance: 'debit',
    statement_type: 'balance_sheet',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 2,
    name: '売上',
    account_type: 'revenue',
    display_order: 1,
    normal_balance: 'credit',
    statement_type: 'income_statement',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
];

describe('AccountTable', () => {
  it('renders account list', () => {
    render(<AccountTable accounts={mockAccounts} onEdit={vi.fn()} onDelete={vi.fn()} />);

    expect(screen.getByText('普通預金')).toBeInTheDocument();
    expect(screen.getByText('売上')).toBeInTheDocument();
    expect(screen.getByText('資産')).toBeInTheDocument();
    expect(screen.getByText('収益')).toBeInTheDocument();
  });

  it('shows empty state message', () => {
    render(<AccountTable accounts={[]} onEdit={vi.fn()} onDelete={vi.fn()} />);

    expect(screen.getByText('勘定科目が登録されていません')).toBeInTheDocument();
  });

  it('renders edit and delete buttons for each row', () => {
    render(<AccountTable accounts={mockAccounts} onEdit={vi.fn()} onDelete={vi.fn()} />);

    const editButtons = screen.getAllByText('編集');
    const deleteButtons = screen.getAllByText('削除');

    expect(editButtons).toHaveLength(2);
    expect(deleteButtons).toHaveLength(2);
  });
});
