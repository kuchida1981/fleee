import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { AccountTable } from '@/components/AccountTable';
import { AccountForm } from '@/components/AccountForm';
import { ImportDialog } from '@/components/ImportDialog';
import { JournalEntryTable } from '@/components/JournalEntryTable';
import { JournalEntryForm } from '@/components/JournalEntryForm';
import { listAccounts, createAccount, updateAccount, deleteAccount } from '@/api/accounts';
import {
  listJournalEntries,
  createJournalEntry,
  updateJournalEntry,
  deleteJournalEntry,
} from '@/api/journalEntries';
import type { Account, AccountType } from '@/types/account';
import type { JournalEntry, CreateJournalEntryRequest } from '@/types/journalEntry';

type Page = 'accounts' | 'journal-entries';

function App() {
  const [page, setPage] = useState<Page>('journal-entries');

  // Accounts state
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isImportOpen, setIsImportOpen] = useState(false);
  const [editingAccount, setEditingAccount] = useState<Account | undefined>(undefined);

  // Journal entries state
  const [entries, setEntries] = useState<JournalEntry[]>([]);
  const [isEntriesLoading, setIsEntriesLoading] = useState(false);
  const [entriesError, setEntriesError] = useState<string | null>(null);
  const [isEntryFormOpen, setIsEntryFormOpen] = useState(false);
  const [editingEntry, setEditingEntry] = useState<JournalEntry | undefined>(undefined);

  // Fetch accounts
  const fetchAccounts = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await listAccounts();
      const sorted = [...data].sort((a, b) => a.display_order - b.display_order);
      setAccounts(sorted);
    } catch (err) {
      console.error(err);
      setError('データの取得に失敗しました。');
    } finally {
      setIsLoading(false);
    }
  };

  // Fetch journal entries
  const fetchEntries = async () => {
    setIsEntriesLoading(true);
    setEntriesError(null);
    try {
      const data = await listJournalEntries();
      setEntries(data);
    } catch (err) {
      console.error(err);
      setEntriesError('データの取得に失敗しました。');
    } finally {
      setIsEntriesLoading(false);
    }
  };

  useEffect(() => {
    let active = true;
    const load = async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
      if (!active) return;
      setIsLoading(true);
      setIsEntriesLoading(true);
      setError(null);
      setEntriesError(null);
      try {
        const [accountData, entryData] = await Promise.all([listAccounts(), listJournalEntries()]);
        if (!active) return;
        const sorted = [...accountData].sort((a, b) => a.display_order - b.display_order);
        setAccounts(sorted);
        setEntries(entryData);
      } catch (err) {
        if (!active) return;
        console.error(err);
        setError('データの取得に失敗しました。');
        setEntriesError('データの取得に失敗しました。');
      } finally {
        if (active) {
          setIsLoading(false);
          setIsEntriesLoading(false);
        }
      }
    };
    load();
    return () => {
      active = false;
    };
  }, []);

  // Account handlers
  const handleAccountFormSubmit = async (data: {
    name: string;
    account_type: AccountType;
    display_order: number;
  }) => {
    try {
      if (editingAccount) {
        await updateAccount(editingAccount.id, data);
      } else {
        await createAccount(data);
      }
      setIsFormOpen(false);
      fetchAccounts();
    } catch (err) {
      console.error(err);
      alert('保存に失敗しました。');
    }
  };

  const handleEditAccount = (account: Account) => {
    setEditingAccount(account);
    setIsFormOpen(true);
  };

  const handleDeleteAccount = async (account: Account) => {
    if (window.confirm(`勘定科目「${account.name}」を削除してもよろしいですか？`)) {
      try {
        await deleteAccount(account.id);
        fetchAccounts();
      } catch (err) {
        console.error(err);
        alert('削除に失敗しました。仕訳で使用中の科目は削除できません。');
      }
    }
  };

  const handleCreateNewAccount = () => {
    setEditingAccount(undefined);
    setIsFormOpen(true);
  };

  // Journal entry handlers
  const handleEntryFormSubmit = async (data: CreateJournalEntryRequest) => {
    try {
      if (editingEntry) {
        await updateJournalEntry(editingEntry.id, data);
      } else {
        await createJournalEntry(data);
      }
      setIsEntryFormOpen(false);
      fetchEntries();
    } catch (err) {
      console.error(err);
      alert('保存に失敗しました。');
    }
  };

  const handleEditEntry = (entry: JournalEntry) => {
    setEditingEntry(entry);
    setIsEntryFormOpen(true);
  };

  const handleDeleteEntry = async (entry: JournalEntry) => {
    if (window.confirm(`仕訳「${entry.description}」を削除してもよろしいですか？`)) {
      try {
        await deleteJournalEntry(entry.id);
        fetchEntries();
      } catch (err) {
        console.error(err);
        alert('削除に失敗しました。');
      }
    }
  };

  const handleCreateNewEntry = () => {
    setEditingEntry(undefined);
    setIsEntryFormOpen(true);
  };

  return (
    <div className="min-h-screen bg-neutral-50 font-sans text-neutral-900 dark:bg-neutral-900 dark:text-neutral-50">
      <header className="sticky top-0 z-50 border-b border-neutral-200 bg-white dark:border-neutral-800 dark:bg-neutral-950">
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-6">
            <span className="text-2xl font-bold tracking-tight text-neutral-950 dark:text-neutral-50">
              fleee
            </span>
            <nav className="flex gap-1">
              <Button
                variant={page === 'journal-entries' ? 'default' : 'ghost'}
                size="sm"
                onClick={() => setPage('journal-entries')}
              >
                仕訳
              </Button>
              <Button
                variant={page === 'accounts' ? 'default' : 'ghost'}
                size="sm"
                onClick={() => setPage('accounts')}
              >
                勘定科目
              </Button>
            </nav>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {page === 'accounts' && (
          <div className="space-y-6">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
              <h1 className="text-3xl font-bold tracking-tight text-neutral-950 dark:text-neutral-50">
                勘定科目管理
              </h1>
              <div className="flex items-center gap-2">
                <Button variant="outline" onClick={() => setIsImportOpen(true)}>
                  インポート
                </Button>
                <Button onClick={handleCreateNewAccount}>勘定科目を作成</Button>
              </div>
            </div>

            {error && (
              <div className="bg-destructive/10 text-destructive rounded-lg p-4 text-sm font-medium">
                {error}
              </div>
            )}

            {isLoading ? (
              <div className="py-12 text-center text-neutral-500">読み込み中...</div>
            ) : (
              <AccountTable
                accounts={accounts}
                onEdit={handleEditAccount}
                onDelete={handleDeleteAccount}
              />
            )}
          </div>
        )}

        {page === 'journal-entries' && (
          <div className="space-y-6">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
              <h1 className="text-3xl font-bold tracking-tight text-neutral-950 dark:text-neutral-50">
                仕訳管理
              </h1>
              <Button onClick={handleCreateNewEntry}>仕訳を作成</Button>
            </div>

            {entriesError && (
              <div className="bg-destructive/10 text-destructive rounded-lg p-4 text-sm font-medium">
                {entriesError}
              </div>
            )}

            {isEntriesLoading ? (
              <div className="py-12 text-center text-neutral-500">読み込み中...</div>
            ) : (
              <JournalEntryTable
                entries={entries}
                onEdit={handleEditEntry}
                onDelete={handleDeleteEntry}
              />
            )}
          </div>
        )}
      </main>

      {/* Account form dialog */}
      {isFormOpen && (
        <AccountForm
          key={editingAccount?.id ?? 'new'}
          open={isFormOpen}
          onOpenChange={setIsFormOpen}
          account={editingAccount}
          onSubmit={handleAccountFormSubmit}
        />
      )}

      {/* Account import dialog */}
      {isImportOpen && (
        <ImportDialog
          key="import-dialog"
          open={isImportOpen}
          onOpenChange={setIsImportOpen}
          onImportComplete={fetchAccounts}
        />
      )}

      {/* Journal entry form dialog */}
      {isEntryFormOpen && (
        <JournalEntryForm
          key={editingEntry?.id ?? 'new-entry'}
          open={isEntryFormOpen}
          onOpenChange={setIsEntryFormOpen}
          accounts={accounts}
          entry={editingEntry}
          onSubmit={handleEntryFormSubmit}
        />
      )}
    </div>
  );
}

export default App;
