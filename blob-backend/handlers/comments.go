package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaoleau/blob/models"
	"github.com/joaoleau/blob/usecases"
)

type CommentHandler struct {
	commentUseCase usecases.CommentUseCase
}

func NewCommentHandler(commentUseCase usecases.CommentUseCase) CommentHandler {
	return CommentHandler{
		commentUseCase: commentUseCase,
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	blobID := c.Param("blobId")
	
	blobUUID, err := uuid.Parse(blobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID. Must be in UUID format."})
		return
	}
	
	var comment models.Comment
	comment.BlobID = blobUUID
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	newComment, err := h.commentUseCase.AddComment(c, &comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newComment)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")

	commentUUID, err := uuid.Parse(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID. Must be in UUID format."})
		return
	}

	if err := h.commentUseCase.RemoveComment(c, commentUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}


func (h *CommentHandler) ListCommentsByBlobID(c *gin.Context) {
	blobID := c.Param("blobId")
	
	blobUUID, err := uuid.Parse(blobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blob ID. Must be in UUID format."})
		return
	}

	comments, err := h.commentUseCase.ListCommentsByBlobID(c, blobUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user email not found in context"})
		return
	}

	user, err := h.commentUseCase.BlobUseCase.UserUseCase.GetUserByEmail(c, email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
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
		"content": comments,
	}

	c.JSON(http.StatusOK, response)
}