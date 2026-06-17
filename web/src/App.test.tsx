import { render, screen } from '@testing-library/react';

vi.mock('@/api/accounts', () => ({
  listAccounts: vi.fn().mockResolvedValue([]),
  createAccount: vi.fn(),
  updateAccount: vi.fn(),
  deleteAccount: vi.fn(),
}));

describe('App', () => {
  it('renders header and page title', async () => {
    const { default: App } = await import('./App');
    render(<App />);

    expect(screen.getByText('fleee')).toBeInTheDocument();
    expect(screen.getByText('勘定科目管理')).toBeInTheDocument();
  });

  it('renders create and import buttons', async () => {
    const { default: App } = await import('./App');
    render(<App />);

    expect(screen.getByText('勘定科目を作成')).toBeInTheDocument();
    expect(screen.getByText('インポート')).toBeInTheDocument();
  });
});
