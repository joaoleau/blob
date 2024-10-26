package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joaoleau/blob/config"
	"github.com/joaoleau/blob/internal/controller/routes"
	"github.com/joho/godotenv"
)

func init() {
	config.LoadConfig()
	config.ConnectToDB()
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	g := gin.Default()
	routes.InitRoutes(&g.RouterGroup)
	
	if err := g.Run(config.GetServerPort()); err != nil {
		log.Fatal(err)
	}
}