package handler

import (
	"enuma-elish/internal/school/service"
	"enuma-elish/internal/school/service/data/request"
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

func (h *Handler) CreateSchool(c *gin.Context) {
	data := request.CreateSchoolRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.CreatSchool(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("create school success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetDetailSchool(c *gin.Context) {
	schoolIDStr := c.Param("school_id")
	schoolID, err := uuid.Parse(schoolIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	schoolDetail, err := h.service.GetDetailSchool(c.Request.Context(), schoolID)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get detail school success").
		SetData(schoolDetail)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ListSchool(c *gin.Context) {
	httpQuery := request.GetListSchoolQuery{}
	httpQuery.Query = commonHttp.DefaultQuery()
	err := c.BindQuery(&httpQuery)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	data, meta, err := h.service.GetListSchool(c.Request.Context(), httpQuery)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("get list school error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("get list school success").
		SetData(data).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) SwitchSchool(c *gin.Context) {
	schoolIDStr := c.Param("school_id")
	schoolID, err := uuid.Parse(schoolIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	token, err := h.service.SwitchSchool(c.Request.Context(), schoolID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage(err.Error()).
			SetErrors(err)
		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("switch school success").
		SetData(gin.H{"access_token": token})

	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteSchool(c *gin.Context) {
	schoolIDStr := c.Param("school_id")
	schoolID, err := uuid.Parse(schoolIDStr)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	err = h.service.DeleteSchool(c.Request.Context(), schoolID)
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage("delete school error").
			SetErrors([]error{err})

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("delete school success")

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateSchoolProfile(c *gin.Context) {
	schoolIDParam := c.Param("school_id")
	schoolID, err := uuid.Parse(schoolIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid school ID"})
		return
	}

	var req request.UpdateSchoolProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	school, err := h.service.UpdateSchoolProfile(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "School profile updated successfully",
		"data":    school,
	})
}
