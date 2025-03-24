-- Table: Content (for notes/snippets/bookmarks/docs)
CREATE TABLE IF NOT EXISTS content (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,               -- note, snippet, bookmark, doc
    title TEXT,
    body TEXT,
    source_url TEXT,                  -- for bookmarks/docs
    file_path TEXT,                   -- for local docs
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Table: Embeddings
CREATE TABLE IF NOT EXISTS embeddings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id INTEGER NOT NULL,
    embedding BLOB NOT NULL,          -- store as binary or array of floats
    model TEXT NOT NULL,              -- which embedding model was used
    dimensions INTEGER NOT NULL,      -- number of dimensions in the embedding
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE
);

-- Optional: Tags table (for organization)
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS content_tags (
    content_id INTEGER,
    tag_id INTEGER,
    PRIMARY KEY (content_id, tag_id),
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Create indexes for faster querying
CREATE INDEX IF NOT EXISTS idx_content_type ON content(type);
CREATE INDEX IF NOT EXISTS idx_content_title ON content(title);
CREATE INDEX IF NOT EXISTS idx_embeddings_content_id ON embeddings(content_id);
