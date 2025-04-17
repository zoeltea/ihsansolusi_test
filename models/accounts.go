package models

import "time"

type Account struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	NIK        string    `json:"nik"`
	NoHP       string    `json:"no_hp"`
	NoRekening string    `json:"no_rekening"`
	Saldo      float64   `json:"saldo"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateAccountRequest struct {
	Name string `json:"name" validate:"required"`
	NIK  string `json:"nik" validate:"required"`
	NoHP string `json:"no_hp" validate:"required"`
}

type SaldoResponse struct {
	NoRekening string  `json:"no_rekening"`
	Saldo      float64 `json:"saldo"`
}

type TransactionRequest struct {
	NoRekening string  `json:"no_rekening" validate:"required"`
	Nominal    float64 `json:"nominal" validate:"required,gt=0"`
	Reference  string  `json:"reference"`
}
