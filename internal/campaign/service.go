package campaign

type Repository interface {
	FindAllCampaign() ([]Campaign, error)
	FindCampaignByUserID(userID int) ([]Campaign, error)
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
