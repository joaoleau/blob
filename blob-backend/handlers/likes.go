package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaoleau/blob/usecases"
)

type LikeHandler struct {
	likeUseCase usecases.LikeUseCase
}

func NewLikeHandler(likeUseCase usecases.LikeUseCase) LikeHandler {
	return LikeHandler{
		likeUseCase: likeUseCase,
	}
}

func (h *LikeHandler) AddLike(ctx *gin.Context) {
	blobID := ctx.Param("blobId")

	blobUUID, err := uuid.Parse(blobID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blob ID. Must be in UUID format."})
		return
	}

	newLike, err := h.likeUseCase.AddLike(ctx, blobUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, newLike)
}

func (h *LikeHandler) RemoveLike(ctx *gin.Context) {
	blobID := ctx.Param("blobId")

	blobUUID, err := uuid.Parse(blobID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blob ID. Must be in UUID format."})
		return
	}

	if err := h.likeUseCase.RemoveLike(ctx, blobUUID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove like."})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h *LikeHandler) ListLike(ctx *gin.Context) {
	blobID := ctx.Param("blobId")

	blobUUID, err := uuid.Parse(blobID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blob ID. Must be in UUID format."})
		return
	}

	likes, err := h.likeUseCase.ListLikesByBlobID(ctx, blobUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	email, exists := ctx.Get("email")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email não encontrado no contexto"})
		return
	}

	user, err := h.likeUseCase.BlobUseCase.UserUseCase.GetUserByEmail(ctx, email.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recuperar informações do usuário"})
		return
	}

	response := gin.H{
		"user_logon": map[string]interface{}{
			"username":    user.Username,
			"avatar_icon": user.AvatarIcon,
			"avatar_color": user.AvatarColor,
			"id": user.ID,
			"email": user.Email,
		},
		"content": likes,
	}

	ctx.JSON(http.StatusOK, response)
}