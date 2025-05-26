package handler

import (
	"enuma-elish/internal/student/service"
	"enuma-elish/internal/student/service/data/request"
	commonHttp "enuma-elish/pkg/http"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Handler struct {
	service   service.Service
	validator *validator.Validate
}

func New(service service.Service, validator *validator.Validate) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

func (h *Handler) InviteStudent(c *gin.Context) {
	data := request.InviteStudentRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.InviteStudent(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("invite teacher success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) VerifyStudentEmail(c *gin.Context) {
	data := request.VerifyStudentEmailRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.VerifyStudentEmail(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("verify teacher email success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateStudentAfterInvite(c *gin.Context) {
	data := request.UpdateStudentAfterVerifyEmailRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.UpdateStudentAfterVerifyEmail(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("update student success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ListStudent(c *gin.Context) {
	httpQuery := request.GetListStudentQuery{}
	httpQuery.Query = commonHttp.DefaultQuery()
	err := c.BindQuery(&httpQuery)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data, meta, err := h.service.GetListStudent(c.Request.Context(), httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get list teacher error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get list teacher success").
		SetData(data).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteStudent(c *gin.Context) {
	studentIDStr := c.Param("student_id")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	schoolIDStr := c.Query("school_id")
	if schoolIDStr == "" {
		c.Error(errors.New("school_id is required")).SetType(gin.ErrorTypeBind)
		return
	}

	schoolID, err := uuid.Parse(schoolIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	err = h.service.DeleteStudent(c.Request.Context(), studentID, schoolID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("delete student error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("delete student success")

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetDetailStudent(c *gin.Context) {
	studentIDStr := c.Param("student_id")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data, err := h.service.GetDetailStudent(c.Request.Context(), studentID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get detail student error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get detail student success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateStudentClass(c *gin.Context) {
	data := request.UpdateStudentClassRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.UpdateStudentClass(c.Request.Context(), data)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("update student class error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("student class updated successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}
