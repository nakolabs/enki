package student

import (
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/internal/student/handler"
	"enuma-elish/internal/student/repository"
	"enuma-elish/internal/student/service"
	"enuma-elish/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Teacher struct {
	*gin.Engine
	c *config.Config
	i *infra.Infra
	v *validator.Validate
}

func New(c *config.Config, i *infra.Infra, r *gin.Engine, v *validator.Validate) *Teacher {
	return &Teacher{
		c:      c,
		i:      i,
		Engine: r,
		v:      v,
	}
}

func (sc *Teacher) Init() {
	r := repository.New(sc.i.Postgres, sc.i.Redis)
	s := service.New(r, sc.c)
	h := handler.New(s, sc.v)

	authMiddleware := middleware.Auth(sc.c.JWT.Secret)

	v1 := sc.Group("/api/v1/student").Use(authMiddleware)

	v1.POST("/invite", h.InviteStudent)
	v1.POST("/invite/verify", h.VerifyStudentEmail)
	v1.POST("/invite/update", h.UpdateStudentAfterInvite)
}
