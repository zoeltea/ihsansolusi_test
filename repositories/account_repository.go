package repositories

import (
	"accounts-service/models"
	"accounts-service/utils"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	BeginTx(ctx context.Context) (*sql.Tx, error) // Add this new method
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
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}
	return tx, nil
}

func (r *accountRepository) CreateAccount(ctx context.Context, account *models.Account) error {
	queryInsert := `
		INSERT INTO accounts (name, nik, no_hp, no_rekening, saldo)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	queryUpdateNoRek := `
		UPDATE accounts
		SET no_rekening = $1
		WHERE id = $2
	`

	err := r.db.QueryRowContext(ctx, queryInsert,
		account.Name,
		account.NIK,
		account.NoHP,
		account.NoRekening,
		account.Saldo,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		r.logger.Error("Error creating account: %v", err)
		return fmt.Errorf("error creating account: %w", err)
	}

	noRekening := generateNoRek(account.ID)

	_, err = r.db.ExecContext(ctx, queryUpdateNoRek, noRekening, account.ID)

	if err != nil {
		r.logger.Error("Error update account for no rekening: %v", err)
		return fmt.Errorf("error update account for no rekening: %w", err)
	}

	account.NoRekening = noRekening

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
		return nil, fmt.Errorf("error getting account by no rekening: %w", err)
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
		return nil, fmt.Errorf("error getting account by NIK: %w", err)
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
		return nil, fmt.Errorf("error getting account by no hp: %w", err)
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
		r.logger.Error("Error updating account saldo: %v", err)
		return fmt.Errorf("error updating account saldo: %w", err)
	}

	return nil
}

func generateNoRek(id uint) string {

	lengthNoRekEnv := os.Getenv("LENGTH_NO_REKENING")
	lengthNorekInt, err := strconv.Atoi(lengthNoRekEnv)
	if err != nil || lengthNoRekEnv == "" {
		lengthNorekInt = 12
	}
	// Convert the id to string
	idString := strconv.FormatUint(uint64(id), 10)

	// Calculate the format lenght
	characterNeeded := lengthNorekInt - len(idString)
	noRekening := strings.Repeat("0", characterNeeded)
	noRekening = noRekening + idString

	return noRekening
}
