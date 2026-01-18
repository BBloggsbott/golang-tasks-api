package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHealthCheck(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Errorf("failed to create mock: %v", err)
	}

	defer db.Close()

	t.Run("successful health check", func(t *testing.T) {
		mock.ExpectPing()

		err := HealthCheck(db)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("unsuccessful health check", func(t *testing.T) {
		mock.ExpectPing().WillReturnError(sqlmock.ErrCancelled)

		err := HealthCheck(db)
		if err == nil {
			t.Error("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

}

func TestClose(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("failed to create mock: %v", err)
	}

	mock.ExpectClose()

	err = Close(db)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

}
