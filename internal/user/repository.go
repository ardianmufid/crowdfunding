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

func (r repository) FindByEmail(email string) (User, error) {

	var user User
	err := r.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, err
}

func (r repository) FindByID(Id int) (User, error) {

	var user User
	err := r.db.Where("id = ?", Id).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r repository) Update(user User) (User, error) {

	err := r.db.Save(&user).Error
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
