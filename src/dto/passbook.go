package dto

type PassbookRequest struct {
	Description string `json:"description"`
	Amount      int32  `json:"amount"`
}
