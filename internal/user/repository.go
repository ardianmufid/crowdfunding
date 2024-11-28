package user

import (
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository {
	return repository{
		db: db,
	}
}

func (r repository) Save(user User) (User, error) {

	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

// type repository struct {
// 	db *sqlx.DB
// }

// func NewRepository(db *sqlx.DB) repository {
// 	return repository{
// 		db: db,
// 	}
// }

// func (r repository) Save(ctx context.Context, user UserEntity) error {

// 	query := `
//             INSERT INTO users (
//                 name, email, occupation, password_hash, role, avatar_file_name, created_at, updated_at
//             ) VALUES (
//                 :name, :email, :occupation, :password_hash, :role, :avatar_file_name, :created_at, :updated_at
//             )
//     `

// 	stmt, err := r.db.PrepareNamedContext(ctx, query)
// 	if err != nil {
// 		return err
// 	}

// 	defer stmt.Close()

// 	_, err = stmt.ExecContext(ctx, user)

// 	return err
// }
