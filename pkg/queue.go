package pkg

import (
	"database/sql"
	"errors"
)

func Enqueue(db *sql.DB, message string) error {
	query := `INSERT INTO queue (message) VALUES ($1)`
	_, err := db.Exec(query, message)
	return err
}

func Dequeue(db *sql.DB) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	var id int
	var message string
	query := `SELECT id, message FROM queue WHERE processed = FALSE ORDER BY created_at LIMIT 1 FOR UPDATE`
	err = tx.QueryRow(query).Scan(&id, &message)
	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil {
			err = errors.Join(err, rErr)
		}
		return "", err
	}

	updateQuery := `UPDATE queue SET processed = TRUE WHERE id = $1`
	_, err = tx.Exec(updateQuery, id)
	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil {
			err = errors.Join(err, rErr)
		}
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return message, nil
}
