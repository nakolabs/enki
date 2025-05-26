package handler

import (
	"enuma-elish/internal/ppdb/service"
	"enuma-elish/internal/ppdb/service/data/request"
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

func (h *Handler) CreatePPDB(c *gin.Context) {
	data := request.CreatePPDBRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.CreatePPDB(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusCreated).
		SetMessage("PPDB created successfully")

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) UpdatePPDB(c *gin.Context) {
	data := request.UpdatePPDBRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	ppdbID, err := uuid.Parse(c.Param("ppdb_id"))
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	data.ID = ppdbID

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.UpdatePPDB(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("PPDB updated successfully")

	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeletePPDB(c *gin.Context) {
	ppdbID, err := uuid.Parse(c.Param("ppdb_id"))
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	err = h.service.DeletePPDB(c.Request.Context(), ppdbID)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("PPDB deleted successfully")

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetPPDBByID(c *gin.Context) {
	ppdbID, err := uuid.Parse(c.Param("ppdb_id"))
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	ppdb, err := h.service.GetPPDBByID(c.Request.Context(), ppdbID)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("PPDB retrieved successfully").
		SetData(ppdb)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetListPPDB(c *gin.Context) {
	query := request.GetListPPDBQuery{}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	ppdbs, meta, err := h.service.GetListPPDB(c.Request.Context(), query)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("PPDB list retrieved successfully").
		SetData(ppdbs).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) RegisterPPDB(c *gin.Context) {
	data := request.RegisterPPDBRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	result, err := h.service.RegisterPPDB(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusCreated).
		SetMessage("Successfully registered for PPDB").
		SetData(result)

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetPPDBRegistrants(c *gin.Context) {
	query := request.GetPPDBRegistrantsQuery{}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	registrants, meta, err := h.service.GetPPDBRegistrants(c.Request.Context(), query)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("PPDB registrants retrieved successfully").
		SetData(registrants).
		SetMeta(meta)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) SelectPPDBStudents(c *gin.Context) {
	data := request.PPDBSelectionRequest{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		c.Error(err)
		return
	}

	err = h.service.SelectPPDBStudents(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("PPDB selection completed successfully")

	c.JSON(http.StatusOK, response)
}
