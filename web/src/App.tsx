import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { AccountTable } from '@/components/AccountTable';
import { AccountForm } from '@/components/AccountForm';
import { ImportDialog } from '@/components/ImportDialog';
import { listAccounts, createAccount, updateAccount, deleteAccount } from '@/api/accounts';
import type { Account, AccountType } from '@/types/account';

function App() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // ダイアログ開閉ステート
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isImportOpen, setIsImportOpen] = useState(false);

  // 編集対象（新規作成時は null または undefined）
  const [editingAccount, setEditingAccount] = useState<Account | undefined>(undefined);

  // 一覧取得
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

  useEffect(() => {
    let active = true;
    const load = async () => {
      // 最初の状態更新を非同期コンテキストで実行し、同期的な setState 警告を回避する
      await new Promise((resolve) => setTimeout(resolve, 0));
      if (!active) return;
      setIsLoading(true);
      setError(null);
      try {
        const data = await listAccounts();
        if (!active) return;
        const sorted = [...data].sort((a, b) => a.display_order - b.display_order);
        setAccounts(sorted);
      } catch (err) {
        if (!active) return;
        console.error(err);
        setError('データの取得に失敗しました。');
      } finally {
        if (active) setIsLoading(false);
      }
    };
    load();
    return () => {
      active = false;
    };
  }, []);

  // 作成または編集
  const handleFormSubmit = async (data: {
    name: string;
    account_type: AccountType;
    display_order: number;
  }) => {
    try {
      if (editingAccount) {
        // 編集
        await updateAccount(editingAccount.id, data);
      } else {
        // 新規作成
        await createAccount(data);
      }
      setIsFormOpen(false);
      fetchAccounts();
    } catch (err) {
      console.error(err);
      alert('保存に失敗しました。');
    }
  };

  // 編集開始
  const handleEdit = (account: Account) => {
    setEditingAccount(account);
    setIsFormOpen(true);
  };

  // 削除
  const handleDelete = async (account: Account) => {
    if (window.confirm(`勘定科目「${account.name}」を削除してもよろしいですか？`)) {
      try {
        await deleteAccount(account.id);
        fetchAccounts();
      } catch (err) {
        console.error(err);
        alert('削除に失敗しました。');
      }
    }
  };

  // 新規作成ダイアログを開く
  const handleCreateNew = () => {
    setEditingAccount(undefined);
    setIsFormOpen(true);
  };

  return (
    <div className="min-h-screen bg-neutral-50 dark:bg-neutral-900 text-neutral-900 dark:text-neutral-50 font-sans">
      <header className="border-b border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-950 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-2xl font-bold tracking-tight text-neutral-950 dark:text-neutral-50">fleee</span>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="space-y-6">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <h1 className="text-3xl font-bold tracking-tight text-neutral-950 dark:text-neutral-50">
              勘定科目管理
            </h1>
            <div className="flex items-center gap-2">
              <Button variant="outline" onClick={() => setIsImportOpen(true)}>
                インポート
              </Button>
              <Button onClick={handleCreateNew}>
                勘定科目を作成
              </Button>
            </div>
          </div>

          {error && (
            <div className="rounded-lg bg-destructive/10 text-destructive p-4 text-sm font-medium">
              {error}
            </div>
          )}

          {isLoading ? (
            <div className="text-center py-12 text-neutral-500">
              読み込み中...
            </div>
          ) : (
            <AccountTable
              accounts={accounts}
              onEdit={handleEdit}
              onDelete={handleDelete}
            />
          )}
        </div>
      </main>

      {/* フォームダイアログ */}
      {isFormOpen && (
        <AccountForm
          key={editingAccount?.id ?? 'new'}
          open={isFormOpen}
          onOpenChange={setIsFormOpen}
          account={editingAccount}
          onSubmit={handleFormSubmit}
        />
      )}

      {/* インポートダイアログ */}
      {isImportOpen && (
        <ImportDialog
          key="import-dialog"
          open={isImportOpen}
          onOpenChange={setIsImportOpen}
          onImportComplete={fetchAccounts}
        />
      )}
    </div>
  );
}

export default App;
