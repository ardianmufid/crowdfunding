package main

import (
	"crowdfunding/config"
	"crowdfunding/database"
	"crowdfunding/internal/campaign"
	"crowdfunding/internal/transaction"
	"crowdfunding/internal/user"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	filename := "cmd/api/config.yaml"
	if err := config.LoadConfig(filename); err != nil {
		panic(err)
	}

	db, err := database.ConnectPostgresSqlx(config.Cfg.DB)
	if err != nil {
		panic(err)
	}

	if db != nil {
		log.Println("db connected")
	}

	router := gin.Default()
	router.Use(cors.Default())

	api := router.Group("api/v1")

	api.Static("/images", "./images")

	user.Init(api, db)
	campaign.Init(api, db)
	transaction.Init(api, db)

	router.Run(config.Cfg.App.Port)
}
