package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/joaoleau/blob/models"
	"github.com/joaoleau/blob/usecases"
	"net/http"
	"github.com/google/uuid"
)

type userController struct {
	authUserCase usecases.AuthUseCase
}

func NewUserontroller(usecases usecases.AuthUseCase) userController {
	return userController{
		authUserCase: usecases,
	}
}

func (u *userController) FindUserByID(ctx *gin.Context) {
	userId := ctx.Param("userId")
	
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID. Must be in UUID format."})
		return
	}

	user, err := u.authUserCase.GetByID(ctx, userUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user."})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
		return
	}

	ctx.JSON(http.StatusOK, user)
}


func (u *userController) CreateUser(ctx *gin.Context) {
	var user models.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data. Ensure JSON format is correct."})
		return
	}

	insertedUser, err := u.authUserCase.Register(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user. Please try again."})
		return
	}

	ctx.JSON(http.StatusCreated, insertedUser)
}


func (u *userController) UpdateUser(ctx *gin.Context) {
}

func (u *userController) DeleteUser(ctx *gin.Context) {
}

// func (p *userController) Getusers(ctx *gin.Context) {

// 	users, err := p.authUserCase.Getusers()
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, users)
// }