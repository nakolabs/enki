package auth

import (
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/internal/auth/handler"
	"enuma-elish/internal/auth/repository"
	"enuma-elish/internal/auth/service"
	"enuma-elish/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Auth struct {
	*gin.Engine
	c *config.Config
	i *infra.Infra
	v *validator.Validate
}

func New(c *config.Config, i *infra.Infra, r *gin.Engine, v *validator.Validate) *Auth {
	return &Auth{
		c:      c,
		i:      i,
		Engine: r,
		v:      v,
	}
}

func (a *Auth) Init() {
	r := repository.New(a.i.Postgres, a.i.Redis)
	s := service.New(r, a.c)
	h := handler.New(s, a.v)

	v1 := a.Group("/api/v1/auth")
	v1.POST("/register", h.Register)
	v1.POST("/login", h.Login)
	v1.POST("/register/verify-email", h.VerifyEmail)
	v1.POST("/forgot-password", h.ForgotPassword)
	v1.POST("/forgot-password/verify", h.ForgotPasswordVerify)
	v1.POST("/refresh-token", h.RefreshToken)

	v1.GET("/me", middleware.Auth(a.c.JWT.Secret), h.Me)
	v1.GET("/profile", middleware.Auth(a.c.JWT.Secret), h.Profile)
	v1.PUT("/profile", middleware.Auth(a.c.JWT.Secret), h.UpdateProfile)
}
