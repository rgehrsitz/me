package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaFS embed.FS

// DB is the database connection
type DB struct {
	*sql.DB
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
	// Ensure directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set pragmas for better performance
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to set foreign_keys pragma: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("failed to set journal_mode pragma: %w", err)
	}

	// Initialize schema
	schemaSQL, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to read schema: %w", err)
	}

	if _, err := db.Exec(string(schemaSQL)); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("Database initialized successfully")
	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// GetContent retrieves a content item by ID
func (db *DB) GetContent(id int64) (*Content, error) {
	row := db.QueryRow(`
		SELECT id, type, title, body, source_url, file_path, created_at, updated_at 
		FROM content 
		WHERE id = ?`, id)

	var content Content
	err := row.Scan(
		&content.ID,
		&content.Type,
		&content.Title,
		&content.Body,
		&content.SourceURL,
		&content.FilePath,
		&content.CreatedAt,
		&content.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get tags for this content
	rows, err := db.Query(`
		SELECT t.name 
		FROM tags t 
		JOIN content_tags ct ON t.id = ct.tag_id 
		WHERE ct.content_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	content.Tags = []string{}
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		content.Tags = append(content.Tags, tag)
	}

	return &content, nil
}

// Content represents a piece of content in the PKB
type Content struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	SourceURL string `json:"source_url,omitempty"`
	FilePath  string `json:"file_path,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Tags      []string `json:"tags,omitempty"`
}

// CreateContent creates a new content item
func (db *DB) CreateContent(content *Content) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		INSERT INTO content (type, title, body, source_url, file_path) 
		VALUES (?, ?, ?, ?, ?)`,
		content.Type, content.Title, content.Body, content.SourceURL, content.FilePath)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Add tags if any
	if len(content.Tags) > 0 {
		for _, tag := range content.Tags {
			// Insert tag if it doesn't exist
			res, err := tx.Exec("INSERT OR IGNORE INTO tags (name) VALUES (?)", tag)
			if err != nil {
				return 0, err
			}

			// Get tag ID
			var tagID int64
			row := tx.QueryRow("SELECT id FROM tags WHERE name = ?", tag)
			if err := row.Scan(&tagID); err != nil {
				return 0, err
			}

			// Link tag to content
			_, err = tx.Exec("INSERT INTO content_tags (content_id, tag_id) VALUES (?, ?)", id, tagID)
			if err != nil {
				return 0, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

// ListContent retrieves a list of content items with optional filtering
func (db *DB) ListContent(contentType string, limit, offset int) ([]Content, error) {
	query := `
		SELECT id, type, title, body, source_url, file_path, created_at, updated_at 
		FROM content`
	args := []interface{}{}

	if contentType != "" {
		query += " WHERE type = ?"
		args = append(args, contentType)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contents := []Content{}
	for rows.Next() {
		var content Content
		err := rows.Scan(
			&content.ID,
			&content.Type,
			&content.Title,
			&content.Body,
			&content.SourceURL,
			&content.FilePath,
			&content.CreatedAt,
			&content.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get tags for this content
		tagRows, err := db.Query(`
			SELECT t.name 
			FROM tags t 
			JOIN content_tags ct ON t.id = ct.tag_id 
			WHERE ct.content_id = ?`, content.ID)
		if err != nil {
			return nil, err
		}
		defer tagRows.Close()

		content.Tags = []string{}
		for tagRows.Next() {
			var tag string
			if err := tagRows.Scan(&tag); err != nil {
				return nil, err
			}
			content.Tags = append(content.Tags, tag)
		}

		contents = append(contents, content)
	}

	return contents, nil
}

// UpdateContent updates an existing content item
func (db *DB) UpdateContent(content *Content) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE content 
		SET type = ?, title = ?, body = ?, source_url = ?, file_path = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`,
		content.Type, content.Title, content.Body, content.SourceURL, content.FilePath, content.ID)
	if err != nil {
		return err
	}

	// Delete existing tag links
	_, err = tx.Exec("DELETE FROM content_tags WHERE content_id = ?", content.ID)
	if err != nil {
		return err
	}

	// Add tags if any
	if len(content.Tags) > 0 {
		for _, tag := range content.Tags {
			// Insert tag if it doesn't exist
			_, err := tx.Exec("INSERT OR IGNORE INTO tags (name) VALUES (?)", tag)
			if err != nil {
				return err
			}

			// Get tag ID
			var tagID int64
			row := tx.QueryRow("SELECT id FROM tags WHERE name = ?", tag)
			if err := row.Scan(&tagID); err != nil {
				return err
			}

			// Link tag to content
			_, err = tx.Exec("INSERT INTO content_tags (content_id, tag_id) VALUES (?, ?)", content.ID, tagID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// DeleteContent deletes a content item by ID
func (db *DB) DeleteContent(id int64) error {
	_, err := db.Exec("DELETE FROM content WHERE id = ?", id)
	return err
}

// StoreEmbedding stores an embedding for a content item
func (db *DB) StoreEmbedding(contentID int64, embedding []byte, model string, dimensions int) (int64, error) {
	// Check if embedding already exists for this content
	var embeddingID int64
	row := db.QueryRow("SELECT id FROM embeddings WHERE content_id = ? AND model = ?", contentID, model)
	err := row.Scan(&embeddingID)
	if err == nil {
		// Update existing embedding
		_, err = db.Exec(`
			UPDATE embeddings 
			SET embedding = ?, dimensions = ? 
			WHERE id = ?`,
			embedding, dimensions, embeddingID)
		return embeddingID, err
	}

	// Insert new embedding
	res, err := db.Exec(`
		INSERT INTO embeddings (content_id, embedding, model, dimensions) 
		VALUES (?, ?, ?, ?)`,
		contentID, embedding, model, dimensions)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
