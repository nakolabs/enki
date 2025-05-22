package school

import (
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/internal/school/handler"
	"enuma-elish/internal/school/repository"
	"enuma-elish/internal/school/service"
	"enuma-elish/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type School struct {
	*gin.Engine
	c *config.Config
	i *infra.Infra
	v *validator.Validate
}

func New(c *config.Config, i *infra.Infra, r *gin.Engine, v *validator.Validate) *School {
	return &School{
		c:      c,
		i:      i,
		Engine: r,
		v:      v,
	}
}

func (sc *School) Init() {
	r := repository.New(sc.i.Postgres)
	s := service.New(r, sc.c)
	h := handler.New(s, sc.v)

	authMiddleware := middleware.Auth(sc.c.JWT.Secret)

	v1 := sc.Group("/api/v1/school").Use(authMiddleware)
	v1.POST("", h.CreateSchool)
	v1.GET("/:school_id", h.GetDetailSchool)
	v1.GET("", h.ListSchool)
	v1.GET("/switch/:school_id", h.SwitchSchool)
}
