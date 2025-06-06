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

func (s *School) Init() {
	r := repository.New(s.i.Postgres)
	svc := service.New(r, s.c)
	h := handler.New(svc, s.v)

	authMiddleware := middleware.Auth(s.c.JWT.Secret)

	v1 := s.Group("/api/v1/school").Use(authMiddleware)
	v1.POST("", h.CreateSchool)
	v1.GET("/:school_id", h.GetDetailSchool)
	v1.GET("", h.ListSchool)
	v1.DELETE("/:school_id", h.DeleteSchool)
	v1.GET("/:school_id/switch", h.SwitchSchool)
	v1.PUT("/:school_id", h.UpdateSchoolProfile)
}
