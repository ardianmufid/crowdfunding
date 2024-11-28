package user

import "time"

type User struct {
	Id             int       `db:"id"`
	Name           string    `db:"name"`
	Email          string    `db:"email"`
	PasswordHash   string    `db:"password_hash"`
	Occupation     string    `db:"occupation"`
	AvatarFileName string    `db:"avatar_file_name"`
	Role           string    `db:"role"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
