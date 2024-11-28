package user

import (
	"crowdfunding/config"
	"crowdfunding/database"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

var svc service

func init() {
	filename := "../../cmd/api/config.yaml"
	err := config.LoadConfig(filename)
	if err != nil {
		panic(err)
	}

	db, err := database.ConnectPostgres(config.Cfg.DB)
	if err != nil {
		panic(err)
	}

	repo := NewRepository(db)
	svc = NewService(repo)
}

func TestRegister_Success(t *testing.T) {
	req := RegisterUserRequest{
		Name: "anto",
		// Email:      fmt.Sprintf("%v@google.id", uuid.NewString()),
		Email:      "anto@gmail.com",
		Password:   "mysecretpassword",
		Occupation: "qa",
	}
	user, err := svc.RegisterUser(req)
	require.Nil(t, err)
	require.NotNil(t, user)
	log.Println(user)
}

func TestLogin_Success(t *testing.T) {
	req := LoginUserRequest{
		Email:    "ardian@gmail.com",
		Password: "mysecretpassword",
	}

	user, err := svc.LoginUser(req)
	require.Nil(t, err)
	require.NotNil(t, user)
	log.Println(user)
}

func TestIsEmailAvailable_Success(t *testing.T) {
	req := CheckEmailRequest{
		Email: "ardian@gmail.com",
	}
	isEmailAvailable, err := svc.IsEmailAvailable(req)
	require.Nil(t, err)
	require.Equal(t, false, isEmailAvailable)
}
