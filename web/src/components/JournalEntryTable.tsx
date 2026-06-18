import type { JournalEntry } from '@/types/journalEntry';
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';

interface JournalEntryTableProps {
  entries: JournalEntry[];
  onEdit: (entry: JournalEntry) => void;
  onDelete: (entry: JournalEntry) => void;
}

export function JournalEntryTable({ entries, onEdit, onDelete }: JournalEntryTableProps) {
  if (entries.length === 0) {
    return (
      <div className="rounded-lg border border-dashed border-neutral-200 p-12 text-center text-neutral-500 dark:border-neutral-800 dark:text-neutral-400">
        仕訳が登録されていません
      </div>
    );
  }

  return (
    <div className="rounded-md border border-neutral-200 bg-white dark:border-neutral-800 dark:bg-neutral-950">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>日付</TableHead>
            <TableHead>摘要</TableHead>
            <TableHead className="text-right">金額</TableHead>
            <TableHead className="text-right">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {entries.map((entry) => {
            const totalDebit = entry.lines.reduce((sum, l) => sum + l.debit_amount, 0);
            return (
              <TableRow key={entry.id}>
                <TableCell>{entry.date}</TableCell>
                <TableCell className="font-medium">{entry.description}</TableCell>
                <TableCell className="text-right font-mono">
                  {totalDebit.toLocaleString()}
                </TableCell>
                <TableCell className="space-x-2 text-right">
                  <Button variant="outline" size="sm" onClick={() => onEdit(entry)}>
                    編集
                  </Button>
                  <Button variant="destructive" size="sm" onClick={() => onDelete(entry)}>
                    削除
                  </Button>
                </TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
}
