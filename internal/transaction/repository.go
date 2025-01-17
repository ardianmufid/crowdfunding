package transaction

import (
	"context"
	"crowdfunding/internal/campaign"
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

func (r repository) GetCampaignID(ctx context.Context, CampaignID int) (transactions []Transaction, err error) {

	query := `
        SELECT 
            t.id, t.user_id, t.campaign_id, t.amount, t.status, t.code, t.created_at, t.updated_at, 
            u.id, u.name, u.occupation, u.email, u.password_hash, u.avatar_file_name, u.role, u.created_at, u.updated_at
        FROM 
            transactions t
        LEFT JOIN 
            users u
        ON t.user_id = u.id
        WHERE t.campaign_id = $1
        ORDER BY t.id DESC
    `

	rows, err := r.db.QueryxContext(ctx, query, CampaignID)
	if err != nil {
		log.Printf("Scan error : %v", err)

		return []Transaction{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var transaction Transaction
		var user user.User

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.CampaignID,
			&transaction.Amount,
			&transaction.Status,
			&transaction.Code,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&user.Id,
			&user.Name,
			&user.Occupation,
			&user.Email,
			&user.PasswordHash,
			&user.AvatarFileName,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Printf("error : %v", err)

			return []Transaction{}, err
		}

		transaction.User = user

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// func (r repository) GetByUserID(ctx context.Context, userID int) (transactions []Transaction, err error) {

// 	query := `
// 		SELECT
// 			t.id, t.user_id, t.campaign_id, t.amount, t.status, t.code, t.created_at, t.updated_at,
// 			c.id, c.user_id, c.name, c.short_description, c.description,
//             c.goal_amount, c.current_amount, c.perks, c.becker_count,
//             c.slug, c.created_at, c.updated_at,
//             ci.id, ci.campaign_id, ci.file_name, ci.is_primary, ci.created_at, ci.updated_at
// 		FROM
// 			transactions t
// 		LEFT JOIN
// 			campaigns c
// 		ON t.campaign_id = c.id
// 		LEFT JOIN
// 			campaign_images ci
// 		ON c.id = ci.campaign_id AND is_primary = true
// 		WHERE
// 			t.user_id = $1
// 		ORDER BY
// 			t.id DESC
// 	`

// 	rows, err := r.db.QueryxContext(ctx, query, userID)
// 	if err != nil {
// 		log.Printf("Scan error : %v", err)

// 		return []Transaction{}, err
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		var transaction Transaction
// 		var kampanye campaign.Campaign
// 		var campaignImages campaign.CampaignImage

// 		err := rows.Scan(
// 			&transaction.ID,
// 			&transaction.User_ID,
// 			&transaction.Campaign_ID,
// 			&transaction.Amount,
// 			&transaction.Status,
// 			&transaction.Code,
// 			&transaction.CreatedAt,
// 			&transaction.UpdatedAt,
// 			&kampanye.ID,
// 			&kampanye.UserID,
// 			&kampanye.Name,
// 			&kampanye.ShortDescription,
// 			&kampanye.Description,
// 			&kampanye.GoalAmount,
// 			&kampanye.CurrentAmount,
// 			&kampanye.Perks,
// 			&kampanye.BeckerCount,
// 			&kampanye.Slug,
// 			&kampanye.CreatedAt,
// 			&kampanye.UpdatedAt,
// 			&campaignImages.ID,
// 			&campaignImages.CampaignId,
// 			&campaignImages.FileName,
// 			&campaignImages.IsPrimary,
// 			&campaignImages.CreatedAt,
// 			&campaignImages.UpdatedAt,
// 		)
// 		if err != nil {
// 			log.Printf("error : %v", err)

// 			return []Transaction{}, err
// 		}

// 		transaction.Campaign = kampanye

// 		if kampanye.CampaignImages == nil {
// 			kampanye.CampaignImages = &[]campaign.CampaignImage{} // Inisialisasi slice jika nil
// 		}

// 		// Set Campaign Image only if not NULL
// 		if campaignImages.ID != nil {
// 			*kampanye.CampaignImages = append(*kampanye.CampaignImages, campaignImages)

// 		}

// 		transactions = append(transactions, transaction)
// 	}

// 	return
// }

func (r repository) GetByUserID(ctx context.Context, userID int) ([]Transaction, error) {
	var transactions []Transaction

	// Query untuk mengambil data Transactions berdasarkan user_id
	transactionsQuery := `
		SELECT 
			id, user_id, campaign_id, amount, status, code, created_at, updated_at
		FROM 
			transactions
		WHERE 
			user_id = $1
		ORDER BY 
			id DESC
	`
	err := r.db.SelectContext(ctx, &transactions, transactionsQuery, userID)
	if err != nil {
		return nil, err
	}

	// Iterasi setiap transaksi untuk melengkapi data campaign dan campaign images
	for i, transaction := range transactions {
		// Query untuk mengambil data Campaign terkait
		campaignQuery := `
			SELECT
				id, user_id, name, short_description, description,
				goal_amount, current_amount, perks, becker_count,
				slug, created_at, updated_at
			FROM campaigns
			WHERE id = $1
		`
		var kampanye campaign.Campaign
		err = r.db.GetContext(ctx, &kampanye, campaignQuery, transaction.CampaignID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("campaign with ID %d not found", transaction.CampaignID)
			}
			return nil, err
		}

		// Query untuk mengambil CampaignImages yang terkait
		imagesQuery := `
			SELECT
				id, campaign_id, file_name, is_primary, created_at, updated_at
			FROM campaign_images
			WHERE campaign_id = $1
		`
		var campaignImages []campaign.CampaignImage
		err = r.db.SelectContext(ctx, &campaignImages, imagesQuery, kampanye.ID)
		if err != nil {
			return nil, err
		}
		kampanye.CampaignImages = &campaignImages

		// Assign campaign ke transaction
		transactions[i].Campaign = kampanye
	}

	return transactions, nil
}

func (r repository) Save(ctx context.Context, model Transaction) (transaction Transaction, err error) {

	log.Printf("Query parameters: user_id=%d, campaign_id=%d, amount=%d, status=%s, code=%v, created_at=%v, updated_at=%v",
		model.UserID,
		model.CampaignID,
		model.Amount,
		model.Status,
		model.Code,
		model.CreatedAt,
		model.UpdatedAt,
	)

	query := `
		INSERT INTO transactions (
			user_id, campaign_id, amount, status, code, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, user_id, campaign_id, amount, status, code, created_at, updated_at
	`

	err = r.db.QueryRowxContext(
		ctx,
		query,
		model.UserID,
		model.CampaignID,
		model.Amount,
		model.Status,
		model.Code,
		model.CreatedAt,
		model.UpdatedAt,
	).StructScan(&transaction)

	if err != nil {
		return
	}
	return
}

func (r repository) Update(ctx context.Context, model Transaction) (transaction Transaction, err error) {

	query := `
		UPDATE transactions
		SET
			user_id = $1,
			campaign_id = $2,
			amount = $3,
			status = $4,
			code = $5,
			payment_url = $6,
			updated_at = $7
		WHERE id = $8
		RETURNING
			id, user_id, campaign_id, amount, status, code, payment_url, created_at, updated_at
	`

	err = r.db.QueryRowxContext(
		ctx,
		query,
		model.UserID,
		model.CampaignID,
		model.Amount,
		model.Status,
		model.Code,
		model.PaymentURL,
		model.UpdatedAt,
		model.ID,
	).StructScan(&transaction)
	if err != nil {
		log.Printf("Scan error : %v", err)
		if err == sql.ErrNoRows {
			return
		}
		return
	}

	return
}

func (r repository) GetByID(ctx context.Context, ID int) (Transaction, error) {
	var transaction Transaction

	// Query untuk mengambil data Transactions berdasarkan user_id
	transactionsQuery := `
		SELECT 
			id, user_id, campaign_id, amount, status, code, created_at, updated_at
		FROM 
			transactions
		WHERE 
			id = $1
		ORDER BY 
			id DESC
	`
	err := r.db.SelectContext(ctx, &transaction, transactionsQuery, ID)
	if err != nil {
		return Transaction{}, err
	}

	// Iterasi setiap transaksi untuk melengkapi data campaign dan campaign images
	// for i, transaction := range transactions {
	// Query untuk mengambil data Campaign terkait
	campaignQuery := `
			SELECT
				id, user_id, name, short_description, description,
				goal_amount, current_amount, perks, becker_count,
				slug, created_at, updated_at
			FROM campaigns
			WHERE id = $1
		`
	var kampanye campaign.Campaign
	err = r.db.GetContext(ctx, &kampanye, campaignQuery, transaction.CampaignID)
	if err != nil {
		if err == sql.ErrNoRows {
			return Transaction{}, fmt.Errorf("campaign with ID %d not found", transaction.CampaignID)
		}
		return Transaction{}, err
	}

	// Query untuk mengambil CampaignImages yang terkait
	imagesQuery := `
			SELECT
				id, campaign_id, file_name, is_primary, created_at, updated_at
			FROM campaign_images
			WHERE campaign_id = $1
		`
	var campaignImages []campaign.CampaignImage
	err = r.db.SelectContext(ctx, &campaignImages, imagesQuery, kampanye.ID)
	if err != nil {
		return Transaction{}, err
	}
	kampanye.CampaignImages = &campaignImages

	// Assign campaign ke transaction
	transaction.Campaign = kampanye
	// }

	return transaction, nil
}
