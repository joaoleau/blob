package handlers

import (
	"net/http"
	"github.com/joaoleau/blob/models"
	"github.com/gin-gonic/gin"
	"github.com/joaoleau/blob/usecases"
)

type UserHandler struct {
	userUseCase *usecases.UserUseCase
}

func NewUserHandler(useCase *usecases.UserUseCase) UserHandler {
	return UserHandler{
		userUseCase: useCase,
	}
}

func (h *UserHandler) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	userWithBlobs, err := h.userUseCase.GetUserByUsername(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user."})
		return
	}

	if userWithBlobs == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
		return
	}

	ctx.JSON(http.StatusOK, userWithBlobs)
}

func (h *UserHandler) GetUserProfile(ctx *gin.Context) {
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email não encontrado no contexto"})
		return
	}

	user, err := h.userUseCase.GetUserByEmail(ctx, email.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recuperar informações do usuário"})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email não encontrado no contexto"})
		return
	}
	
	var userData models.User
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.userUseCase.UpdateUser(ctx, email.(string), userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}