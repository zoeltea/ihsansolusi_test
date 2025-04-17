package models

import "accounts-service/utils"

var (
	AccountWithNIKIsExist         = "ACCOUNT_WITH_NIK_IS_EXIST"
	AccountWithNoHpKIsExist       = "ACCOUNT_WITH_NO_HP_IS_EXIST"
	AccountWithNoRekeningNotFound = "ACCOUNT_WITH_NO_REK_NOT_FOUND"
	AccountNameEmpty              = "ACCOUNT_NAME_EMPTY"
	AccountNikEmpty               = "ACCOUNT_NIK_EMPTY"
	AccountNoHpEmpty              = "ACCOUNT_NO_HP_EMPTY"
	AccountParamNoRekeningEmpty   = "ACCOUNT_PARAM_NO_REKENING_EMPTY"
	AccountParamNominalLessZero   = "ACCOUNT_PARAM_NOMINAL_LESS_THAN_ZERO"
	Accountinsufficient           = "ACCOUNT_INSUFFICIENT_SALDO"
	AccountInvalidRequest         = "ACCOUNT_CREATE_INVALID_REQUEST"
	CreditInvalidRequest          = "CREDIT_INVALID_REQUEST"
	DebitInvalidRequest           = "DEBIT_INVALID_REQUEST"
	GetAccountError               = "GET_ACCOUNT_ERROR"
	UpdateSaldoError              = "UPDATE_SALDO_ERROR"
	CreateMutationError           = "CREATE_MUTATION_ERROR"
	CreateAccountError            = "CREATE_ACCOUNT_ERROR"
	CreateTransactionDBError      = "CREATE_TRANSACTION_DB_ERROR"
	CommitTransactionDBError      = "COMMIT_TRANSACTION_DB_ERROR"

	AccountWithNIKIsExistErr         = utils.NewRemark("Account with NIK is already exist", AccountWithNIKIsExist, "nik", nil)
	AccountWithNoHpKIsExistErr       = utils.NewRemark("Account with No HP is already exist", AccountWithNoHpKIsExist, "name", nil)
	AccountWithNoRekeningNotFoundErr = utils.NewRemark("Account with No Rekening not found", AccountWithNoRekeningNotFound, "no_rekening", nil)
	AccountParamNoRekeningEmptyErr   = utils.NewRemark("Param No rekening empty", AccountParamNoRekeningEmpty, "no_rekening", nil)
	AccountParamNominalErr           = utils.NewRemark("Param nominal less than 0", AccountParamNominalLessZero, "nominal", nil)
	AccountNameEmptyErr              = utils.NewRemark("Parameter Account name is empty", AccountNameEmpty, "name", nil)
	AccountNikEmptyErr               = utils.NewRemark("Parameter Account NIK is empty", AccountNikEmpty, "nik", nil)
	AccountNoHpEmptyErr              = utils.NewRemark("Parameter Account No Hp is empty", AccountNoHpEmpty, "no hp", nil)
	AccountinsufficientErr           = utils.NewRemark("Saldo not enough / Insufficient balance", AccountNoHpEmpty, "nominal", nil)
	AccountInvalidRequestErr         = utils.NewRemark("Invalid parameter create account", AccountInvalidRequest, "name, nik, no_hp", nil)
	CreditInvalidRequestErr          = utils.NewRemark("Invalid parameter credit/tabung", CreditInvalidRequest, "no_rekening, nominal", nil)
	DebitInvalidRequestErr           = utils.NewRemark("Invalid parameter debit/tarik", DebitInvalidRequest, "no_rekening, nominal", nil)
)
