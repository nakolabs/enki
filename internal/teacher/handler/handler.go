package handler

import (
	"enuma-elish/internal/teacher/service"
	"enuma-elish/internal/teacher/service/data/request"
	commonHttp "enuma-elish/pkg/http"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
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

func (h *Handler) InviteTeacher(c *gin.Context) {
	data := request.InviteTeacherRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.InviteTeacher(c.Request.Context(), data)
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

func (h *Handler) VerifyTeacherEmail(c *gin.Context) {
	data := request.VerifyTeacherEmailRequest{}
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

func (h *Handler) UpdateTeacherAfterInvite(c *gin.Context) {
	data := request.UpdateTeacherAfterVerifyEmailRequest{}
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
		SetMessage("update teacher success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ListTeachers(c *gin.Context) {
	schoolIdStr := c.Param("school_id")
	schoolID, err := uuid.Parse(schoolIdStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	httpQuery := commonHttp.DefaultQuery()
	err = c.BindQuery(&httpQuery)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data, meta, err := h.service.ListTeachers(c.Request.Context(), schoolID, httpQuery)
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
