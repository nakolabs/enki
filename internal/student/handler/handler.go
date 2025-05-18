package handler

import (
	"enuma-elish/internal/student/service"
	"enuma-elish/internal/student/service/data/request"
	commonHttp "enuma-elish/pkg/http"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
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

	err = h.service.VerifyTeacherEmail(c.Request.Context(), data)
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

	err = h.service.UpdateTeacherAfterVerifyEmail(c.Request.Context(), data)
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
