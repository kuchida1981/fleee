import { render, screen } from '@testing-library/react';

vi.mock('@/api/accounts', () => ({
  listAccounts: vi.fn().mockResolvedValue([]),
  createAccount: vi.fn(),
  updateAccount: vi.fn(),
  deleteAccount: vi.fn(),
}));

vi.mock('@/api/journalEntries', () => ({
  listJournalEntries: vi.fn().mockResolvedValue([]),
  createJournalEntry: vi.fn(),
  updateJournalEntry: vi.fn(),
  deleteJournalEntry: vi.fn(),
}));

describe('App', () => {
  it('renders header and navigation', async () => {
    const { default: App } = await import('./App');
    render(<App />);

    expect(screen.getByText('fleee')).toBeInTheDocument();
    expect(screen.getByText('仕訳')).toBeInTheDocument();
    expect(screen.getByText('勘定科目')).toBeInTheDocument();
  });

  it('renders journal entries page by default', async () => {
    const { default: App } = await import('./App');
    render(<App />);

    expect(screen.getByText('仕訳管理')).toBeInTheDocument();
    expect(screen.getByText('仕訳を作成')).toBeInTheDocument();
  });
});
