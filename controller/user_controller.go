package controller

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaoleau/blob/models"
	"github.com/joaoleau/blob/usecases"
	"github.com/joaoleau/blob/utils"
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


func (u *userController) FindUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	
	email = strings.TrimSpace(email)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format."})
		return
	}

	user, err := u.authUserCase.GetByEmail(ctx, email)
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

/// tem q arrumar
func (u *userController) UpdateUser(ctx *gin.Context) {
	var user models.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data. Ensure JSON format is correct."})
		return
	}

	insertedUser, err := u.authUserCase.Update(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to uptade user. Please try again."})
		return
	}

	ctx.JSON(http.StatusCreated, insertedUser)
}


func (u *userController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID. Must be in UUID format."})
		return
	}

	resp := u.authUserCase.Delete(ctx, userUUID)
	if resp != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error."})
		return
	}

	ctx.JSON(http.StatusOK,  gin.H{"message": "Removed user success"})
}


func (u *userController) GetUsers(ctx *gin.Context) {
	pagination, err := utils.GetPaginationFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid pagination parameters."})
		return
	}

	users, err := u.authUserCase.GetUsers(ctx, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users."})
		return
	}

	ctx.JSON(http.StatusOK, users)
}


func (u *userController) FindByName(ctx *gin.Context) {
	name := ctx.Param("nickname")
	name = strings.TrimSpace(name)
	
	pagination, err := utils.GetPaginationFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid pagination parameters."})
		return
	}

	users, err := u.authUserCase.FindByName(ctx, name, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users."})
		return
	}

	ctx.JSON(http.StatusOK, users)
}