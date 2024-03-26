package models

import (
	"encoding/json"
	"time"
)

type SimpleAccountData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OrderData struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    json.Number    `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}

type BalanceData struct {
	Current   json.Number `json:"current"`
	Withdrawn json.Number `json:"withdrawn"`
}

type WithdrawInputData struct {
	Order string `json:"order"`
	Sum   json.Number    `json:"sum"`
}

type WithdrawData struct {
	Order       string    `json:"order"`
	Sum         json.Number       `json:"sum"`
	ProceededAt time.Time `json:"proceeded_at"`
}

type AccrualData struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual json.Number    `json:"accrual"`
}
