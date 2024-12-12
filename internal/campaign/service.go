package campaign

import (
	"crowdfunding/internal/utils"
	"fmt"
)

type Repository interface {
	FindAllCampaign() ([]Campaign, error)
	FindCampaignByUserID(userID int) ([]Campaign, error)
	FindCampaignByID(ID int) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) service {
	return service{
		repo: repo,
	}
}

func (s service) GetAllCampaign(userID int) ([]Campaign, error) {
	if userID != 0 {
		campaigns, err := s.repo.FindCampaignByUserID(userID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	campaigns, err := s.repo.FindAllCampaign()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s service) GetCampaignByID(request CampaignDetailRequest) (Campaign, error) {
	campaign, err := s.repo.FindCampaignByID(request.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s service) CreateCampaign(request CreateCampaignRequest) (Campaign, error) {

	campaign := Campaign{}
	campaign.Name = request.Name
	campaign.ShortDescription = request.ShortDescription
	campaign.Description = request.Description
	campaign.Perks = request.Perks
	campaign.GoalAmount = request.GoalAmount
	campaign.UserID = request.User.Id

	// slug
	slugString := fmt.Sprintf("%s %d", request.Name, request.User.Id)
	campaign.Slug = utils.NewSlug(slugString)

	newCampaign, err := s.repo.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}
