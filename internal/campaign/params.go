package campaign

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

func NewMapperCampaignResponse(campaign Campaign) CampaignsResponse {

	ImageUrl := ""
	if len(campaign.CampaignImages) > 0 {
		ImageUrl = campaign.CampaignImages[0].FileName
	}
	return CampaignsResponse{
		ID:               campaign.GoalAmount,
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
