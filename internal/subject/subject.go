package subject

import (
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/internal/subject/handler"
	"enuma-elish/internal/subject/repository"
	"enuma-elish/internal/subject/service"
	"enuma-elish/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Subject struct {
	*gin.Engine
	c *config.Config
	i *infra.Infra
	v *validator.Validate
}

func New(c *config.Config, i *infra.Infra, r *gin.Engine, v *validator.Validate) *Subject {
	return &Subject{
		c:      c,
		i:      i,
		Engine: r,
		v:      v,
	}
}

func (s *Subject) Init() {
	r := repository.New(s.i.Postgres)
	svc := service.New(r, s.c)
	h := handler.New(svc, s.v)

	authMiddleware := middleware.Auth(s.c.JWT.Secret)

	v1 := s.Group("/api/v1/subject").Use(authMiddleware)
	v1.POST("", h.CreateSubject)
	v1.GET("", h.ListSubject)
	v1.GET("/:subject_id", h.GetDetailSubject)
	v1.PUT("/:subject_id", h.UpdateSubject)
	v1.DELETE("/:subject_id", h.DeleteSubject)
	v1.POST("/assign-teachers", h.AssignTeachersToSubject)
	v1.GET("/:subject_id/teachers", h.GetTeachersBySubject)
	v1.PUT("/class", h.UpdateSubjectClass)
}
