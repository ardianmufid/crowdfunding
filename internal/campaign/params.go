package campaign

import (
	"crowdfunding/internal/user"
	"strings"
)

type CreateCampaignRequest struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	Perks            string `json:"perks" binding:"required"`
	GoalAmount       int    `json:"goal_amount" binding:"required"`
	User             user.User
}

type CreateCampaignImageRequest struct {
	CampaignID int  `form:"campaign_id" binding:"required"`
	IsPrimary  bool `form:"is_primary"`
	User       user.User
}

type CampaignDetailRequest struct {
	ID int `uri:"id"`
}

type CampaignsResponse struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

type CampaignDetailResponse struct {
	ID               int                     `json:"id"`
	Name             string                  `json:"name"`
	ShortDescription string                  `json:"short_description"`
	Description      string                  `json:"description"`
	ImageURL         string                  `json:"image_url"`
	GoalAmount       int                     `json:"goal_amount"`
	CurrentAmount    int                     `json:"current_amount"`
	UserID           int                     `json:"user_id"`
	Slug             string                  `json:"slug"`
	Perks            []string                `json:"perks"`
	User             CampaignUserResponse    `json:"user"`
	Images           []CampaignImageResponse `json:"images"`
}

type CampaignUserResponse struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignImageResponse struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func NewMapperCampaignResponse(campaign Campaign) CampaignsResponse {

	// Image URL
	ImageUrl := ""
	if campaign.CampaignImages != nil && len(*campaign.CampaignImages) > 0 {
		// Cari gambar yang is_primary = true
		for _, img := range *campaign.CampaignImages {
			if *img.IsPrimary {
				ImageUrl = *img.FileName
				break // Hentikan loop setelah menemukan gambar utama
			}
		}

		// Jika tidak ditemukan gambar dengan is_primary = true, pilih gambar pertama
		if ImageUrl == "" {
			ImageUrl = *(*campaign.CampaignImages)[0].FileName
		}
	}

	return CampaignsResponse{
		ID:               campaign.ID,
		UserID:           campaign.UserID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Slug:             campaign.Slug,
		ImageURL:         ImageUrl,
	}
}

func NewMapperCampaignsResponse(campaigns []Campaign) []CampaignsResponse {

	campaignsFormatter := []CampaignsResponse{}

	for _, campaign := range campaigns {
		campaignFormatter := NewMapperCampaignResponse(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter
}

func NewMapperCampaignDetailResponse(campaign Campaign) CampaignDetailResponse {

	// Image URL
	// Image URL
	ImageUrl := ""
	if campaign.CampaignImages != nil && len(*campaign.CampaignImages) > 0 {
		// Cari gambar yang is_primary = true
		for _, img := range *campaign.CampaignImages {
			if *img.IsPrimary {
				ImageUrl = *img.FileName
				break // Hentikan loop setelah menemukan gambar utama
			}
		}

		// Jika tidak ditemukan gambar dengan is_primary = true, pilih gambar pertama
		if ImageUrl == "" {
			ImageUrl = *(*campaign.CampaignImages)[0].FileName
		}
	}

	// Perks
	var perks []string

	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, strings.TrimSpace(perk))
	}

	// User
	user := campaign.User

	campaignUserFormatter := CampaignUserResponse{}
	campaignUserFormatter.Name = user.Name
	campaignUserFormatter.ImageURL = user.AvatarFileName

	// Images
	campaignImagesFormatter := []CampaignImageResponse{}

	for _, image := range *campaign.CampaignImages {
		campaignImageFormatter := CampaignImageResponse{}
		campaignImageFormatter.ImageURL = *image.FileName
		campaignImageFormatter.IsPrimary = *image.IsPrimary

		campaignImagesFormatter = append(campaignImagesFormatter, campaignImageFormatter)
	}

	return CampaignDetailResponse{
		ID:               campaign.ID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		Description:      campaign.Description,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Slug:             campaign.Slug,
		ImageURL:         ImageUrl,
		UserID:           campaign.UserID,
		Perks:            perks,
		User:             campaignUserFormatter,
		Images:           campaignImagesFormatter,
	}

}
