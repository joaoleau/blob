package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaoleau/blob/models"
	"github.com/joaoleau/blob/usecases"
	"github.com/joaoleau/blob/utils"
)

type BlobHandler struct {
	blobUseCase usecases.BlobUseCase
}

func NewBlobHandler(useCase usecases.BlobUseCase) BlobHandler {
	return BlobHandler{
		blobUseCase: useCase,
	}
}


func (h *BlobHandler) RegisterBlob(ctx *gin.Context) {
	var blob models.Blob
	if err := ctx.ShouldBindJSON(&blob); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data."})
		return
	}

	createdBlob, err := h.blobUseCase.RegisterBlob(ctx, &blob)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blob."})
		return
	}

	ctx.JSON(http.StatusCreated, createdBlob)
}


func (h *BlobHandler) UpdateBlob(ctx *gin.Context) {
	blobID := ctx.Param("blobId")

	blobUUID, err := uuid.Parse(blobID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blob ID. Must be in UUID format."})
		return
	}

	var blob models.Blob
	if err := ctx.ShouldBindJSON(&blob); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data."})
		return
	}
	blob.BlobID = blobUUID

	updatedBlob, err := h.blobUseCase.UpdateBlob(ctx, &blob)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blob."})
		return
	}

	ctx.JSON(http.StatusOK, updatedBlob)
}


func (h *BlobHandler) DeleteBlob(ctx *gin.Context) {
	blobID := ctx.Param("blobId")

	blobUUID, err := uuid.Parse(blobID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blob ID. Must be in UUID format."})
		return
	}

	if err := h.blobUseCase.DeleteBlob(ctx, blobUUID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blob."})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}


func (h *BlobHandler) GetBlobByID(ctx *gin.Context) {
	blobID := ctx.Param("blobId")

	blobUUID, err := uuid.Parse(blobID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blob ID. Must be in UUID format."})
		return
	}

	blob, err := h.blobUseCase.GetBlobByID(ctx, blobUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blob."})
		return
	}

	if blob == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Blob not found."})
		return
	}

	ctx.JSON(http.StatusOK, blob)
}


func (h *BlobHandler) ListBlobs(ctx *gin.Context) {
	pq, err := utils.GetPaginationFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination query."})
		return
	}

	blobList, err := h.blobUseCase.ListBlobs(ctx, pq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blobs."})
		return
	}

	ctx.JSON(http.StatusOK, blobList)
}
