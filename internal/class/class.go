package class

import (
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/internal/class/handler"
	"enuma-elish/internal/class/repository"
	"enuma-elish/internal/class/service"
	"enuma-elish/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Class struct {
	*gin.Engine
	c *config.Config
	i *infra.Infra
	v *validator.Validate
}

func New(c *config.Config, i *infra.Infra, r *gin.Engine, v *validator.Validate) *Class {
	return &Class{
		c:      c,
		i:      i,
		Engine: r,
		v:      v,
	}
}

func (cl *Class) Init() {
	r := repository.New(cl.i.Postgres)
	s := service.New(r, cl.c)
	h := handler.New(s, cl.v)

	authMiddleware := middleware.Auth(cl.c.JWT.Secret)

	v1 := cl.Group("/api/v1/class").Use(authMiddleware)
	v1.POST("", h.CreateClass)
	v1.GET("", h.ListClass)
	v1.GET("/:class_id", h.GetDetailClass)
	v1.PUT("/:class_id", h.UpdateClass)
	v1.DELETE("/:class_id", h.DeleteClass)
	v1.POST("/add-students", h.AddStudentsToClass)
	v1.POST("/assign-teachers", h.AddTeachersToClass)
	v1.POST("/add-subjects", h.AddSubjectsToClass)
	v1.GET("/:class_id/students", h.GetStudentsByClass)
	v1.GET("/:class_id/teachers", h.GetTeachersByClass)
	v1.GET("/:class_id/subjects", h.GetSubjectsByClass)
	v1.DELETE("/teacher", h.RemoveTeachersFromClass)
	v1.DELETE("/student", h.RemoveStudentsFromClass)
	v1.DELETE("/subject", h.RemoveSubjectsFromClass)
}
