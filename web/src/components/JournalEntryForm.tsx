import { useState } from 'react';
import type { Account } from '@/types/account';
import type { JournalEntry, CreateJournalEntryRequest } from '@/types/journalEntry';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@/components/ui/select';
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from '@/components/ui/table';

interface JournalEntryFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  accounts: Account[];
  entry?: JournalEntry;
  onSubmit: (data: CreateJournalEntryRequest) => void;
}

interface LineInput {
  accountId: string;
  debitAmount: string;
  creditAmount: string;
}

function emptyLine(): LineInput {
  return { accountId: '', debitAmount: '', creditAmount: '' };
}

export function JournalEntryForm({
  open,
  onOpenChange,
  accounts,
  entry,
  onSubmit,
}: JournalEntryFormProps) {
  const isEdit = !!entry;
  const isCompoundEntry =
    entry && entry.lines.length > 2
      ? true
      : entry && entry.lines.length === 2
        ? !(
            (entry.lines[0].debit_amount > 0 && entry.lines[1].credit_amount > 0) ||
            (entry.lines[0].credit_amount > 0 && entry.lines[1].debit_amount > 0)
          )
        : false;

  const [date, setDate] = useState(entry?.date ?? new Date().toISOString().slice(0, 10));
  const [description, setDescription] = useState(entry?.description ?? '');
  const [receiptRequired, setReceiptRequired] = useState(entry?.receipt_required ?? false);
  const [memo, setMemo] = useState(entry?.memo ?? '');
  const [isCompound, setIsCompound] = useState(isCompoundEntry);

  // Single mode state
  const [debitAccountId, setDebitAccountId] = useState<string>(
    entry && !isCompoundEntry
      ? String(entry.lines.find((l) => l.debit_amount > 0)?.account_id ?? '')
      : '',
  );
  const [creditAccountId, setCreditAccountId] = useState<string>(
    entry && !isCompoundEntry
      ? String(entry.lines.find((l) => l.credit_amount > 0)?.account_id ?? '')
      : '',
  );
  const [amount, setAmount] = useState<string>(
    entry && !isCompoundEntry
      ? String(entry.lines.find((l) => l.debit_amount > 0)?.debit_amount ?? '')
      : '',
  );

  // Compound mode state
  const [lines, setLines] = useState<LineInput[]>(() => {
    if (entry && isCompoundEntry) {
      return entry.lines.map((l) => ({
        accountId: String(l.account_id),
        debitAmount: l.debit_amount > 0 ? String(l.debit_amount) : '',
        creditAmount: l.credit_amount > 0 ? String(l.credit_amount) : '',
      }));
    }
    return [emptyLine(), emptyLine(), emptyLine()];
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const totalDebit = lines.reduce((sum, l) => sum + (parseInt(l.debitAmount) || 0), 0);
  const totalCredit = lines.reduce((sum, l) => sum + (parseInt(l.creditAmount) || 0), 0);
  const isBalanced = totalDebit === totalCredit && totalDebit > 0;

  const switchToCompound = () => {
    const newLines: LineInput[] = [];
    if (debitAccountId && amount) {
      newLines.push({ accountId: debitAccountId, debitAmount: amount, creditAmount: '' });
    } else {
      newLines.push(emptyLine());
    }
    if (creditAccountId && amount) {
      newLines.push({ accountId: creditAccountId, debitAmount: '', creditAmount: amount });
    } else {
      newLines.push(emptyLine());
    }
    newLines.push(emptyLine());
    setLines(newLines);
    setIsCompound(true);
  };

  const switchToSingle = () => {
    const debitLine = lines.find((l) => parseInt(l.debitAmount) > 0);
    const creditLine = lines.find((l) => parseInt(l.creditAmount) > 0);
    if (debitLine) {
      setDebitAccountId(debitLine.accountId);
      setAmount(debitLine.debitAmount);
    }
    if (creditLine) {
      setCreditAccountId(creditLine.accountId);
    }
    setIsCompound(false);
  };

  const updateLine = (index: number, field: keyof LineInput, value: string) => {
    setLines((prev) => prev.map((l, i) => (i === index ? { ...l, [field]: value } : l)));
  };

  const addLine = () => setLines((prev) => [...prev, emptyLine()]);
  const removeLine = (index: number) => setLines((prev) => prev.filter((_, i) => i !== index));

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const newErrors: Record<string, string> = {};

    if (!date) newErrors.date = '日付は必須です';
    if (!description.trim()) newErrors.description = '摘要は必須です';

    if (!isCompound) {
      if (!debitAccountId) newErrors.debitAccountId = '借方科目は必須です';
      if (!creditAccountId) newErrors.creditAccountId = '貸方科目は必須です';
      if (debitAccountId && creditAccountId && debitAccountId === creditAccountId)
        newErrors.creditAccountId = '借方と貸方に同じ科目は指定できません';
      const parsedAmount = parseInt(amount);
      if (!amount || isNaN(parsedAmount) || parsedAmount <= 0)
        newErrors.amount = '金額は1以上の整数を入力してください';
    } else {
      if (!isBalanced) newErrors.balance = '借方合計と貸方合計が一致していません';
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    let requestLines: CreateJournalEntryRequest['lines'];
    if (isCompound) {
      requestLines = lines
        .filter((l) => l.accountId && (parseInt(l.debitAmount) > 0 || parseInt(l.creditAmount) > 0))
        .map((l) => ({
          account_id: parseInt(l.accountId),
          debit_amount: parseInt(l.debitAmount) || 0,
          credit_amount: parseInt(l.creditAmount) || 0,
        }));
    } else {
      const parsedAmount = parseInt(amount);
      requestLines = [
        { account_id: parseInt(debitAccountId), debit_amount: parsedAmount, credit_amount: 0 },
        { account_id: parseInt(creditAccountId), debit_amount: 0, credit_amount: parsedAmount },
      ];
    }

    onSubmit({
      date,
      description: description.trim(),
      receipt_required: receiptRequired,
      memo,
      lines: requestLines,
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? '仕訳を編集' : '仕訳を作成'}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4 py-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label htmlFor="je-date">日付</Label>
              <Input
                id="je-date"
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
              />
              {errors.date && <p className="text-destructive text-xs">{errors.date}</p>}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="je-description">摘要</Label>
              <Input
                id="je-description"
                type="text"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="例: 携帯電話料金6月分"
              />
              {errors.description && (
                <p className="text-destructive text-xs">{errors.description}</p>
              )}
            </div>
          </div>

          {!isCompound ? (
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-1.5">
                  <Label htmlFor="je-debit">借方科目</Label>
                  <Select value={debitAccountId} onValueChange={setDebitAccountId}>
                    <SelectTrigger id="je-debit" className="w-full">
                      <SelectValue placeholder="借方科目を選択" />
                    </SelectTrigger>
                    <SelectContent>
                      {accounts.map((a) => (
                        <SelectItem key={a.id} value={String(a.id)}>
                          {a.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {errors.debitAccountId && (
                    <p className="text-destructive text-xs">{errors.debitAccountId}</p>
                  )}
                </div>
                <div className="space-y-1.5">
                  <Label htmlFor="je-credit">貸方科目</Label>
                  <Select value={creditAccountId} onValueChange={setCreditAccountId}>
                    <SelectTrigger id="je-credit" className="w-full">
                      <SelectValue placeholder="貸方科目を選択" />
                    </SelectTrigger>
                    <SelectContent>
                      {accounts.map((a) => (
                        <SelectItem key={a.id} value={String(a.id)}>
                          {a.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {errors.creditAccountId && (
                    <p className="text-destructive text-xs">{errors.creditAccountId}</p>
                  )}
                </div>
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="je-amount">金額</Label>
                <Input
                  id="je-amount"
                  type="number"
                  min="1"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  placeholder="例: 10000"
                />
                {errors.amount && <p className="text-destructive text-xs">{errors.amount}</p>}
              </div>
            </div>
          ) : (
            <div className="space-y-2">
              <div className="rounded-md border border-neutral-200 dark:border-neutral-800">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>科目</TableHead>
                      <TableHead className="w-32">借方</TableHead>
                      <TableHead className="w-32">貸方</TableHead>
                      <TableHead className="w-12" />
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {lines.map((line, i) => (
                      <TableRow key={i}>
                        <TableCell>
                          <Select
                            value={line.accountId}
                            onValueChange={(v) => updateLine(i, 'accountId', v)}
                          >
                            <SelectTrigger className="w-full">
                              <SelectValue placeholder="科目を選択" />
                            </SelectTrigger>
                            <SelectContent>
                              {accounts.map((a) => (
                                <SelectItem key={a.id} value={String(a.id)}>
                                  {a.name}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        </TableCell>
                        <TableCell>
                          <Input
                            type="number"
                            min="0"
                            value={line.debitAmount}
                            onChange={(e) => updateLine(i, 'debitAmount', e.target.value)}
                            placeholder="0"
                          />
                        </TableCell>
                        <TableCell>
                          <Input
                            type="number"
                            min="0"
                            value={line.creditAmount}
                            onChange={(e) => updateLine(i, 'creditAmount', e.target.value)}
                            placeholder="0"
                          />
                        </TableCell>
                        <TableCell>
                          {lines.length > 2 && (
                            <Button
                              type="button"
                              variant="destructive"
                              size="sm"
                              onClick={() => removeLine(i)}
                            >
                              ×
                            </Button>
                          )}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
              <div className="flex items-center justify-between">
                <Button type="button" variant="outline" size="sm" onClick={addLine}>
                  + 行追加
                </Button>
                <div className="text-sm">
                  借方合計: <span className="font-mono font-bold">{totalDebit.toLocaleString()}</span>
                  {' / '}
                  貸方合計: <span className="font-mono font-bold">{totalCredit.toLocaleString()}</span>
                  {totalDebit > 0 || totalCredit > 0 ? (
                    isBalanced ? (
                      <span className="ml-2 text-green-600">✓ 一致</span>
                    ) : (
                      <span className="ml-2 text-red-600">✗ 不一致</span>
                    )
                  ) : null}
                </div>
              </div>
              {errors.balance && (
                <p className="text-destructive text-xs">{errors.balance}</p>
              )}
            </div>
          )}

          <div className="flex items-center justify-between">
            <Button type="button" variant="ghost" size="sm" onClick={isCompound ? switchToSingle : switchToCompound}>
              {isCompound ? '単一仕訳に切替' : '複合仕訳に切替'}
            </Button>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="flex items-center space-x-2">
              <Checkbox
                id="je-receipt"
                checked={receiptRequired}
                onCheckedChange={(checked) => setReceiptRequired(checked === true)}
              />
              <Label htmlFor="je-receipt">領収書の要否</Label>
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="je-memo">メモ</Label>
            <Input
              id="je-memo"
              type="text"
              value={memo}
              onChange={(e) => setMemo(e.target.value)}
              placeholder="自由入力"
            />
          </div>

          <DialogFooter className="pt-4">
            <Button type="submit" className="w-full sm:w-auto">
              {isEdit ? '更新' : '作成'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
