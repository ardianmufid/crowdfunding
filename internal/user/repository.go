package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Save(ctx context.Context, model User) (user User, err error)
	FindByEmail(ctx context.Context, email string) (user User, err error)
	FindByID(ctx context.Context, Id int) (user User, err error)
	Update(ctx context.Context, model User) (user User, err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository {
	return repository{
		db: db,
	}
}

func (r repository) Save(ctx context.Context, model User) (user User, err error) {

	// err := r.db.Create(&user).Error
	// if err != nil {
	// 	return user, err
	// }

	// return user, nil

	query := `
			INSERT INTO users (
					name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at
			) VALUES (
			 		$1, $2, $3, $4, $5, $6, $7, $8
			) RETURNING
			 		id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at
	`
	err = r.db.QueryRowxContext(
		ctx,
		query,
		model.Name,
		model.Occupation,
		model.Email,
		model.PasswordHash,
		model.AvatarFileName,
		model.Role,
		model.CreatedAt,
		model.UpdatedAt).StructScan(&user)
	if err != nil {
		return
	}

	return

}

func (r repository) FindByEmail(ctx context.Context, email string) (user User, err error) {

	// var user User
	// err := r.db.Where("email = ?", email).Find(&user).Error
	// if err != nil {
	// 	return user, err
	// }

	// return user, err

	query := `
			SELECT 
				id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at
			FROM users
			WHERE email=$1
	`

	err = r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("error not found")
			return
		}
		return
	}
	return
}

func (r repository) FindByID(ctx context.Context, Id int) (user User, err error) {

	query := `
		SELECT 
			id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at
		FROM users
		WHERE id=$1
	`

	err = r.db.GetContext(ctx, &user, query, Id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("error not found")
			return
		}
	}

	return

	// var user User
	// err := r.db.Where("id = ?", Id).Find(&user).Error
	// if err != nil {
	// 	return user, err
	// }
	// return user, nil

}

func (r repository) Update(ctx context.Context, model User) (user User, err error) {

	query := `
	UPDATE users
		SET 
			name = $1,
			email = $2,
			password_hash = $3,
			occupation = $4,
			avatar_file_name = $5,
			role = $6,
			updated_at = $7
	WHERE id = $8
	RETURNING id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at
`

	err = r.db.QueryRowxContext(
		ctx,
		query,
		model.Name,
		model.Email,
		model.PasswordHash,
		model.Occupation,
		model.AvatarFileName,
		model.Role,
		model.UpdatedAt,
		model.Id,
	).StructScan(&user)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New("error not found")
			return
		}
		return
	}

	return
	// err := r.db.Save(&user).Error
	// if err != nil {
	// 	return user, err
	// }
	// return user, nil
}
