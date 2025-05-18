package ppdb

import (
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/internal/ppdb/handler"
	"enuma-elish/internal/ppdb/repository"
	"enuma-elish/internal/ppdb/service"
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

func (sc *PPDB) Init() {
	r := repository.New(sc.i.Postgres)
	s := service.New(r, sc.c)
	h := handler.New(s, sc.v)

	v1 := sc.Group("/api/v1/ppdb")
	v1.POST("/create", h.CreateSchool)
}
