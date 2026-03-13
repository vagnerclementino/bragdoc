-- Initial schema for Bragdoc
-- Complete schema with users, brags, tags, categories, and job_titles

-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    job_title TEXT,
    company TEXT,
    locale TEXT NOT NULL DEFAULT 'en-US',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);

-- Categories table (system-managed)
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL CHECK(length(name) <= 20 AND name = UPPER(name)),
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    UNIQUE(name)
);

-- Insert default categories
INSERT INTO categories (name, description) VALUES
    ('UNKNOWN', 'GENERAL CATEGORY'),
    ('PROJECT', 'PROJECT DELIVERABLES'),
    ('ACHIEVEMENT', 'MEASURABLE ACHIEVEMENTS'),
    ('SKILL', 'SKILLS AND LEARNING'),
    ('LEADERSHIP', 'TEAM OR LEADERSHIP ACTS'),
    ('INNOVATION', 'INNOVATIONS AND IMPROVEMENTS');

-- Job titles table (historical roles per user)
CREATE TABLE job_titles (
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

-- Unique constraint: only one active job title per user
CREATE UNIQUE INDEX idx_job_titles_user_active 
ON job_titles(user_id) 
WHERE end_date IS NULL;

-- Brags table
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

-- Tags table
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL CHECK(length(name) >= 2 AND length(name) <= 20),
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, owner_id),
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

-- Brag tags junction table
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
CREATE INDEX idx_brags_category_id ON brags(category_id);
CREATE INDEX idx_brags_job_title_id ON brags(job_title_id);
