package user

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	ROLE_ADMIN Role = "admin"
	ROLE_USER  Role = "user"
)

type User struct {
	Id             int       `db:"id"`
	Name           string    `db:"name"`
	Email          string    `db:"email"`
	PasswordHash   string    `db:"password_hash"`
	Occupation     string    `db:"occupation"`
	AvatarFileName string    `db:"avatar_file_name"`
	Role           Role      `db:"role"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func NewFromRegisterRequest(req RegisterUserRequest) User {
	return User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
		Occupation:   req.Occupation,
		Role:         ROLE_USER,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func NewFromLoginRequest(req LoginUserRequest) User {
	return User{
		Email:        req.Email,
		PasswordHash: req.Password,
	}
}

func (u *User) EncryptPassword(salt int) (err error) {

	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), salt)
	if err != nil {
		return
	}

	u.PasswordHash = string(encryptedPass)
	return nil
}

func (u User) VerifyPasswordFromPlain(encrypted string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(u.PasswordHash))
}
func (u User) VerifyPasswordFromEncrypted(plain string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plain))
}
func (u User) VerifyPassword(plainPassword string) error {
	if u.PasswordHash == "" {
		return errors.New("no password hash stored")
	}
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword))
}
