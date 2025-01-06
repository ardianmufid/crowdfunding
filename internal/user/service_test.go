package user

import (
	"context"
	"crowdfunding/config"
	"crowdfunding/database"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var svc service

func init() {
	filename := "../../cmd/api/config.yaml"
	err := config.LoadConfig(filename)
	if err != nil {
		panic(err)
	}

	db, err := database.ConnectPostgresSqlx(config.Cfg.DB)
	if err != nil {
		panic(err)
	}

	repo := NewRepository(db)
	svc = NewService(repo)
}

func TestRegister_Success(t *testing.T) {
	// Arrange
	req := RegisterUserRequest{
		Name:       "novsky",
		Email:      "novsky_" + uuid.NewString() + "@google.id", // Menggunakan email unik untuk mencegah duplikasi
		Password:   "mysecretpassword",
		Occupation: "developer",
	}

	// Act
	user, err := svc.RegisterUser(context.Background(), req)

	// Assert
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, req.Name, user.Name)
	require.Equal(t, req.Email, user.Email)
	require.NotEmpty(t, user.PasswordHash) // Pastikan password sudah terenkripsi
	log.Printf("User registered successfully: %+v\n", user)
}

func TestRegister_EmailDuplicate(t *testing.T) {
	// Arrange
	email := "duplicate@gmail.com"
	req1 := RegisterUserRequest{
		Name:       "User1",
		Email:      email,
		Password:   "password1",
		Occupation: "developer",
	}
	req2 := RegisterUserRequest{
		Name:       "User2",
		Email:      email, // Email yang sama dengan req1
		Password:   "password2",
		Occupation: "tester",
	}

	// Insert user pertama
	user1, err := svc.RegisterUser(context.Background(), req1)
	require.Nil(t, err)
	require.NotNil(t, user1)

	// Act
	user2, err := svc.RegisterUser(context.Background(), req2)

	// Assert
	require.NotNil(t, err)                                  // Harus ada error karena email duplikat
	require.Contains(t, err.Error(), "duplicate key value") // Pastikan error terkait email
	require.Empty(t, user2)
	log.Printf("Expected error due to duplicate email: %v\n", err)
}

func TestLoginUser(t *testing.T) {
	// Arrange: Setup initial user for testing
	email := "test_login@gmail.com"
	password := "securepassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	initialUser := User{
		Name:         "Test User",
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	_, err := svc.repo.Save(context.Background(), initialUser)
	require.Nil(t, err)

	t.Run("LoginUser_Success", func(t *testing.T) {
		// Arrange
		req := LoginUserRequest{
			Email:    email,
			Password: password,
		}

		// Act
		user, err := svc.LoginUser(context.Background(), req)

		// Assert
		require.Nil(t, err)
		require.NotNil(t, user)
		require.Equal(t, email, user.Email)
		log.Printf("Login successful: %+v\n", user)
	})

	t.Run("LoginUser_Error_InvalidEmail", func(t *testing.T) {
		// Arrange
		req := LoginUserRequest{
			Email:    "nonexistent@gmail.com",
			Password: password,
		}

		// Act
		user, err := svc.LoginUser(context.Background(), req)

		// Assert
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "error not found") // Error PostgreSQL untuk data tidak ditemukan
		require.Empty(t, user)
		log.Printf("Error as expected for invalid email: %v\n", err)
	})

	t.Run("LoginUser_Error_InvalidPassword", func(t *testing.T) {
		// Arrange
		req := LoginUserRequest{
			Email:    email,
			Password: "wrongpassword",
		}

		// Act
		user, err := svc.LoginUser(context.Background(), req)

		// Assert
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "password not match")
		require.Empty(t, user)
		log.Printf("Error as expected for invalid password: %v\n", err)
	})
}

// func TestRegister_Success(t *testing.T) {
// 	req := RegisterUserRequest{
// 		Name: "novsky",
// 		// Email:      fmt.Sprintf("%v@google.id", uuid.NewString()),
// 		Email:      "novsky@gmail.com",
// 		Password:   "mysecretpassword",
// 		Occupation: "dev",
// 	}
// 	user, err := svc.RegisterUser(context.Background(), req)
// 	require.Nil(t, err)
// 	require.NotNil(t, user)
// 	log.Println(user)
// }

// func TestLogin_Success(t *testing.T) {
// 	req := LoginUserRequest{
// 		Email:    "isnan@gmail.com",
// 		Password: "mysecretpassword",
// 	}

// 	user, err := svc.LoginUser(context.Background(), req)
// 	require.Nil(t, err)
// 	require.NotNil(t, user)
// 	log.Println(user)
// }
// func TestLogin_Error(t *testing.T) {
// 	req := LoginUserRequest{
// 		Email:    "isnan@gmail.com",
// 		Password: "",
// 	}

// 	user, err := svc.LoginUser(context.Background(), req)
// 	require.NotNil(t, err)
// 	require.Nil(t, user)
// 	// log.Println(user)
// }

// func TestIsEmailAvailable_Success(t *testing.T) {
// 	req := CheckEmailRequest{
// 		Email: "ardian@gmail.com",
// 	}
// 	isEmailAvailable, err := svc.IsEmailAvailable(req)
// 	require.Nil(t, err)
// 	require.Equal(t, false, isEmailAvailable)
// }
