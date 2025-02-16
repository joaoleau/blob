package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joaoleau/blob/db"
	"github.com/joaoleau/blob/handlers"
	"github.com/joaoleau/blob/middleware"
	"github.com/joaoleau/blob/repository"
	"github.com/joaoleau/blob/usecases"
	"github.com/joho/godotenv"
)

func SetupRouters(server *gin.Engine, dbConnection *sqlx.DB) {

	userRepository := repository.NewUserRepository(dbConnection)
	userUseCase := usecases.NewUserUseCase(userRepository)
	userHandler := handlers.NewUserHandler(userUseCase)
	
	blobRepository := repository.NewBlobRepository(dbConnection)
	blobUseCase := usecases.NewBlobUseCase(blobRepository, userUseCase)
	blobHandler := handlers.NewBlobHandler(blobUseCase)

	likeRepository := repository.NewLikeRepository(dbConnection, &blobRepository)
	likeUseCase := usecases.NewLikeUseCase(likeRepository, &blobUseCase)
	likeHandler := handlers.NewLikeHandler(likeUseCase)

	commentsRepository := repository.NewCommentRepository(dbConnection, &blobRepository)
	commentsUseCase := usecases.NewCommentUseCase(commentsRepository, &blobUseCase)
	commentsHandler := handlers.NewCommentHandler(commentsUseCase)


	protected := server.Group("/api")
	protected.Use(middleware.AuthMiddleware(dbConnection))


	protected.GET("/secure", func(c *gin.Context) {
		email, exists := c.Get("email")
		if !exists {
			c.JSON(500, gin.H{"error": "Email not found in context"})
			return
		}

		c.JSON(200, gin.H{
			"message": "You have access to this route.",
			"email":  email,
		})
	})
 
	protected.POST("/blob", blobHandler.RegisterBlob)
	protected.DELETE("/blob/:blobId", blobHandler.DeleteBlob)
	protected.GET("/blob/:blobId", blobHandler.GetBlobByID)
	protected.GET("/blob", blobHandler.ListBlobs)

	protected.GET("/interest", blobHandler.ListInterests)

	protected.POST("/blob/:blobId/like", likeHandler.AddLike)
	protected.GET("/blob/:blobId/like", likeHandler.ListLike)
	protected.DELETE("/blob/:blobId/like", likeHandler.RemoveLike)

	protected.POST("/blob/:blobId/comment", commentsHandler.CreateComment)
	protected.GET("/blob/:blobId/comment", commentsHandler.ListCommentsByBlobID)
	protected.DELETE("/blob/:blobId/comment/:commentId", commentsHandler.DeleteComment)

	protected.GET("/user", userHandler.GetUserProfile)
	protected.GET("/user/:username", userHandler.GetUserByUsername)
	protected.PUT("/user", userHandler.UpdateUser)
}


func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	server := gin.Default()
	dbConnection, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConnection.Close()
	
	SetupRouters(server, dbConnection)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
