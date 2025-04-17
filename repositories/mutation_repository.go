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
	GetMutationsByAccountID(ctx context.Context, accountID int) ([]models.Mutation, error)
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

func (r *mutationRepository) GetMutationsByAccountID(ctx context.Context, accountID int) ([]models.Mutation, error) {
	query := `
		SELECT id, account_id, nominal, type, reference, created_at
		FROM mutations
		WHERE account_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		r.logger.Error("Error getting mutations by account ID: %v", err)
		return nil, fmt.Errorf("error getting mutations by account ID: %w", err)
	}
	defer rows.Close()

	var mutations []models.Mutation
	for rows.Next() {
		var mutation models.Mutation
		err := rows.Scan(
			&mutation.ID,
			&mutation.AccountID,
			&mutation.Nominal,
			&mutation.Type,
			&mutation.Reference,
			&mutation.CreatedAt,
		)
		if err != nil {
			r.logger.Error("Error scanning mutation row: %v", err)
			return nil, fmt.Errorf("error scanning mutation row: %w", err)
		}
		mutations = append(mutations, mutation)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("Error after scanning mutation rows: %v", err)
		return nil, fmt.Errorf("error after scanning mutation rows: %w", err)
	}

	return mutations, nil
}
