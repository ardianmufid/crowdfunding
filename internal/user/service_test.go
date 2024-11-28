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
		Email:      "antp@gmail.com",
		Password:   "mysecretpassword",
		Occupation: "qa",
	}
	user, err := svc.RegisterUser(req)
	require.Nil(t, err)
	require.NotNil(t, user)
	log.Println(user)
}
