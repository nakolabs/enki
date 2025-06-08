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

type Student struct {
	*gin.Engine
	c *config.Config
	i *infra.Infra
	v *validator.Validate
}

func New(c *config.Config, i *infra.Infra, r *gin.Engine, v *validator.Validate) *Student {
	return &Student{
		c:      c,
		i:      i,
		Engine: r,
		v:      v,
	}
}

func (s *Student) Init() {
	r := repository.New(s.i.Postgres, s.i.Redis)
	svc := service.New(r, s.c)
	h := handler.New(svc, s.v)

	authMiddleware := middleware.Auth(s.c.JWT.Secret)

	v1 := s.Group("/api/v1/student").Use(authMiddleware)

	v1.GET("", h.ListStudent)
	v1.GET("/:student_id", h.GetDetailStudent)
	v1.DELETE("/:student_id", h.DeleteStudent)

	v1.POST("/invite", h.InviteStudent)
	v1.POST("/invite/verify", h.VerifyStudentEmail)
	v1.POST("/invite/complete", h.UpdateStudentAfterInvite)
	v1.PUT("/class", h.UpdateStudentClass)
}
