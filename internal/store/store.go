package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lukeyeager/school-lunch-schedule/internal/healthepro"
	_ "modernc.org/sqlite" // register sqlite driver
)

const schema = `
CREATE TABLE IF NOT EXISTS menus (
    date       TEXT PRIMARY KEY,
    fetched_at TEXT NOT NULL,
    source     TEXT NOT NULL,
    items      TEXT NOT NULL
)`

// MenuRecord is a row from the menus table.
type MenuRecord struct {
	Date      string
	FetchedAt time.Time
	Source    string
	Items     []healthepro.DisplayItem
}

// Store wraps a SQLite database for menu persistence.
type Store struct {
	db *sql.DB
}

// New opens (or creates) the SQLite database at path and applies the schema.
// Existing databases with extra columns (e.g. the old `changed` column) are
// left intact; the new code simply ignores those columns.
func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("applying schema: %w", err)
	}
	return &Store{db: db}, nil
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// Upsert inserts or replaces the menu record for the given date.
func (s *Store) Upsert(date string, entry *healthepro.DayEntry) error {
	items, err := json.Marshal(entry.Items)
	if err != nil {
		return fmt.Errorf("marshaling items: %w", err)
	}
	_, err = s.db.Exec(`
		INSERT INTO menus (date, fetched_at, source, items)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(date) DO UPDATE SET
			fetched_at = excluded.fetched_at,
			source     = excluded.source,
			items      = excluded.items`,
		date,
		time.Now().UTC().Format(time.RFC3339),
		entry.Source,
		string(items),
	)
	return err
}

// Get retrieves the menu record for the given ISO date, or nil if not found.
func (s *Store) Get(date string) (*MenuRecord, error) {
	row := s.db.QueryRow(
		`SELECT date, fetched_at, source, items FROM menus WHERE date = ?`, date)

	var r MenuRecord
	var fetchedAt, itemsJSON string

	err := row.Scan(&r.Date, &fetchedAt, &r.Source, &itemsJSON)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scanning menu row: %w", err)
	}

	r.FetchedAt, err = time.Parse(time.RFC3339, fetchedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing fetched_at: %w", err)
	}
	if err := json.Unmarshal([]byte(itemsJSON), &r.Items); err != nil {
		return nil, fmt.Errorf("parsing items JSON: %w", err)
	}
	return &r, nil
}
