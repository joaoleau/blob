package main

import (
	"context"
	"log"
	"os"

	oidc "github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/joaoleau/blob/db"
	"github.com/joaoleau/blob/handlers"
	"github.com/joaoleau/blob/middlewares"
	"github.com/joaoleau/blob/repository"
	"github.com/joaoleau/blob/usecases"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
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

	ctx := context.Background()
	keycloakAuth := os.Getenv("KEYCLOAK_AUTH")
	provider, err := oidc.NewProvider(ctx, keycloakAuth)
	if err != nil {
		log.Fatal(err)
	}
	
	authConfig := oauth2.Config{
		ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  os.Getenv("REDIRECT_AUTH"),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "roles", "email"},
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: "account"})

	state := os.Getenv("STATE_AUTH")
	
	// Da um jeito nisso /////////////////////
	blobRepository := repository.NewBlobRepository(dbConnection)
	blobUseCase := usecases.NewBlobUseCase(blobRepository)
	blobHandler := handlers.NewBlobHandler(blobUseCase)

	server.POST("/blobs", blobHandler.RegisterBlob)
	server.PUT("/blobs/:blobId", blobHandler.UpdateBlob)
	server.DELETE("/blobs/:blobId", blobHandler.DeleteBlob)
	server.GET("/blobs/:blobId", blobHandler.GetBlobByID)
	server.GET("/blobs", blobHandler.ListBlobs)

	server.GET("/", handlers.InitAuth(&authConfig, state))
	server.GET("/callback", handlers.CallbackAuth(&authConfig, ctx, state))

	server.GET("/protected", middlewares.AuthMiddleware(verifier, ctx), func(c *gin.Context) {
		claims := c.MustGet("claims").(map[string]interface{})
		c.JSON(200, gin.H{"message": "Acesso autorizado!", "claims": claims})
	})
	/////////////////////////////////////////////
	
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
