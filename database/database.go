// database/database.go
package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mhassaninmsft/gogitcrawler/models"
)

type DB struct {
	db *sql.DB
}

func Connect(connectionString string) (*DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (db *DB) SaveContributor(contributor *models.Contributor) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert contributor
	query := `INSERT INTO contributors (id, login) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET login = $2`
	if _, err := tx.Exec(query, contributor.ID, contributor.Login); err != nil {
		return err
	}

	// Insert repos
	for _, repo := range contributor.Repos {
		query := `INSERT INTO repos (id, name, full_name, contributor_id) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET name = $2, full_name = $3, contributor_id = $4`
		if _, err := tx.Exec(query, repo.ID, repo.Name, repo.FullName, contributor.ID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
