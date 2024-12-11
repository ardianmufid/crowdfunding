package campaign

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}

func (r repository) FindAllCampaign() ([]Campaign, error) {

	var campaigns []Campaign
	if err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error; err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r repository) FindCampaignByUserID(userID int) ([]Campaign, error) {

	var campaigns []Campaign
	if err := r.db.Where("id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error; err != nil {
		return campaigns, err
	}

	return campaigns, nil
}
