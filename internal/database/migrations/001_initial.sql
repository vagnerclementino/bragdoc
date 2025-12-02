-- Initial schema for Bragdoc
-- Creates users, brags, tags, and brag_tags tables

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    job_title TEXT,
    company TEXT,
    locale TEXT NOT NULL DEFAULT 'en-US', -- Locale format: language-COUNTRY (e.g., en-US, pt-BR, pt-PT)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);

CREATE TABLE brags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL CHECK(length(name) >= 2 AND length(name) <= 20), -- Tag name: 2-20 characters
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, owner_id),
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE brag_tags (
    brag_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (brag_id, tag_id),
    FOREIGN KEY (brag_id) REFERENCES brags(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Indexes for optimizing searches
CREATE INDEX idx_brag_tags_brag_id ON brag_tags(brag_id);
CREATE INDEX idx_brag_tags_tag_id ON brag_tags(tag_id);
CREATE INDEX idx_tags_owner_id ON tags(owner_id);
CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_brags_owner_id ON brags(owner_id);
CREATE INDEX idx_brags_category ON brags(category);
