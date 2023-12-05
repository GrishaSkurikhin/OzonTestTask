package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	urlsTable = "urls"

	QueryTimeout = 5 * time.Second
)

type OrderStorage struct {
	*sql.DB
}

func New(source string) (*OrderStorage, error) {
	const op = "storage.postgresql.New"

	conn, err := sql.Open("postgres", source)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: ping failed: %v", op, err)
	}
	return &OrderStorage{conn}, nil
}

func (s *OrderStorage) Disconnect(ctx context.Context) error {
	const op = "storage.postgresql.Disconnect"

	err := s.Close()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}

func (s *OrderStorage) SaveURL(longURL string, shortURL string) error {
	const op = "storage.postgresql.SaveURL"

	query := fmt.Sprintf(`
		INSERT INTO %s (long_url, short_url)
		VALUES ($1, $2)
	`, urlsTable)

	tx, err := s.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %v", op, err)
	}

	addReq, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %v", op, err)
	}
	defer addReq.Close()

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	_, err = addReq.ExecContext(ctx, longURL, shortURL)
	if err != nil {
		return fmt.Errorf("%s: failed to execute statement: %v", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %v", op, err)
	}

	return nil
}

func (s *OrderStorage) IsShortURLExists(shortURL string) (bool, error) {
	const op = "storage.postgresql.IsShortURLExists"

	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM %s
			WHERE short_url = $1
		)
	`, urlsTable)

	checkReq, err := s.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: failed to prepare statement: %v", op, err)
	}
	defer checkReq.Close()

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	var exists bool
	err = checkReq.QueryRowContext(ctx, shortURL).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: failed to execute statement: %v", op, err)
	}

	return exists, nil
}

func (s *OrderStorage) GetURL(shortURL string) (string, error) {
	const op = "storage.postgresql.GetURL"

	query := fmt.Sprintf(`
		SELECT long_url
		FROM %s
		WHERE short_url = $1
	`, urlsTable)

	getReq, err := s.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("%s: failed to prepare statement: %v", op, err)
	}
	defer getReq.Close()

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	var longURL string
	err = getReq.QueryRowContext(ctx, shortURL).Scan(&longURL)
	if err != nil {
		return "", fmt.Errorf("%s: failed to execute statement: %v", op, err)
	}

	return longURL, nil
}
