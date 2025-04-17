package usecases

import (
	"accounts-service/models"
	"accounts-service/repositories"
	"accounts-service/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type AccountUsecase interface {
	CreateAccount(ctx context.Context, req *models.CreateAccountRequest) (*models.Account, error)
	GetAccountByNoRekening(ctx context.Context, noRekening string) (*models.Account, error)
	GetSaldo(ctx context.Context, noRekening string) (*models.SaldoResponse, error)
	Debit(ctx context.Context, req *models.TransactionRequest) error
	Credit(ctx context.Context, req *models.TransactionRequest) error
}

type accountUsecase struct {
	accountRepo  repositories.AccountRepository
	mutationRepo repositories.MutationRepository
	logger       utils.Logger
}

func NewAccountUsecase(accountRepo repositories.AccountRepository, mutationRepo repositories.MutationRepository, logger utils.Logger) AccountUsecase {
	return &accountUsecase{
		accountRepo:  accountRepo,
		mutationRepo: mutationRepo,
		logger:       logger,
	}
}

func (u *accountUsecase) CreateAccount(ctx context.Context, req *models.CreateAccountRequest) (*models.Account, error) {
	// Check if account with same nik already exists
	existingAccount, err := u.accountRepo.GetAccountByNik(ctx, req.NIK)
	if err != nil {
		u.logger.Error("Error checking existing account: %v", err)
		return nil, fmt.Errorf("error checking existing account: %w", err)
	}

	if existingAccount != nil {
		return nil, models.AccountWithNIKIsExistErr
	}

	// Check if account with same no_hp already exists
	existingAccount, err = u.accountRepo.GetAccountByNoHp(ctx, req.NoHP)
	if err != nil {
		u.logger.Error("Error checking existing account: %v", err)
		return nil, err
	}

	if existingAccount != nil {
		return nil, models.AccountWithNoHpKIsExistErr
	}

	// generate no rekening from timestamp
	currentTime := time.Now()
	unixTime := currentTime.Unix()
	noRekening := strconv.FormatInt(unixTime, 10)

	account := &models.Account{
		Name:       req.Name,
		NIK:        req.NIK,
		NoHP:       req.NoHP,
		NoRekening: noRekening,
	}

	err = u.accountRepo.CreateAccount(ctx, account)
	if err != nil {
		u.logger.Error("Error creating account: %v", err)
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	return account, nil
}

func (u *accountUsecase) GetAccountByNoRekening(ctx context.Context, noRekening string) (*models.Account, error) {
	account, err := u.accountRepo.GetAccountByNoRekening(ctx, noRekening)
	if err != nil {
		u.logger.Error("Error getting account by no rekening: %v", err)
		return nil, fmt.Errorf("error getting account by no rekening: %w", err)
	}

	if account == nil {
		return nil, errors.New("account not found")
	}

	return account, nil
}

func (u *accountUsecase) GetSaldo(ctx context.Context, noRekening string) (*models.SaldoResponse, error) {
	account, err := u.accountRepo.GetAccountByNoRekening(ctx, noRekening)
	if err != nil {
		u.logger.Error("Error getting account saldo: %v", err)
		return nil, fmt.Errorf("error getting account saldo: %w", err)
	}

	if account == nil {
		return nil, errors.New("account not found")
	}

	return &models.SaldoResponse{
		NoRekening: account.NoRekening,
		Saldo:      account.Saldo,
	}, nil
}

func (u *accountUsecase) Debit(ctx context.Context, req *models.TransactionRequest) error {
	// Start transaction1
	tx, err := u.accountRepo.BeginTx(ctx)
	if err != nil {
		u.logger.Error("Error starting transaction: %v", err)
		return fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				u.logger.Error("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	// Get account
	account, err := u.accountRepo.GetAccountByNoRekening(ctx, req.NoRekening)
	if err != nil {
		u.logger.Error("Error getting account for debit/tarik: %v", err)
		return fmt.Errorf("error getting account for debit/tarik: %w", err)
	}

	if account == nil {
		return models.AccountWithNoRekeningNotFoundErr
	}

	// Check if saldo is enough
	if account.Saldo < req.Nominal {
		return models.AccountinsufficientErr
	}

	// Update saldo (debit/tarik)
	err = u.accountRepo.UpdateSaldo(ctx, tx, account.ID, -req.Nominal)
	if err != nil {
		u.logger.Error("Error updating saldo for debit/tarik: %v", err)
		return fmt.Errorf("error updating saldo for debit/tarik: %w", err)
	}

	// Create mutation record
	mutation := &models.Mutation{
		AccountID: account.ID,
		Nominal:   req.Nominal,
		Type:      "debit/tarik",
		Reference: req.Reference,
	}

	err = u.mutationRepo.CreateMutation(ctx, tx, mutation)
	if err != nil {
		u.logger.Error("Error creating mutation for debit/tarik: %v", err)
		return fmt.Errorf("error creating mutation for debit/tarik: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		u.logger.Error("Error committing transaction: %v", err)
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (u *accountUsecase) Credit(ctx context.Context, req *models.TransactionRequest) error {
	// Start transaction
	tx, err := u.accountRepo.BeginTx(ctx)
	if err != nil {
		u.logger.Error("Error starting transaction: %v", err)
		return fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				u.logger.Error("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	// Get account
	account, err := u.accountRepo.GetAccountByNoRekening(ctx, req.NoRekening)
	if err != nil {
		u.logger.Error("Error getting account for credit/tabung: %v", err)
		return fmt.Errorf("error getting account for credit/tabung: %w", err)
	}
	if account == nil {
		return models.AccountWithNoRekeningNotFoundErr
	}

	// Update saldo (credit/tabung)
	err = u.accountRepo.UpdateSaldo(ctx, tx, account.ID, req.Nominal)
	if err != nil {
		u.logger.Error("Error updating saldo for credit/tabung: %v", err)
		return fmt.Errorf("error updating saldo for credit/tabung: %w", err)
	}

	// Create mutation record
	mutation := &models.Mutation{
		AccountID: account.ID,
		Nominal:   req.Nominal,
		Type:      "credit/tabung",
		Reference: req.Reference,
	}

	err = u.mutationRepo.CreateMutation(ctx, tx, mutation)
	if err != nil {
		u.logger.Error("Error creating mutation for credit/tabung: %v", err)
		return fmt.Errorf("error creating mutation for credit/tabung: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		u.logger.Error("Error committing transaction: %v", err)
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
