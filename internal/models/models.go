package models

import "time"

type SimpleAccountData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OrderData struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    string `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}

type BalanceData struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

type WithdrawInputData struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

type WithdrawData struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProceededAt time.Time `json:"proceeded_at"`
}

type AccrualData struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual string `json:"accrual"`
}
