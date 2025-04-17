package repositories_test

import (
	"accounts-service/models"
	"accounts-service/repositories"
	"accounts-service/utils"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMutationRepository_CreateMutation(t *testing.T) {

	t.Run("repo not null", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		logger := utils.NewLogger("info")
		repo := repositories.NewMutationRepository(db, logger)
		assert.NotNil(t, repo)
	})

	t.Run("success with transaction", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		logger := utils.NewLogger("info")
		repo := repositories.NewMutationRepository(db, logger)

		// Create a transaction
		mock.ExpectBegin()
		tx, err := db.Begin()
		assert.NoError(t, err)

		// Test data
		mutation := &models.Mutation{
			AccountID: 1,
			Nominal:   10000,
			Type:      "credit/tabung",
			Reference: "",
		}

		// Mock expectation
		mock.ExpectQuery(`
			INSERT INTO mutations \(account_id, nominal, type, reference\)
			VALUES \(\$1, \$2, \$3, \$4\)
			RETURNING id, created_at
		`).
			WithArgs(mutation.AccountID, mutation.Nominal, mutation.Type, mutation.Reference).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))

		// Execute test
		err = repo.CreateMutation(context.Background(), tx, mutation)

		// Assertions result
		assert.NoError(t, err)
		assert.Equal(t, uint(1), mutation.ID)
		assert.NotZero(t, mutation.CreatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success without transaction", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		logger := utils.NewLogger("info")
		repo := repositories.NewMutationRepository(db, logger)

		// Test data
		mutation := &models.Mutation{
			AccountID: 1,
			Nominal:   10000,
			Type:      "debit/tarik",
			Reference: "withdrawal",
		}

		// Mock expectation
		mock.ExpectQuery(`
			INSERT INTO mutations \(account_id, nominal, type, reference\)
			VALUES \(\$1, \$2, \$3, \$4\)
			RETURNING id, created_at
		`).
			WithArgs(mutation.AccountID, mutation.Nominal, mutation.Type, mutation.Reference).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(2, time.Now()))

		// Execute
		err = repo.CreateMutation(context.Background(), nil, mutation)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, uint(2), mutation.ID)
		assert.NotZero(t, mutation.CreatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error create mutation", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		logger := utils.NewLogger("info")
		repo := repositories.NewMutationRepository(db, logger)

		// Test data
		mutation := &models.Mutation{
			AccountID: 1,
			Nominal:   10000,
			Type:      "credit/tabung",
			Reference: "",
		}

		// Mock expectation
		mock.ExpectQuery(`
			INSERT INTO mutations \(account_id, nominal, type, reference\)
			VALUES \(\$1, \$2, \$3, \$4\)
			RETURNING id, created_at
		`).
			WithArgs(mutation.AccountID, mutation.Nominal, mutation.Type, mutation.Reference).
			WillReturnError(errors.New("database error"))

		// Execute
		err = repo.CreateMutation(context.Background(), nil, mutation)

		// Assertions
		assert.Error(t, err)
		assert.EqualError(t, err, "error creating mutation: database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
