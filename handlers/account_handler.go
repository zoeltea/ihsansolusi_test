package handlers

import (
	"accounts-service/models"
	"accounts-service/usecases"
	"accounts-service/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	accountUsecase usecases.AccountUsecase
	logger         utils.Logger
}

func NewAccountHandler(accountUsecase usecases.AccountUsecase, logger utils.Logger) *AccountHandler {
	return &AccountHandler{
		accountUsecase: accountUsecase,
		logger:         logger,
	}
}

func (h *AccountHandler) CreateAccount(ctx echo.Context) error {
	var req models.CreateAccountRequest

	if err := ctx.Bind(&req); err != nil {
		h.logger.Warning("Error param request: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"remark": "invalid request"})
	}

	if req.NIK == "" {
		return ctx.JSON(http.StatusBadRequest, models.AccountNikEmptyErr)
	}

	if req.Name == "" {
		return ctx.JSON(http.StatusBadRequest, models.AccountNameEmptyErr)
	}

	if req.NoHP == "" {
		return ctx.JSON(http.StatusBadRequest, models.AccountNoHpEmptyErr)
	}

	account, err := h.accountUsecase.CreateAccount(ctx.Request().Context(), &req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"remark": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) GetSaldo(ctx echo.Context) error {
	noRekening := ctx.Param("no_rekening")

	saldo, err := h.accountUsecase.GetSaldo(ctx.Request().Context(), noRekening)
	if err != nil {
		h.logger.Error("Error getting saldo: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"remark": err.Error()})
	}

	return ctx.JSON(http.StatusOK, saldo)
}

func (h *AccountHandler) Debit(ctx echo.Context) error {
	var req models.TransactionRequest
	if err := ctx.Bind(&req); err != nil {
		h.logger.Warning("Error param request: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"remark": "invalid request"})
	}

	if req.NoRekening == "" {
		return ctx.JSON(http.StatusBadRequest, models.AccountParamNoRekeningEmptyErr)
	}

	if req.Nominal <= 0 {
		return ctx.JSON(http.StatusBadRequest, models.AccountParamNominalErr)
	}

	err := h.accountUsecase.Debit(ctx.Request().Context(), &req)
	if err != nil {
		h.logger.Error("Error processing debit: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"remark": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "penarikan saldo successful"})
}

func (h *AccountHandler) Credit(ctx echo.Context) error {
	var req models.TransactionRequest
	if err := ctx.Bind(&req); err != nil {
		h.logger.Warning("Error binding credit/tabung request: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"remark": "invalid request"})
	}

	if req.NoRekening == "" {
		return ctx.JSON(http.StatusBadRequest, models.AccountParamNoRekeningEmptyErr)
	}

	if req.Nominal <= 0 {
		return ctx.JSON(http.StatusBadRequest, models.AccountParamNominalErr)
	}

	err := h.accountUsecase.Credit(ctx.Request().Context(), &req)
	if err != nil {
		h.logger.Error("Error processing credit/tabung: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"remark": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "menabung successful"})
}
