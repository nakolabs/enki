package handler

import (
	"enuma-elish/internal/class/service"
	"enuma-elish/internal/class/service/data/request"
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

func (h *Handler) CreateClass(c *gin.Context) {
	data := request.CreateClassRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.CreateClass(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("create class success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetDetailClass(c *gin.Context) {
	classIDStr := c.Param("class_id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	classDetail, err := h.service.GetDetailClass(c.Request.Context(), classID)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get detail class success").
		SetData(classDetail)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ListClass(c *gin.Context) {
	httpQuery := request.GetListClassQuery{}
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

	data, meta, err := h.service.GetListClass(c.Request.Context(), httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get list class error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get list class success").
		SetData(data).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateClass(c *gin.Context) {
	classIDStr := c.Param("class_id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data := request.UpdateClassRequest{}
	err = c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.UpdateClass(c.Request.Context(), classID, data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("update class success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteClass(c *gin.Context) {
	classIDStr := c.Param("class_id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	err = h.service.DeleteClass(c.Request.Context(), classID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("delete class error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("delete class success")

	c.JSON(http.StatusOK, response)
}

func (h *Handler) AddStudentsToClass(c *gin.Context) {
	data := request.AddStudentToClassRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.AddStudentsToClass(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("students added to class successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) AddTeachersToClass(c *gin.Context) {
	data := request.AddTeacherToClassRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.AddTeachersToClass(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("teachers added to class successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) AddSubjectsToClass(c *gin.Context) {
	data := request.AddSubjectToClassRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.AddSubjectsToClass(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("subjects added to class successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetStudentsByClass(c *gin.Context) {
	classIDStr := c.Param("class_id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	httpQuery := request.GetStudentsByClassQuery{}
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

	data, meta, err := h.service.GetStudentsByClass(c.Request.Context(), classID, httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get students by class error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get students by class success").
		SetData(data).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetTeachersByClass(c *gin.Context) {
	classIDStr := c.Param("class_id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	httpQuery := request.GetTeachersByClassQuery{}
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

	data, meta, err := h.service.GetTeachersByClass(c.Request.Context(), classID, httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get teachers by class error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get teachers by class success").
		SetData(data).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetSubjectsByClass(c *gin.Context) {
	classIDStr := c.Param("class_id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	httpQuery := request.GetSubjectsByClassQuery{}
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

	data, meta, err := h.service.GetSubjectsByClass(c.Request.Context(), classID, httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get subjects by class error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get subjects by class success").
		SetData(data).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) RemoveTeachersFromClass(c *gin.Context) {
	data := request.RemoveTeacherFromClassRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.RemoveTeachersFromClass(c.Request.Context(), data)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("remove teachers from class error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("teachers removed from class successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) RemoveStudentsFromClass(c *gin.Context) {
	data := request.RemoveStudentFromClassRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.RemoveStudentsFromClass(c.Request.Context(), data)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("remove students from class error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("students removed from class successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) RemoveSubjectsFromClass(c *gin.Context) {
	data := request.RemoveSubjectFromClassRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.RemoveSubjectsFromClass(c.Request.Context(), data)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("remove subjects from class error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("subjects removed from class successfully").
		SetData(data)

	c.JSON(http.StatusOK, response)
}
