package repositories

import (
	"accounts-service/models"
	"accounts-service/utils"
	"context"
	"database/sql"
	"fmt"
)

type MutationRepository interface {
	CreateMutation(ctx context.Context, tx *sql.Tx, mutation *models.Mutation) error
}

type mutationRepository struct {
	db     *sql.DB
	logger utils.Logger
}

func NewMutationRepository(db *sql.DB, logger utils.Logger) MutationRepository {
	return &mutationRepository{
		db:     db,
		logger: logger,
	}
}

func (r *mutationRepository) CreateMutation(ctx context.Context, tx *sql.Tx, mutation *models.Mutation) error {
	query := `
		INSERT INTO mutations (account_id, nominal, type, reference)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query,
			mutation.AccountID,
			mutation.Nominal,
			mutation.Type,
			mutation.Reference,
		).Scan(&mutation.ID, &mutation.CreatedAt)
	} else {
		err = r.db.QueryRowContext(ctx, query,
			mutation.AccountID,
			mutation.Nominal,
			mutation.Type,
			mutation.Reference,
		).Scan(&mutation.ID, &mutation.CreatedAt)
	}

	if err != nil {
		r.logger.Error("Error creating mutation: %v", err)
		return fmt.Errorf("error creating mutation: %w", err)
	}

	return nil
}
