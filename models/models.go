package models

type TransferRequest struct {
	RecipientID string `json:"id" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
	Valute      string `json:"valute" binding:"required"`
}