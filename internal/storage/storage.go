package storage

import (
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/internal/storage/handler"
	"enuma-elish/internal/storage/repository"
	"enuma-elish/internal/storage/service"
	"enuma-elish/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Storage struct {
	*gin.Engine
	c *config.Config
	i *infra.Infra
	v *validator.Validate
}

func New(c *config.Config, i *infra.Infra, r *gin.Engine, v *validator.Validate) *Storage {
	return &Storage{
		c:      c,
		i:      i,
		Engine: r,
		v:      v,
	}
}

func (s *Storage) Init() {
	r := repository.New(s.i.Postgres)
	svc := service.New(s.i.Cloudinary, r, s.c)
	h := handler.New(svc, s.v)

	storage := s.Group("/api/v1/storage").Use(middleware.Auth(s.c.JWT.Secret))

	// Storage endpoints
	storage.POST("/image", h.StoreImage)
	storage.POST("/video", h.StoreVideo)
	storage.POST("/document", h.StoreDocument)

	// File management endpoints
	storage.DELETE("/file", h.DeleteFile)
	storage.GET("/file/:publicId", h.GetFile)
	storage.GET("/serve/:publicId", h.ServeFile)
	storage.GET("/history", h.GetStorageHistory)
}
