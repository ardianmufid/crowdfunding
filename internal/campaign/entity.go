package campaign

import (
	"time"
)

type Campaign struct {
	Id               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	GoalAmount       int
	CurrentAmount    int
	Perks            string
	BeckerCount      int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage
}

type CampaignImage struct {
	Id         int
	CampaignId int
	FileName   string
	IsPrimary  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
