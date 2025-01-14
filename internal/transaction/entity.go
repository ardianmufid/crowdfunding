package transaction

import (
	"crowdfunding/internal/campaign"
	"crowdfunding/internal/user"
	"time"
)

type Transaction struct {
	ID         int               `db:"id"`
	CampaignID int               `db:"campaign_id"`
	UserID     int               `db:"user_id"`
	Amount     int               `db:"amount"`
	Status     string            `db:"status"`
	Code       *string           `db:"code"`
	PaymentURL *string           `db:"payment_url"`
	User       user.User         `db:"user"`
	Campaign   campaign.Campaign `db:"campaign"`
	CreatedAt  time.Time         `db:"created_at"`
	UpdatedAt  time.Time         `db:"updated_at"`
}

func NewFromCreateTransactionRequest(request CreateTransactionRequest) Transaction {

	code := "TRX-000111"

	return Transaction{
		CampaignID: request.CampaignID,
		UserID:     request.User.Id,
		Amount:     request.Amount,
		Status:     "pending",
		Code:       &code,
		User:       request.User,
		Campaign:   campaign.Campaign{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
