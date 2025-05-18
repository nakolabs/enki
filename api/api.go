package api

import (
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/internal/auth"
	"enuma-elish/internal/school"
	"enuma-elish/internal/student"
	"enuma-elish/internal/teacher"
	"enuma-elish/pkg/middleware"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net/http"
	"time"
)

type API struct {
	*gin.Engine
	config *config.Config
	infra  *infra.Infra
}

func New(c *config.Config, infra *infra.Infra) *API {
	g := gin.New()

	api := &API{g, c, infra}
	validate := validator.New()
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
	api.Use(gin.Recovery(), middleware.ErrorParser(), corsMiddleware)

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	auth.New(api.config, api.infra, api.Engine, validate).Init()
	school.New(api.config, api.infra, api.Engine, validate).Init()
	teacher.New(api.config, api.infra, api.Engine, validate).Init()
	student.New(api.config, api.infra, api.Engine, validate).Init()

	return api
}

func (api *API) Run() {
	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", api.config.Http.Host, api.config.Http.Port),
		Handler:      api,
		ReadTimeout:  time.Duration(api.config.Http.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(api.config.Http.WriteTimeout) * time.Second,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("server error")
		panic(err)
	}
}

func (api *API) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "api stuff",
	}

	run := &cobra.Command{
		Use:   "run",
		Short: "run api",
		Run: func(cmd *cobra.Command, args []string) {
			api.Run()
		},
	}

	cmd.AddCommand(run)
	return cmd
}
