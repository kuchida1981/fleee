import type { Account, ImportResult, AccountType } from '@/types/account';

export async function listAccounts(): Promise<Account[]> {
  const response = await fetch('/api/accounts');
  if (!response.ok) {
    throw new Error(`Failed to fetch accounts: ${response.statusText}`);
  }
  return response.json();
}

export async function createAccount(data: {
  name: string;
  account_type: AccountType;
  display_order: number;
}): Promise<Account> {
  const response = await fetch('/api/accounts', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  if (!response.ok) {
    throw new Error(`Failed to create account: ${response.statusText}`);
  }
  return response.json();
}

export async function getAccount(id: number): Promise<Account> {
  const response = await fetch(`/api/accounts/${id}`);
  if (!response.ok) {
    throw new Error(`Failed to get account ${id}: ${response.statusText}`);
  }
  return response.json();
}

export async function updateAccount(
  id: number,
  data: {
    name: string;
    account_type: AccountType;
    display_order: number;
  }
): Promise<Account> {
  const response = await fetch(`/api/accounts/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  if (!response.ok) {
    throw new Error(`Failed to update account ${id}: ${response.statusText}`);
  }
  return response.json();
}

export async function deleteAccount(id: number): Promise<void> {
  const response = await fetch(`/api/accounts/${id}`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error(`Failed to delete account ${id}: ${response.statusText}`);
  }
}

export async function importAccounts(file: File): Promise<ImportResult> {
  const formData = new FormData();
  formData.append('file', file);

  const response = await fetch('/api/accounts/import', {
    method: 'POST',
    body: formData,
  });
  if (!response.ok) {
    throw new Error(`Failed to import accounts: ${response.statusText}`);
  }
  return response.json();
}

