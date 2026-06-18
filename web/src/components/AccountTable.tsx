import type { Account, AccountType } from '@/types/account';
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';

interface AccountTableProps {
  accounts: Account[];
  onEdit: (account: Account) => void;
  onDelete: (account: Account) => void;
}

const ACCOUNT_TYPE_LABELS: Record<AccountType, string> = {
  asset: '資産',
  liability: '負債',
  equity: '純資産',
  revenue: '収益',
  expense: '費用',
};

export function AccountTable({ accounts, onEdit, onDelete }: AccountTableProps) {
  if (accounts.length === 0) {
    return (
      <div className="rounded-lg border border-dashed border-neutral-200 p-12 text-center text-neutral-500 dark:border-neutral-800 dark:text-neutral-400">
        勘定科目が登録されていません
      </div>
    );
  }

  return (
    <div className="rounded-md border border-neutral-200 bg-white dark:border-neutral-800 dark:bg-neutral-950">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>科目名</TableHead>
            <TableHead>勘定区分</TableHead>
            <TableHead>表示順序</TableHead>
            <TableHead className="text-right">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {accounts.map((account) => (
            <TableRow key={account.id}>
              <TableCell className="font-medium">{account.name}</TableCell>
              <TableCell>{ACCOUNT_TYPE_LABELS[account.account_type]}</TableCell>
              <TableCell>{account.display_order}</TableCell>
              <TableCell className="space-x-2 text-right">
                <Button variant="outline" size="sm" onClick={() => onEdit(account)}>
                  編集
                </Button>
                <Button variant="destructive" size="sm" onClick={() => onDelete(account)}>
                  削除
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
