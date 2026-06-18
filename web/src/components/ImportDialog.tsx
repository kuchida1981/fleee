import { useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { importAccounts } from '@/api/accounts';
import type { ImportResult } from '@/types/account';

interface ImportDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onImportComplete: () => void;
}

export function ImportDialog({ open, onOpenChange, onImportComplete }: ImportDialogProps) {
  const [file, setFile] = useState<File | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [result, setResult] = useState<ImportResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      setFile(e.target.files[0]);
      setError(null);
      setResult(null);
    }
  };

  const handleImport = async () => {
    if (!file) {
      setError('ファイルを選択してください');
      return;
    }

    setIsSubmitting(true);
    setError(null);
    setResult(null);

    try {
      const importResult = await importAccounts(file);
      setResult(importResult);
      onImportComplete();
    } catch (err) {
      console.error(err);
      setError('インポートに失敗しました。ファイル形式やデータ内容を確認してください。');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>勘定科目のインポート</DialogTitle>
        </DialogHeader>

        <div className="space-y-4 py-4">
          <div className="space-y-1.5">
            <Input
              type="file"
              accept=".csv,.tsv"
              onChange={handleFileChange}
              disabled={isSubmitting}
            />
            <p className="text-xs text-neutral-500 dark:text-neutral-400">
              CSVまたはTSVファイルを選択してください。
            </p>
          </div>

          {error && <div className="text-destructive text-sm font-medium">{error}</div>}

          {result && (
            <div className="rounded-lg border border-green-200 bg-green-50 p-3 text-sm text-green-800 dark:border-green-900/50 dark:bg-green-950/30 dark:text-green-300">
              {result.total}件中 {result.success}件をインポートしました（{result.skipped}
              件スキップ）
            </div>
          )}
        </div>

        <DialogFooter>
          <Button
            onClick={handleImport}
            disabled={!file || isSubmitting}
            className="w-full sm:w-auto"
          >
            {isSubmitting ? 'インポート中...' : 'インポートを実行'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
