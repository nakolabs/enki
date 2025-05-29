package handler

import (
	"enuma-elish/internal/teacher/service"
	"enuma-elish/internal/teacher/service/data/request"
	commonHttp "enuma-elish/pkg/http"
	"errors"
	"fmt"
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

	httpQuery := request.GetListTeacherQuery{}
	httpQuery.Query = commonHttp.DefaultQuery()
	err := c.BindQuery(&httpQuery)
	if err != nil {
		fmt.Println(err)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data, meta, err := h.service.ListTeachers(c.Request.Context(), httpQuery)
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

func (h *Handler) DeleteTeacher(c *gin.Context) {
	teacherIDStr := c.Param("teacher_id")
	teacherID, err := uuid.Parse(teacherIDStr)
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

	err = h.service.DeleteTeacher(c.Request.Context(), teacherID, schoolID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("delete teacher error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("delete teacher success")

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetDetailTeacher(c *gin.Context) {
	teacherIDStr := c.Param("teacher_id")
	teacherID, err := uuid.Parse(teacherIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data, err := h.service.GetDetailTeacher(c.Request.Context(), teacherID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get detail teacher error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get detail teacher success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateTeacherClass(c *gin.Context) {
	data := request.UpdateTeacherClassRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.UpdateTeacherClass(c.Request.Context(), data)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("update teacher class error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("teacher class updated successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetTeacherSubjects(c *gin.Context) {
	teacherIDStr := c.Param("teacher_id")
	teacherID, err := uuid.Parse(teacherIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	subjects, err := h.service.GetTeacherSubjects(c.Request.Context(), teacherID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get teacher subjects error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get teacher subjects success").
		SetData(subjects)

	c.JSON(http.StatusOK, response)
}
