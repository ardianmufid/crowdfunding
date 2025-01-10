package transaction

import (
	"crowdfunding/internal/campaign"
	"crowdfunding/internal/user"
	"time"
)

type Transaction struct {
	ID          int       `db:"id"`
	Campaign_ID int       `db:"campaign_id"`
	User_ID     int       `db:"user_id"`
	Amount      int       `db:"amount"`
	Status      string    `db:"status"`
	Code        *string   `db:"code"`
	User        user.User `db:"user"`
	Campaign    campaign.Campaign
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
