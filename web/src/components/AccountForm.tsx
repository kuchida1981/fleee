import { useState } from 'react';
import type { Account, AccountType } from '@/types/account';
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
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@/components/ui/select';

interface AccountFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  account?: Account;
  onSubmit: (data: {
    name: string;
    account_type: AccountType;
    display_order: number;
  }) => void;
}

export function AccountForm({
  open,
  onOpenChange,
  account,
  onSubmit,
}: AccountFormProps) {
  const [name, setName] = useState(account?.name ?? '');
  const [accountType, setAccountType] = useState<AccountType | ''>(account?.account_type ?? '');
  const [displayOrder, setDisplayOrder] = useState<number>(account?.display_order ?? 0);
  const [errors, setErrors] = useState<{ name?: string; accountType?: string }>({});

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const newErrors: { name?: string; accountType?: string } = {};

    if (!name.trim()) {
      newErrors.name = '科目名は必須です';
    }
    if (!accountType) {
      newErrors.accountType = '勘定区分は必須です';
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    onSubmit({
      name: name.trim(),
      account_type: accountType as AccountType,
      display_order: displayOrder,
    });
  };

  const isEdit = !!account;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{isEdit ? '勘定科目を編集' : '勘定科目を作成'}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4 py-4">
          <div className="space-y-1.5">
            <Label htmlFor="name">科目名</Label>
            <Input
              id="name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="例: 普通預金"
            />
            {errors.name && (
              <p className="text-xs text-destructive">{errors.name}</p>
            )}
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="account_type">勘定区分</Label>
            <Select
              value={accountType}
              onValueChange={(val) => setAccountType(val as AccountType)}
            >
              <SelectTrigger id="account_type" className="w-full">
                <SelectValue placeholder="勘定区分を選択してください" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="asset">資産</SelectItem>
                <SelectItem value="liability">負債</SelectItem>
                <SelectItem value="equity">純資産</SelectItem>
                <SelectItem value="revenue">収益</SelectItem>
                <SelectItem value="expense">費用</SelectItem>
              </SelectContent>
            </Select>
            {errors.accountType && (
              <p className="text-xs text-destructive">{errors.accountType}</p>
            )}
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="display_order">表示順序</Label>
            <Input
              id="display_order"
              type="number"
              value={displayOrder}
              onChange={(e) => setDisplayOrder(Number(e.target.value))}
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
