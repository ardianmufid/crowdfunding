package campaign

import (
	"context"
	"crowdfunding/internal/user"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository {
	return repository{
		db: db,
	}
}

func (r repository) FindAllCampaign(ctx context.Context) (campaigns []Campaign, err error) {
	query := `
        SELECT 
            c.id, c.user_id, c.name, c.short_description, c.description, 
            c.goal_amount, c.current_amount, c.perks, c.becker_count, 
            c.slug, c.created_at, c.updated_at,
            ci.id AS campaign_image_id, ci.campaign_id, ci.file_name, ci.is_primary, 
            ci.created_at AS campaign_image_created_at, ci.updated_at AS campaign_image_updated_at
        FROM 
            campaigns c
        LEFT JOIN 
            campaign_images ci 
        ON 
            c.id = ci.campaign_id AND ci.is_primary = true
    `

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		log.Printf("Scan error: %v", err)

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var campaign Campaign
		var campaignImage CampaignImage

		err := rows.Scan(
			&campaign.ID,
			&campaign.UserID,
			&campaign.Name,
			&campaign.ShortDescription,
			&campaign.Description,
			&campaign.GoalAmount,
			&campaign.CurrentAmount,
			&campaign.Perks,
			&campaign.BeckerCount,
			&campaign.Slug,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
			&campaignImage.ID,
			&campaignImage.CampaignId,
			&campaignImage.FileName,
			&campaignImage.IsPrimary,
			&campaignImage.CreatedAt,
			&campaignImage.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)

			return nil, err
		}

		// Jika CampaignImage.ID nil, kosongkan CampaignImages
		if campaignImage.ID != nil {
			campaign.CampaignImages = &[]CampaignImage{campaignImage}
		} else {
			campaign.CampaignImages = &[]CampaignImage{}
		}

		campaigns = append(campaigns, campaign)
	}
	return
}

func (r repository) Save(ctx context.Context, model Campaign) (campaign Campaign, err error) {

	query := `
        INSERT INTO campaigns (
            user_id, name, short_description, description, goal_amount, current_amount, perks, becker_count, slug, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()
        )
        RETURNING id, user_id, name, short_description, description, goal_amount, current_amount, perks, becker_count, slug, created_at, updated_at
    `

	err = r.db.QueryRowxContext(
		ctx,
		query,
		model.UserID,
		model.Name,
		model.ShortDescription,
		model.Description,
		model.GoalAmount,
		model.CurrentAmount,
		model.Perks,
		model.BeckerCount,
		model.Slug,
	).StructScan(&campaign)
	if err != nil {
		log.Printf("Scan error: %v", err)

		return
	}

	return
}

func (r repository) FindCampaignByUserID(ctx context.Context, userID int) (campaigns []Campaign, err error) {

	query := `
		SELECT
			c.id, c.user_id, c.name, c.short_description, c.description, c.goal_amount, c.current_amount, c.perks, c.becker_count, c.slug, c.created_at, c.updated_at, ci.id, ci.campaign_id, ci.file_name, ci.is_primary, ci.created_at, ci.updated_at
		FROM campaigns c
		LEFT JOIN campaign_images ci
		ON
			c.id = ci.campaign_id AND ci.is_primary = true
		WHERE c.user_id = $1
	`

	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		log.Printf("Scan error: %v", err)

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var campaign Campaign
		var campaignImage CampaignImage

		err := rows.Scan(
			&campaign.ID,
			&campaign.UserID,
			&campaign.Name,
			&campaign.ShortDescription,
			&campaign.Description,
			&campaign.GoalAmount,
			&campaign.CurrentAmount,
			&campaign.Perks,
			&campaign.BeckerCount,
			&campaign.Slug,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
			&campaignImage.ID,
			&campaignImage.CampaignId,
			&campaignImage.FileName,
			&campaignImage.IsPrimary,
			&campaignImage.CreatedAt,
			&campaignImage.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan error : %v", err)
			return nil, err
		}

		campaigns = append(campaigns, campaign)
	}

	return
}

func (r repository) FindCampaignByID(ctx context.Context, ID int) (Campaign, error) {
	var campaign Campaign

	// Query untuk mengambil data Campaign
	campaignQuery := `
		SELECT
			id, user_id, name, short_description, description,
			goal_amount, current_amount, perks, becker_count,
			slug, created_at, updated_at
		FROM campaigns
		WHERE id = $1
	`
	err := r.db.GetContext(
		ctx,
		&campaign,
		campaignQuery,
		ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return Campaign{}, fmt.Errorf("campaign with ID %d not found", ID)
		}
		return Campaign{}, err
	}

	// Query untuk mengambil CampaignImages yang terkait
	imagesQuery := `
		SELECT
			id, campaign_id, file_name, is_primary, created_at, updated_at
		FROM campaign_images
		WHERE campaign_id = $1
	`
	var campaignImages []CampaignImage
	err = r.db.SelectContext(ctx, &campaignImages, imagesQuery, ID)
	if err != nil {
		return Campaign{}, err
	}

	// Assign CampaignImages ke Campaign
	campaign.CampaignImages = &campaignImages

	// Query untuk mengambil User yang terkait
	userQuery := `
		SELECT
			id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user user.User
	err = r.db.GetContext(
		ctx,
		&user,
		userQuery,
		campaign.UserID,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return Campaign{}, err
		}
		// Jika user tidak ditemukan, biarkan nilai default kosong
	}
	campaign.User = user

	return campaign, nil
}

func (r repository) Update(ctx context.Context, model Campaign) (campaign Campaign, err error) {
	query := `
		UPDATE campaigns 
		SET 
			user_id = $1,
			name = $2,
			short_description = $3,
			description = $4,
			goal_amount = $5,
			current_amount = $6,
			perks = $7,
			becker_count = $8,
			slug = $9,
			created_at = $10,
			updated_at = $11
		WHERE id = $12
		RETURNING 
			id, user_id, name, short_description, description, 
			goal_amount, current_amount, perks, becker_count, 
			slug, created_at, updated_at
	`

	// Menggunakan NamedQueryContext untuk mendapatkan data yang diupdate
	err = r.db.QueryRowxContext(
		ctx,
		query,
		model.UserID,
		model.Name,
		model.ShortDescription,
		model.Description,
		model.GoalAmount,
		model.CurrentAmount,
		model.Perks,
		model.BeckerCount,
		model.Slug,
		model.CreatedAt,
		model.UpdatedAt,
		model.ID,
	).StructScan(&campaign)
	if err != nil {
		log.Printf("Scan error : %v", err)
		if err == sql.ErrNoRows {
			return
		}
		return
	}

	return
}

func (r repository) CreateImage(ctx context.Context, models CampaignImage) (campaignImage CampaignImage, err error) {

	query := `
		INSERT INTO campaign_images (
			campaign_id, file_name, is_primary, created_at, updated_at
		) VALUES (
			$1, $2, $3, NOW(), NOW()
		) RETURNING
		 	id, campaign_id, file_name, is_primary, created_at, updated_at
	`

	err = r.db.QueryRowxContext(ctx, query, models.CampaignId, models.FileName, models.IsPrimary).StructScan(&campaignImage)
	if err != nil {
		log.Printf("Scan error : %v", err)
		return
	}

	return
}

func (r repository) MarkAllImagesAsNonPrimary(ctx context.Context, campaignID int) (bool, error) {
	query := `
		UPDATE campaign_images 
		SET is_primary = false 
		WHERE campaign_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, campaignID)
	if err != nil {
		return false, err
	}

	return true, nil
}
