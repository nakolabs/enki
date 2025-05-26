package handler

import (
	"enuma-elish/internal/question/service"
	"enuma-elish/internal/question/service/data/request"
	commonHttp "enuma-elish/pkg/http"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Handler struct {
	service   service.Service
	validator *validator.Validate
}

func New(service service.Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:   service,
		validator: validate,
	}
}

func (h *Handler) CreateQuestion(c *gin.Context) {
	data := request.CreateQuestionRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	if err := data.Validate(); err != nil {
		c.Error(err)
		return
	}

	err = h.service.CreateQuestion(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusCreated).
		SetMessage("question created successfully").
		SetData(data)

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetDetailQuestion(c *gin.Context) {
	questionIDStr := c.Param("question_id")
	questionID, err := uuid.Parse(questionIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data, err := h.service.GetDetailQuestion(c.Request.Context(), questionID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get detail question error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get detail question success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetListQuestions(c *gin.Context) {
	httpQuery := request.GetListQuestionQuery{}
	httpQuery.Query = commonHttp.DefaultQuery()
	err := c.BindQuery(&httpQuery)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data, meta, err := h.service.GetListQuestions(c.Request.Context(), httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get list questions error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get list questions success").
		SetData(data).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateQuestion(c *gin.Context) {
	questionIDStr := c.Param("question_id")
	questionID, err := uuid.Parse(questionIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data := request.UpdateQuestionRequest{}
	err = c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	if err := data.Validate(); err != nil {
		c.Error(err)
		return
	}

	err = h.service.UpdateQuestion(c.Request.Context(), questionID, data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("question updated successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteQuestion(c *gin.Context) {
	questionIDStr := c.Param("question_id")
	questionID, err := uuid.Parse(questionIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	err = h.service.DeleteQuestion(c.Request.Context(), questionID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("delete question error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("question deleted successfully")

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetQuestionsByType(c *gin.Context) {
	httpQuery := request.GetQuestionsByTypeQuery{}
	err := c.BindQuery(&httpQuery)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(httpQuery); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetQuestionsByType(c.Request.Context(), httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get questions by type error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get questions by type success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}
