package handler

import (
	"enuma-elish/internal/school/service"
	"enuma-elish/internal/school/service/data/request"
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
	data, err := h.service.GetListSchool(c.Request.Context())
	if err != nil {
		response := commonHttp.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetMessage(err.Error()).
			SetErrors(err.Error())

		c.JSON(response.Code, response)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("list school success").
		SetData(data)
	
	c.JSON(http.StatusOK, response)
}
