package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joaoleau/blob/db"
	"github.com/joaoleau/blob/repository"
	"github.com/joaoleau/blob/controller"
	"github.com/joaoleau/blob/usecases"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	server := gin.Default()
	dbConnection, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConnection.Close()

	//Arrumar um jeito de centralizar esses handlers
	userRepository := repository.NewAuthRepository(dbConnection)
	userUseCase := usecases.NewAuthUseCase(userRepository)
	userController := controller.NewUserontroller(userUseCase)

	
	server.GET("/getUserById/:userId", userController.FindUserByID)
	server.GET("/findUserByEmail/:email", userController.FindUserByEmail)
	server.GET("/findByName/:nickname", userController.FindByName)
	server.GET("/getUsers", userController.GetUsers)
	server.POST("/createUser", userController.CreateUser)
	server.PATCH("/updateUser/:userId", userController.UpdateUser)
	server.DELETE("/deleteUser/:userId", userController.DeleteUser)


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
