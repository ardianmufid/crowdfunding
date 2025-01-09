package campaign

import (
	"context"
	"errors"
	"log"
)

type Repository interface {
	FindAllCampaign(ctx context.Context) (campaigns []Campaign, err error)
	FindCampaignByUserID(ctx context.Context, userID int) (campaigns []Campaign, err error)
	FindCampaignByID(ctx context.Context, ID int) (campaign Campaign, err error)
	Save(ctx context.Context, model Campaign) (campaign Campaign, err error)
	Update(ctx context.Context, model Campaign) (campaign Campaign, err error)
	CreateImage(ctx context.Context, models CampaignImage) (campaignImage CampaignImage, err error)
	MarkAllImagesAsNonPrimary(ctx context.Context, campaignID int) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) service {
	return service{
		repo: repo,
	}
}

func (s service) CreateCampaign(ctx context.Context, request CreateCampaignRequest) (campaign Campaign, err error) {

	campaign = NewFromCreateCampaignRequest(request)

	newCampaign, err := s.repo.Save(ctx, campaign)
	if err != nil {
		log.Println("error service")
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s service) GetAllCampaign(ctx context.Context, userID int) (campaigns []Campaign, err error) {
	if userID != 0 {
		campaigns, err := s.repo.FindCampaignByUserID(ctx, userID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	campaigns, err = s.repo.FindAllCampaign(ctx)
	if err != nil {
		log.Println("error svc")
		return campaigns, err
	}

	return campaigns, nil
}

func (s service) GetCampaignByID(ctx context.Context, request CampaignDetailRequest) (campaign Campaign, err error) {

	campaign, err = s.repo.FindCampaignByID(ctx, request.ID)
	if err != nil {
		log.Printf("Scan error : %v", err)
		return
	}

	return
}

func (s service) UpdateCampaign(ctx context.Context, requestID CampaignDetailRequest, requestData CreateCampaignRequest) (campaign Campaign, err error) {

	campaign, err = s.repo.FindCampaignByID(ctx, requestID.ID)
	if err != nil {
		log.Printf("service s.repo.FindCampaignByID error : %v", err)
		return
	}

	if campaign.UserID != requestData.User.Id {
		log.Printf("service campaign.UserID != requestData.User.Id error : %v", err)
		return campaign, errors.New("failed to update campaign")
	}

	campaign.Name = requestData.Name
	campaign.ShortDescription = requestData.ShortDescription
	campaign.Description = requestData.Description
	campaign.Perks = requestData.Perks
	campaign.GoalAmount = requestData.GoalAmount

	updatedCampaign, err := s.repo.Update(ctx, campaign)
	if err != nil {
		log.Printf("service s.repo.Update(ctx, campaign) error : %v", err)
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

func (s service) SaveCampaignImage(ctx context.Context, request CreateCampaignImageRequest, fileLocation string) (campaignImage CampaignImage, err error) {
	campaign, err := s.repo.FindCampaignByID(ctx, request.CampaignID)
	if err != nil {
		return
	}

	if campaign.UserID != request.User.Id {
		return CampaignImage{}, errors.New("failed to created campaign image")
	}

	if request.IsPrimary {

		_, err := s.repo.MarkAllImagesAsNonPrimary(ctx, request.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	campaignImage.CampaignId = &request.CampaignID
	campaignImage.IsPrimary = &request.IsPrimary
	campaignImage.FileName = &fileLocation

	newCampaignImage, err := s.repo.CreateImage(ctx, campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}
