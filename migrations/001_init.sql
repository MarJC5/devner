CREATE TABLE IF NOT EXISTS projects (
    name        TEXT PRIMARY KEY,
    type        TEXT NOT NULL,              -- wordpress | laravel | node | nextjs | astro
    db_engine   TEXT NOT NULL DEFAULT '',   -- mysql | postgres | ''
    db_name     TEXT NOT NULL DEFAULT '',
    domain      TEXT NOT NULL,              -- e.g. myapp.localhost
    path        TEXT NOT NULL,              -- absolute path on host
    created_at  INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS hosts (
    domain      TEXT PRIMARY KEY,
    target      TEXT NOT NULL,
    managed_at  INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS history (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    ts          INTEGER NOT NULL,
    tool        TEXT NOT NULL,
    args_json   TEXT NOT NULL,
    result_json TEXT NOT NULL,
    duration_ms INTEGER NOT NULL,
    ok          INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_history_ts ON history(ts);
CREATE INDEX IF NOT EXISTS idx_history_tool ON history(tool);
