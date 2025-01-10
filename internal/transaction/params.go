package transaction

import (
	"crowdfunding/internal/user"
	"time"
)

type GetCampaignTrasactionRequest struct {
	ID   int `uri:"id"`
	User user.User
}

type CampaignTransactionsResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMapperTransactionResponse(transaction Transaction) CampaignTransactionsResponse {
	return CampaignTransactionsResponse{
		ID:        transaction.ID,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}
}

func NewMapperTransactionsResponse(transactions []Transaction) []CampaignTransactionsResponse {
	if len(transactions) == 0 {
		return []CampaignTransactionsResponse{}
	}

	var transactionsMapper []CampaignTransactionsResponse

	for _, transaction := range transactions {
		mapper := NewMapperTransactionResponse(transaction)
		transactionsMapper = append(transactionsMapper, mapper)
	}

	return transactionsMapper
}

// type UserTransactionResponse struct {
// 	ID        int              `json:"id"`
// 	Amount    int              `json:"amount"`
// 	Status    string           `json:"status"`
// 	CreatedAt time.Time        `json:"created_at"`
// 	Campaign  CampaignResponse `json:"campaign"`
// }

// type CampaignResponse struct {
// 	Name     string `json:"name"`
// 	ImageURL string `json:"image_url"`
// }

// func NewMapperUserTransactionResponse(transaction Transaction) UserTransactionResponse {
// 	formatter := UserTransactionResponse{}
// 	formatter.ID = transaction.ID
// 	formatter.Amount = transaction.Amount
// 	formatter.Status = transaction.Status
// 	formatter.CreatedAt = transaction.CreatedAt

// 	campaignResponse := CampaignResponse{}
// 	campaignResponse.Name = transaction.Campaign.Name
// 	campaignResponse.ImageURL = ""

// 	if len(*transaction.Campaign.CampaignImages) > 0 {
// 		campaignResponse.ImageURL = transaction.Campaign.CampaignImages[0].FileName
// 	}

// 	formatter.Campaign = campaignResponse

// 	return formatter
// }

// func NewMapperUserTransactionsResponse(transactions []Transaction) []UserTransactionResponse {
// 	if len(transactions) == 0 {
// 		return []UserTransactionResponse{}
// 	}

// 	var transactionsFormatter []UserTransactionResponse

// 	for _, transaction := range transactions {
// 		formatter := NewMapperUserTransactionResponse(transaction)
// 		transactionsFormatter = append(transactionsFormatter, formatter)
// 	}

// 	return transactionsFormatter
// }
