package api

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/barnbridge/internal-api/config"
	"github.com/barnbridge/internal-api/db"
)

var log = logrus.WithField("module", "api")

type Config struct {
	Port           string
	DevCorsEnabled bool
	DevCorsHost    string
}

type API struct {
	engine *gin.Engine
	db     *db.DB
}

func New(db *db.DB) *API {
	return &API{
		db: db,
	}
}

func (a *API) Run(ctx context.Context) {
	a.engine = gin.Default()

	if config.Store.API.DevCors {
		a.engine.Use(cors.New(cors.Config{
			AllowOrigins:     []string{config.Store.API.DevCorsHost},
			AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))
	}

	a.setRoutes()

	err := a.engine.Run(":" + config.Store.API.Port)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *API) Close() {
}
