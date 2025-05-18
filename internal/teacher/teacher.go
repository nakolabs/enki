package teacher

import (
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/internal/teacher/handler"
	"enuma-elish/internal/teacher/repository"
	"enuma-elish/internal/teacher/service"
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
	s := service.New(sc.c, r)
	h := handler.New(s, sc.v)

	authMiddleware := middleware.Auth(sc.c.JWT.Secret)

	v1 := sc.Group("/api/v1/teacher").Use(authMiddleware)
	v1.GET("/school/:school_id", h.ListTeachers)

	v1.POST("/invite", h.InviteTeacher)
	v1.POST("/invite/verify", h.VerifyTeacherEmail)
	v1.POST("/invite/update", h.UpdateTeacherAfterInvite)
}
