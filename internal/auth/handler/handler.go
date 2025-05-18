package handler

import (
	"enuma-elish/internal/auth/service"
	"enuma-elish/internal/auth/service/data/request"
	commonHttp "enuma-elish/pkg/http"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
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

func (h *Handler) Register(c *gin.Context) {
	data := request.Register{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		log.Err(err).Msg("validator fail")
		c.Error(err)
		return
	}

	err = h.service.Register(c.Request.Context(), data)
	if err != nil {
		c.Error(err)
		return
	}

	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("registration success").
		SetData(data)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Login(c *gin.Context) {
	req := request.LoginRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("login success").
		SetData(data)
	c.JSON(http.StatusOK, response)
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	req := request.VerifyEmailRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.Error(err)
		return
	}
	err = h.service.VerifyEmail(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	response := commonHttp.NewResponse().
		SetCode(http.StatusOK).
		SetMessage("verify email success")
	c.JSON(http.StatusOK, response)
}
