package repositories_test

import (
	"accounts-service/repositories"
	"accounts-service/utils"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAccountnRepository_Account(t *testing.T) {

	t.Run("repo not null", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		logger := utils.NewLogger("info")
		repo := repositories.NewAccountRepository(db, logger)
		assert.NotNil(t, repo)
	})

	t.Run("test begin transaction", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		logger := utils.NewLogger("info")
		repo := repositories.NewAccountRepository(db, logger)

		// Create a transaction
		mock.ExpectBegin()
		assert.NoError(t, err)

		// Execute test
		tx, err := repo.BeginTx(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, tx)
	})
}
