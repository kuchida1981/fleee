import type {
  JournalEntry,
  CreateJournalEntryRequest,
  UpdateJournalEntryRequest,
} from '@/types/journalEntry';

export async function listJournalEntries(): Promise<JournalEntry[]> {
  const response = await fetch('/api/journal-entries');
  if (!response.ok) {
    throw new Error(`Failed to fetch journal entries: ${response.statusText}`);
  }
  return response.json();
}

export async function getJournalEntry(id: number): Promise<JournalEntry> {
  const response = await fetch(`/api/journal-entries/${id}`);
  if (!response.ok) {
    throw new Error(`Failed to get journal entry ${id}: ${response.statusText}`);
  }
  return response.json();
}

export async function createJournalEntry(data: CreateJournalEntryRequest): Promise<JournalEntry> {
  const response = await fetch('/api/journal-entries', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  if (!response.ok) {
    throw new Error(`Failed to create journal entry: ${response.statusText}`);
  }
  return response.json();
}

export async function updateJournalEntry(
  id: number,
  data: UpdateJournalEntryRequest,
): Promise<JournalEntry> {
  const response = await fetch(`/api/journal-entries/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  if (!response.ok) {
    throw new Error(`Failed to update journal entry ${id}: ${response.statusText}`);
  }
  return response.json();
}

export async function deleteJournalEntry(id: number): Promise<void> {
  const response = await fetch(`/api/journal-entries/${id}`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error(`Failed to delete journal entry ${id}: ${response.statusText}`);
  }
}
