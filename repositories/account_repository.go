package repositories

import (
	"accounts-service/models"
	"accounts-service/utils"
	"context"
	"database/sql"
)

const (
	LENGTH_NO_REK = 12
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *models.Account) error
	GetAccountByNoRekening(ctx context.Context, noRekening string) (*models.Account, error)
	GetAccountByNoHp(ctx context.Context, noHp string) (*models.Account, error)
	GetAccountByNik(ctx context.Context, nik string) (*models.Account, error)
	UpdateSaldo(ctx context.Context, tx *sql.Tx, accountID uint, nominal float64) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type accountRepository struct {
	db     *sql.DB
	logger utils.Logger
}

func NewAccountRepository(db *sql.DB, logger utils.Logger) AccountRepository {
	return &accountRepository{
		db:     db,
		logger: logger,
	}
}

func (r *accountRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("Error beginning transaction: %v", err)
		return nil, utils.NewRemark(
			"Error beginning transaction",
			models.CreateTransactionDBError,
			"",
			err,
		)
	}
	return tx, nil
}

func (r *accountRepository) CreateAccount(ctx context.Context, account *models.Account) error {
	queryInsert := `
		INSERT INTO accounts (name, nik, no_hp, saldo, no_rekening)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, queryInsert,
		account.Name,
		account.NIK,
		account.NoHP,
		account.Saldo,
		account.NoRekening,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		r.logger.Error("Error creating account: %v", err)
		return utils.NewRemark(
			"Error creating account",
			models.CreateAccountError,
			"",
			err,
		)
	}

	return nil
}

func (r *accountRepository) GetAccountByNoRekening(ctx context.Context, no_rekening string) (*models.Account, error) {
	query := `
		SELECT id, name, nik, no_hp, no_rekening, saldo, created_at, updated_at
		FROM accounts
		WHERE no_rekening = $1
	`

	var account models.Account
	err := r.db.QueryRowContext(ctx, query, no_rekening).Scan(
		&account.ID,
		&account.Name,
		&account.NIK,
		&account.NoHP,
		&account.NoRekening,
		&account.Saldo,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Error getting account by no rekening: %v", err)
		return nil, utils.NewRemark(
			"Error getting account by no rekening",
			models.GetAccountError,
			"no_rekening",
			err,
		)
	}

	return &account, nil
}

func (r *accountRepository) GetAccountByNik(ctx context.Context, nik string) (*models.Account, error) {
	query := `
		SELECT id, name, nik, no_hp, no_rekening, saldo, created_at, updated_at
		FROM accounts
		WHERE nik = $1
	`

	var account models.Account
	err := r.db.QueryRowContext(ctx, query, nik).Scan(
		&account.ID,
		&account.Name,
		&account.NIK,
		&account.NoHP,
		&account.NoRekening,
		&account.Saldo,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Error getting account by NIK: %v", err)
		return nil, utils.NewRemark(
			"Error getting account by no nik",
			models.GetAccountError,
			"no_rekening",
			err,
		)
	}

	return &account, nil
}

func (r *accountRepository) GetAccountByNoHp(ctx context.Context, no_hp string) (*models.Account, error) {
	query := `
		SELECT id, name, nik, no_hp, no_rekening, saldo, created_at, updated_at
		FROM accounts
		WHERE no_hp = $1
	`

	var account models.Account
	err := r.db.QueryRowContext(ctx, query, no_hp).Scan(
		&account.ID,
		&account.Name,
		&account.NIK,
		&account.NoHP,
		&account.NoRekening,
		&account.Saldo,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Error getting account by no hp: %v", err)
		return nil, utils.NewRemark(
			"Error getting account by no no_hp",
			models.GetAccountError,
			"no_rekening",
			err,
		)
	}

	return &account, nil
}

func (r *accountRepository) UpdateSaldo(ctx context.Context, tx *sql.Tx, accountID uint, nominal float64) error {
	query := `
		UPDATE accounts
		SET saldo = saldo + $1, updated_at = NOW()
		WHERE id = $2
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, nominal, accountID)
	} else {
		_, err = r.db.ExecContext(ctx, query, nominal, accountID)
	}

	if err != nil {
		typeTransaction := "credit/tabung"
		if nominal < 0 {
			typeTransaction = "debit/tarik"
		}
		r.logger.Error("Error updating account saldo: %v", err)
		return utils.NewRemark(
			"error updating account saldo",
			models.UpdateSaldoError,
			"no_rekening",
			map[string]interface{}{
				"error": err,
				"type":  typeTransaction,
			},
		)
	}

	return nil
}
