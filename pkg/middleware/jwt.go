package middleware

import (
	"context"
	commonHttp "enuma-elish/pkg/http"
	"enuma-elish/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := commonHttp.NewResponse().
				SetCode(http.StatusUnauthorized).
				SetMessage("Unauthorized").
				SetErrors([]string{"Authorization header is empty"})

			c.JSON(response.Code, response)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response := commonHttp.NewResponse().
				SetCode(http.StatusUnauthorized).
				SetMessage("Unauthorized").
				SetErrors([]string{"Authorization format must be Bearer {token}"})

			c.JSON(response.Code, response)
			c.Abort()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Verify(tokenStr, secret)
		if err != nil {
			response := commonHttp.NewResponse().
				SetCode(http.StatusUnauthorized).
				SetMessage("Unauthorized").
				SetErrors([]string{err.Error()})

			c.JSON(response.Code, response)
			c.Abort()
			return
		}

		claim, err := jwt.ExtractToken(*token)
		if err != nil {
			response := commonHttp.NewResponse().
				SetCode(http.StatusUnauthorized).
				SetMessage("Unauthorized").
				SetErrors([]string{err.Error()})

			c.JSON(response.Code, response)
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), jwt.ContextKey, claim)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
