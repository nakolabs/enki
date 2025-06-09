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

func (t *Teacher) Init() {
	r := repository.New(t.i.Postgres, t.i.Redis)
	s := service.New(t.c, r)
	h := handler.New(s, t.v)

	authMiddleware := middleware.Auth(t.c.JWT.Secret)

	v1 := t.Group("/api/v1/teacher").Use(authMiddleware)
	v1.GET("", h.ListTeachers)
	v1.GET("/:teacher_id", h.GetDetailTeacher)
	v1.GET("/statistic", h.GetTeacherStatistics)
	v1.DELETE("/:teacher_id", h.DeleteTeacher)
	v1.PUT("/class", h.UpdateTeacherClass)

	v1.POST("/invite", h.InviteTeacher)
	v1.POST("/invite/verify", h.VerifyTeacherEmail)
	v1.POST("/invite/complete", h.UpdateTeacherAfterInvite)

	v1.GET("/:teacher_id/subjects", h.GetTeacherSubjects)
	v1.GET("/:teacher_id/classes", h.GetTeacherClasses)
}
