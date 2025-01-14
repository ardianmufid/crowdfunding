package transaction

import (
	"crowdfunding/internal/user"
	"time"
)

type GetCampaignTrasactionRequest struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}

type CampaignTransactionsResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMapperCampaignTransactionResponse(transaction Transaction) CampaignTransactionsResponse {
	return CampaignTransactionsResponse{
		ID:        transaction.ID,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}
}

func NewMapperCampaignTransactionsResponse(transactions []Transaction) []CampaignTransactionsResponse {
	if len(transactions) == 0 {
		return []CampaignTransactionsResponse{}
	}

	var transactionsMapper []CampaignTransactionsResponse

	for _, transaction := range transactions {
		mapper := NewMapperCampaignTransactionResponse(transaction)
		transactionsMapper = append(transactionsMapper, mapper)
	}

	return transactionsMapper
}

type UserTransactionResponse struct {
	ID        int              `json:"id"`
	Amount    int              `json:"amount"`
	Status    string           `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	Campaign  CampaignResponse `json:"campaign"`
}

type CampaignResponse struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func NewMapperUserTransactionResponse(transaction Transaction) UserTransactionResponse {
	formatter := UserTransactionResponse{}
	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	// Image URL
	ImageUrl := ""
	if transaction.Campaign.CampaignImages != nil && len(*transaction.Campaign.CampaignImages) > 0 {
		// Cari gambar yang is_primary = true
		for _, img := range *transaction.Campaign.CampaignImages {
			if *img.IsPrimary {
				ImageUrl = *img.FileName
				break // Hentikan loop setelah menemukan gambar utama
			}
		}

		// Jika tidak ditemukan gambar dengan is_primary = true, pilih gambar pertama
		if ImageUrl == "" {
			ImageUrl = *(*transaction.Campaign.CampaignImages)[0].FileName
		}
	}

	campaignResponse := CampaignResponse{}
	campaignResponse.Name = transaction.Campaign.Name
	campaignResponse.ImageURL = ImageUrl

	// if len(*transaction.Campaign.CampaignImages) > 0 {
	// 	campaignResponse.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	// }

	formatter.Campaign = campaignResponse

	return formatter
}

func NewMapperUserTransactionsResponse(transactions []Transaction) []UserTransactionResponse {
	if len(transactions) == 0 {
		return []UserTransactionResponse{}
	}

	var transactionsFormatter []UserTransactionResponse

	for _, transaction := range transactions {
		formatter := NewMapperUserTransactionResponse(transaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}

	return transactionsFormatter
}

type CreateTransactionRequest struct {
	Amount     int `json:"amount" binding:"required"`
	CampaignID int `json:"campaign_id" binding:"required"`
	User       user.User
}

type TransactionResponse struct {
	ID         int    `json:"id"`
	CampaignID int    `json:"campaign_id"`
	UserID     int    `json:"user_id"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

func NewMapperTransactionResponse(transaction Transaction) TransactionResponse {

	return TransactionResponse{
		ID:         transaction.ID,
		CampaignID: transaction.CampaignID,
		UserID:     transaction.UserID,
		Amount:     transaction.Amount,
		Status:     transaction.Status,
		Code:       *transaction.Code,
		PaymentURL: *transaction.PaymentURL,
	}
}
