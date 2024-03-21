package models

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
