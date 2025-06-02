package handler

import (
	"enuma-elish/internal/storage/service"
	"enuma-elish/internal/storage/service/data/request"
	"enuma-elish/internal/storage/service/data/response"
	commonError "enuma-elish/pkg/error"
	commonHttp "enuma-elish/pkg/http"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service   service.Service
	validator *validator.Validate
}

func New(s service.Service, validator *validator.Validate) *Handler {
	return &Handler{
		service:   s,
		validator: validator,
	}
}

func (h *Handler) StoreImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	defer file.Close()

	req := request.StoreImageRequest{
		File:   file,
		Header: header,
		Folder: c.PostForm("folder"),
	}

	result, err := h.service.StoreImage(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("image stored successfully").
		SetData(result)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) StoreVideo(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	defer file.Close()

	req := request.StoreVideoRequest{
		File:   file,
		Header: header,
		Folder: c.PostForm("folder"),
	}

	result, err := h.service.StoreVideo(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("video stored successfully").
		SetData(result)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) StoreDocument(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	defer file.Close()

	req := request.StoreDocumentRequest{
		File:   file,
		Header: header,
		Folder: c.PostForm("folder"),
	}

	result, err := h.service.StoreDocument(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("document stored successfully").
		SetData(result)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteFile(c *gin.Context) {
	var req request.DeleteFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	result, err := h.service.DeleteFile(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("file deleted successfully").
		SetData(result)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetFile(c *gin.Context) {
	publicID := c.Param("publicId")
	if publicID == "" {
		c.Error(commonError.New("public id cannot be empty", 422)).SetType(gin.ErrorTypeBind)
		return
	}

	result, err := h.service.GetFile(c.Request.Context(), publicID)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("file retrieved successfully").
		SetData(result)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ServeFile(c *gin.Context) {
	publicID := c.Param("publicId")
	if publicID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "public_id is required"})
		return
	}

	result, err := h.service.GetFile(c.Request.Context(), publicID)
	if err != nil {
		if err == commonError.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		if err == commonError.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file"})
		return
	}

	// Set appropriate headers
	c.Header("Content-Type", result.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", result.Filename))
	c.Header("Cache-Control", "private, max-age=3600") // Cache for 1 hour

	// Serve the raw file content
	c.Data(http.StatusOK, result.ContentType, result.Content)
}

func (h *Handler) GetStorageHistory(c *gin.Context) {
	httpQuery := request.GetStorageHistoryQuery{}
	httpQuery.Query = commonHttp.DefaultQuery()
	if err := c.ShouldBindQuery(&httpQuery); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	fileType := c.Query("type")
	var result *response.StorageHistoryResponse
	var meta *commonHttp.Meta
	var err error

	if fileType != "" {
		result, meta, err = h.service.GetStorageHistoryByType(c.Request.Context(), fileType, httpQuery)
	} else {
		result, meta, err = h.service.GetStorageHistory(c.Request.Context(), httpQuery)
	}

	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("storage history retrieved successfully").
		SetData(result).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}
