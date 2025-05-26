package question

import (
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/internal/question/handler"
	"enuma-elish/internal/question/repository"
	"enuma-elish/internal/question/service"
	"enuma-elish/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Question struct {
	*gin.Engine
	c *config.Config
	i *infra.Infra
	v *validator.Validate
}

func New(c *config.Config, i *infra.Infra, r *gin.Engine, v *validator.Validate) *Question {
	return &Question{
		c:      c,
		i:      i,
		Engine: r,
		v:      v,
	}
}

func (q *Question) Init() {
	r := repository.New(q.i.Postgres, q.i.Redis)
	s := service.New(q.c, r)
	h := handler.New(s, q.v)

	authMiddleware := middleware.Auth(q.c.JWT.Secret)

	v1 := q.Group("/api/v1/question").Use(authMiddleware)

	v1.POST("", h.CreateQuestion)
	v1.GET("", h.GetListQuestions)
	v1.GET("/:question_id", h.GetDetailQuestion)
	v1.PUT("/:question_id", h.UpdateQuestion)
	v1.DELETE("/:question_id", h.DeleteQuestion)
	v1.GET("/by-type", h.GetQuestionsByType)
}
