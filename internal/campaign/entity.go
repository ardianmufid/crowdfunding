package campaign

import (
	"crowdfunding/internal/user"
	"crowdfunding/internal/utils"
	"fmt"
	"time"
)

type Campaign struct {
	ID               int              `db:"id"`
	UserID           int              `db:"user_id"`
	Name             string           `db:"name"`
	ShortDescription string           `db:"short_description"`
	Description      string           `db:"description"`
	GoalAmount       int              `db:"goal_amount"`
	CurrentAmount    int              `db:"current_amount"`
	Perks            string           `db:"perks"`
	BeckerCount      int              `db:"becker_count"`
	Slug             string           `db:"slug"`
	CreatedAt        time.Time        `db:"created_at"`
	UpdatedAt        time.Time        `db:"updated_at"`
	CampaignImages   *[]CampaignImage `db:"campaign_images" json:"campaign_images"`
	User             user.User        `db:"user"`
}

type CampaignImage struct {
	ID         *int       `db:"id"`
	CampaignId *int       `db:"campaign_id"`
	FileName   *string    `db:"file_name"`
	IsPrimary  *bool      `db:"is_primary"`
	CreatedAt  *time.Time `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

func NewFromCreateCampaignRequest(request CreateCampaignRequest) Campaign {

	// slug
	slugString := fmt.Sprintf("%s %d", request.Name, request.User.Id)

	return Campaign{
		UserID:           request.User.Id,
		Name:             request.Name,
		ShortDescription: request.ShortDescription,
		Description:      request.Description,
		GoalAmount:       request.GoalAmount,
		Perks:            request.Perks,
		Slug:             utils.NewSlug(slugString),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		CampaignImages:   &[]CampaignImage{},
		User:             request.User,
	}
}
