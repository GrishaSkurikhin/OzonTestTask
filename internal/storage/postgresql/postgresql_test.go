package postgresql_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/storage/postgresql"
)

func TestSaveURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	storage := &postgresql.OrderStorage{db}

	testCases := []struct {
		name      string
		longURL   string
		shortURL  string
		wantError bool
	}{
		{
			name:      "Successfully",
			longURL:   "https://example.com",
			shortURL:  "localhost:8080/exmpl",
			wantError: false,
		},
		{
			name:      "Failed",
			longURL:   "https://example.com",
			shortURL:  "localhost:8080/exmpl",
			wantError: true,
		},
	}

	query := regexp.QuoteMeta(fmt.Sprintf(`
		INSERT INTO %s (long_url, short_url)
		VALUES ($1, $2)
	`, postgresql.UrlsTable))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectPrepare(query)
			if tc.wantError {
				mock.ExpectExec(query).WithArgs(tc.longURL, tc.shortURL).WillReturnError(err)
				mock.ExpectRollback()
			} else {
				mock.ExpectExec(query).WithArgs(tc.longURL, tc.shortURL).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			}

			err = storage.SaveURL(context.Background(), tc.longURL, tc.shortURL)
			if (err != nil) != tc.wantError {
				t.Errorf("SaveURL() error = %v, wantError %v", err, tc.wantError)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestIsShortURLExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	storage := &postgresql.OrderStorage{db}

	testCases := []struct {
		name      string
		shortURL  string
		wantError bool
	}{
		{
			name:      "Successfully",
			shortURL:  "localhost:8080/exmpl",
			wantError: false,
		},
		{
			name:      "Failed",
			shortURL:  "localhost:8080/exmpl",
			wantError: true,
		},
	}

	query := regexp.QuoteMeta(fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM %s
			WHERE short_url = $1
		)
	`, postgresql.UrlsTable))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectPrepare(query)
			if tc.wantError {
				mock.ExpectQuery(query).WithArgs(tc.shortURL).WillReturnError(err)
			} else {
				mock.ExpectQuery(query).WithArgs(tc.shortURL).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			}

			_, err = storage.IsShortURLExists(context.Background(), tc.shortURL)
			if (err != nil) != tc.wantError {
				t.Errorf("IsShortURLExists() error = %v, wantError %v", err, tc.wantError)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %v", err)
			}
		})
	}

}

func TestGetURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	storage := &postgresql.OrderStorage{db}

	testCases := []struct {
		name      string
		shortURL  string
		wantError bool
	}{
		{
			name:      "Successfully",
			shortURL:  "localhost:8080/exmpl",
			wantError: false,
		},
		{
			name:      "Failed",
			shortURL:  "localhost:8080/exmpl",
			wantError: true,
		},
	}

	query := regexp.QuoteMeta(fmt.Sprintf(`
		SELECT long_url
		FROM %s
		WHERE short_url = $1
	`, postgresql.UrlsTable))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectPrepare(query)
			if tc.wantError {
				mock.ExpectQuery(query).WithArgs(tc.shortURL).WillReturnError(err)
			} else {
				mock.ExpectQuery(query).WithArgs(tc.shortURL).WillReturnRows(sqlmock.NewRows([]string{"long_url"}).AddRow("https://example.com"))
			}

			_, err = storage.GetURL(context.Background(), tc.shortURL)
			if (err != nil) != tc.wantError {
				t.Errorf("GetURL() error = %v, wantError %v", err, tc.wantError)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %v", err)
			}
		})
	}
}
