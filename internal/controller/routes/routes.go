package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joaoleau/blob/internal/controller/routes"
)

func InitRoutes(r *gin.RouterGroup) {

	r.GET("/getUserById/:userId", routes.FindUserByID)
	r.POST("/createUser", routes.CreateUser)
	r.PUT("/updateUser/:userId", routes.UpdateUser)
	r.DELETE("/deleteUser/:userId", routes.DeleteUser)
}