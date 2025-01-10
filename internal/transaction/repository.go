package transaction

import (
	"context"
	"crowdfunding/internal/user"
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
			&transaction.User_ID,
			&transaction.Campaign_ID,
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
