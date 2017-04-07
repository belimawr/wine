package store

import (
	"database/sql"
	"log"

	"time"

	"github.com/belimawr/wine/models"
)

// Store - Interface to abstract database operations in a high level
type Store interface {
	PutWine(wine models.Wine) error
}

// NewSQLiteStore - Returns a new Store that uses SQLite
func NewSQLiteStore(db *sql.DB) Store {
	return sqlite{
		db: db,
	}
}

type sqlite struct {
	db *sql.DB
}

func (s sqlite) PutWine(w models.Wine) error {
	query := `INSERT INTO wine (name, price, grape, description, crawled_at) VALUES (?, ?, ?, ?, ?)`

	now := time.Now()

	_, err := s.db.Exec(query, w.Name, w.Price, w.Grape, w.Description, now)

	if err != nil {
		log.Printf("Could not insert Wine into database, error: %q", err)
	}

	return nil
}
