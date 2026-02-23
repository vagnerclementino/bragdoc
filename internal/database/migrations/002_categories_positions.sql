-- Migration 002: categories and positions support
-- Note: SQLite requires table recreation for structural changes

PRAGMA foreign_keys=OFF;

-- 1. Categories table (system-managed for now)
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL CHECK(length(name) <= 20 AND name = UPPER(name)),
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    UNIQUE(name)
);

INSERT INTO categories (name, description)
SELECT 'UNKNOWN', 'GENERAL CATEGORY'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'UNKNOWN');
INSERT INTO categories (name, description)
SELECT 'PROJECT', 'PROJECT DELIVERABLES'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'PROJECT');
INSERT INTO categories (name, description)
SELECT 'ACHIEVEMENT', 'MEASURABLE ACHIEVEMENTS'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'ACHIEVEMENT');
INSERT INTO categories (name, description)
SELECT 'SKILL', 'SKILLS AND LEARNING'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'SKILL');
INSERT INTO categories (name, description)
SELECT 'LEADERSHIP', 'TEAM OR LEADERSHIP ACTS'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'LEADERSHIP');
INSERT INTO categories (name, description)
SELECT 'INNOVATION', 'INNOVATIONS AND IMPROVEMENTS'
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = 'INNOVATION');

-- 2. Positions table (historical roles per user)
CREATE TABLE IF NOT EXISTS positions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    company TEXT,
    start_date DATE,
    end_date DATE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 3. Rebuild brags table with category/position references
ALTER TABLE brags RENAME TO brags_old;

CREATE TABLE brags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category_id INTEGER NOT NULL,
    position_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (position_id) REFERENCES positions(id)
);

INSERT INTO brags (id, owner_id, title, description, category_id, created_at, updated_at)
SELECT 
    b.id,
    b.owner_id,
    b.title,
    b.description,
    CASE b.category
        WHEN 1 THEN (SELECT id FROM categories WHERE name = 'UNKNOWN')
        WHEN 2 THEN (SELECT id FROM categories WHERE name = 'PROJECT')
        WHEN 3 THEN (SELECT id FROM categories WHERE name = 'ACHIEVEMENT')
        WHEN 4 THEN (SELECT id FROM categories WHERE name = 'SKILL')
        WHEN 5 THEN (SELECT id FROM categories WHERE name = 'LEADERSHIP')
        WHEN 6 THEN (SELECT id FROM categories WHERE name = 'INNOVATION')
        ELSE (SELECT id FROM categories WHERE name = 'UNKNOWN')
    END AS category_id,
    b.created_at,
    b.updated_at
FROM brags_old b;

DROP TABLE brags_old;

CREATE INDEX IF NOT EXISTS idx_brags_owner_id ON brags(owner_id);
CREATE INDEX IF NOT EXISTS idx_brags_category_id ON brags(category_id);
CREATE INDEX IF NOT EXISTS idx_brags_position_id ON brags(position_id);

PRAGMA foreign_keys=ON;
