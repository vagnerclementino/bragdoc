-- Migration 003: Rename positions to job_titles
-- This migration renames the positions table to job_titles for better clarity

PRAGMA foreign_keys=OFF;

-- 1. Rename positions table to job_titles
ALTER TABLE positions RENAME TO job_titles;

-- 2. Rebuild brags table to update foreign key reference
ALTER TABLE brags RENAME TO brags_old;

CREATE TABLE brags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category_id INTEGER NOT NULL,
    job_title_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (job_title_id) REFERENCES job_titles(id)
);

INSERT INTO brags (id, owner_id, title, description, category_id, job_title_id, created_at, updated_at)
SELECT id, owner_id, title, description, category_id, position_id, created_at, updated_at
FROM brags_old;

DROP TABLE brags_old;

CREATE INDEX IF NOT EXISTS idx_brags_owner_id ON brags(owner_id);
CREATE INDEX IF NOT EXISTS idx_brags_category_id ON brags(category_id);
CREATE INDEX IF NOT EXISTS idx_brags_job_title_id ON brags(job_title_id);

PRAGMA foreign_keys=ON;
