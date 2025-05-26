package exam

import (
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/internal/exam/handler"
	"enuma-elish/internal/exam/repository"
	"enuma-elish/internal/exam/service"
	"enuma-elish/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Exam struct {
	*gin.Engine
	c *config.Config
	i *infra.Infra
	v *validator.Validate
}

func New(c *config.Config, i *infra.Infra, r *gin.Engine, v *validator.Validate) *Exam {
	return &Exam{
		c:      c,
		i:      i,
		Engine: r,
		v:      v,
	}
}

func (e *Exam) Init() {
	r := repository.New(e.i.Postgres, e.i.Redis)
	s := service.New(e.c, r)
	h := handler.New(s, e.v)

	authMiddleware := middleware.Auth(e.c.JWT.Secret)

	v1 := e.Group("/api/v1/exam").Use(authMiddleware)

	v1.POST("", h.CreateExam)
	v1.GET("", h.GetListExams)
	v1.GET("/:exam_id", h.GetDetailExam)
	v1.PUT("/:exam_id", h.UpdateExam)
	v1.DELETE("/:exam_id", h.DeleteExam)

	v1.POST("/assign", h.AssignExamToClass)
	v1.POST("/grade", h.GradeExam)
	v1.GET("/:exam_id/students", h.GetExamStudents)

	studentV1 := e.Group("/api/v1/student/exam")
	studentV1.GET("", h.GetStudentExams)
	studentV1.GET("/:exam_id", h.GetStudentExamDetail)
	studentV1.POST("/submit", h.SubmitExamAnswers)
}
