package handler

import (
	"enuma-elish/internal/subject/service"
	"enuma-elish/internal/subject/service/data/request"
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

func New(service service.Service, validator *validator.Validate) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

func (h *Handler) CreateSubject(c *gin.Context) {
	data := request.CreateSubjectRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.CreateSubject(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("create subject success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetDetailSubject(c *gin.Context) {
	subjectIDStr := c.Param("subject_id")
	subjectID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	subjectDetail, err := h.service.GetDetailSubject(c.Request.Context(), subjectID)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get detail subject success").
		SetData(subjectDetail)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ListSubject(c *gin.Context) {
	httpQuery := request.GetListSubjectQuery{}
	httpQuery.Query = commonHttp.DefaultQuery()
	err := c.BindQuery(&httpQuery)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(httpQuery); err != nil {
		c.Error(err)
		return
	}

	data, meta, err := h.service.GetListSubject(c.Request.Context(), httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get list subject error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get list subject success").
		SetData(data).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateSubject(c *gin.Context) {
	subjectIDStr := c.Param("subject_id")
	subjectID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data := request.UpdateSubjectRequest{}
	err = c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.UpdateSubject(c.Request.Context(), subjectID, data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("update subject success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteSubject(c *gin.Context) {
	subjectIDStr := c.Param("subject_id")
	subjectID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	err = h.service.DeleteSubject(c.Request.Context(), subjectID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("delete subject error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("delete subject success")

	c.JSON(http.StatusOK, response)
}

func (h *Handler) AssignTeachersToSubject(c *gin.Context) {
	data := request.AssignTeacherToSubjectRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.AssignTeachersToSubject(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("teachers assigned to subject successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetTeachersBySubject(c *gin.Context) {
	subjectIDStr := c.Param("subject_id")
	subjectID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	httpQuery := request.GetTeachersBySubjectQuery{}
	httpQuery.Query = commonHttp.DefaultQuery()

	err = c.BindQuery(&httpQuery)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(httpQuery); err != nil {
		c.Error(err)
		return
	}

	data, meta, err := h.service.GetTeachersBySubject(c.Request.Context(), subjectID, httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get teachers by subject error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get teachers by subject success").
		SetData(data).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateSubjectClass(c *gin.Context) {
	data := request.UpdateSubjectClassRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.UpdateSubjectClass(c.Request.Context(), data)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("update subject class error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("subject class updated successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}
