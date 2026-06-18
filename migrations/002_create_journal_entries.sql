CREATE TABLE IF NOT EXISTS journal_entries (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    date             TEXT    NOT NULL,
    description      TEXT    NOT NULL,
    receipt_required INTEGER NOT NULL DEFAULT 0,
    memo             TEXT,
    created_at       TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at       TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS journal_lines (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    journal_entry_id INTEGER NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    account_id       INTEGER NOT NULL REFERENCES accounts(id),
    debit_amount     INTEGER NOT NULL DEFAULT 0,
    credit_amount    INTEGER NOT NULL DEFAULT 0,
    created_at       TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at       TEXT    NOT NULL DEFAULT (datetime('now')),
    CHECK (debit_amount >= 0),
    CHECK (credit_amount >= 0),
    CHECK (debit_amount > 0 OR credit_amount > 0),
    CHECK (NOT (debit_amount > 0 AND credit_amount > 0))
);
