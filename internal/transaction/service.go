package transaction

import (
	"context"
	"crowdfunding/internal/campaign"
	"errors"
)

type Repository interface {
	GetCampaignID(ctx context.Context, CampaignID int) (transactions []Transaction, err error)
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

func NewService(repo Repository, campaignRepo campaign.Repository) service {
	return service{
		repository:         repo,
		campaignRepository: campaignRepo,
	}
}

func (s service) GetTransactionByCampaignID(ctx context.Context, request GetCampaignTrasactionRequest) (transactions []Transaction, err error) {

	campaign, err := s.campaignRepository.FindCampaignByID(ctx, request.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != request.User.Id {
		return []Transaction{}, errors.New("failed to update campaign")
	}

	transactions, err = s.repository.GetCampaignID(ctx, request.ID)
	if err != nil {
		return []Transaction{}, err
	}

	return
}
