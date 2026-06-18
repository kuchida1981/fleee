import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { AccountForm } from './AccountForm';

describe('AccountForm', () => {
  it('renders create mode', () => {
    render(<AccountForm open={true} onOpenChange={vi.fn()} onSubmit={vi.fn()} />);

    expect(screen.getByText('勘定科目を作成')).toBeInTheDocument();
    expect(screen.getByLabelText('科目名')).toBeInTheDocument();
    expect(screen.getByText('作成')).toBeInTheDocument();
  });

  it('renders edit mode with pre-filled values', () => {
    const account = {
      id: 1,
      name: '普通預金',
      account_type: 'asset' as const,
      display_order: 0,
      normal_balance: 'debit' as const,
      statement_type: 'balance_sheet' as const,
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    };

    render(<AccountForm open={true} onOpenChange={vi.fn()} account={account} onSubmit={vi.fn()} />);

    expect(screen.getByText('勘定科目を編集')).toBeInTheDocument();
    expect(screen.getByDisplayValue('普通預金')).toBeInTheDocument();
    expect(screen.getByText('更新')).toBeInTheDocument();
  });

  it('shows validation error for empty name', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();

    render(<AccountForm open={true} onOpenChange={vi.fn()} onSubmit={onSubmit} />);

    await user.click(screen.getByText('作成'));

    expect(screen.getByText('科目名は必須です')).toBeInTheDocument();
    expect(onSubmit).not.toHaveBeenCalled();
  });
});
