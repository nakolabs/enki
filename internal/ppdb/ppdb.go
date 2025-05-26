package ppdb

import (
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/internal/ppdb/handler"
	"enuma-elish/internal/ppdb/repository"
	"enuma-elish/internal/ppdb/service"
	"enuma-elish/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PPDB struct {
	*gin.Engine
	c *config.Config
	i *infra.Infra
	v *validator.Validate
}

func New(c *config.Config, i *infra.Infra, r *gin.Engine, v *validator.Validate) *PPDB {
	return &PPDB{
		c:      c,
		i:      i,
		Engine: r,
		v:      v,
	}
}

func (p *PPDB) Init() {
	r := repository.New(p.i.Postgres)
	svc := service.New(r, p.c)
	h := handler.New(svc, p.v)

	authMiddleware := middleware.Auth(p.c.JWT.Secret)

	v1 := p.Group("/api/v1/ppdb").Use(authMiddleware)

	// PPDB management (for school admin)
	v1.POST("", h.CreatePPDB)
	v1.GET("", h.GetListPPDB)
	v1.GET("/:ppdb_id", h.GetPPDBByID)
	v1.PUT("/:ppdb_id", h.UpdatePPDB)
	v1.DELETE("/:ppdb_id", h.DeletePPDB)

	// PPDB registration and management - Now requires authentication
	v1.POST("/register", h.RegisterPPDB)
	v1.GET("/registrants", h.GetPPDBRegistrants)
	v1.POST("/select", h.SelectPPDBStudents)
}
