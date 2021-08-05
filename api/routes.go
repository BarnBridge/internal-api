package api

import (
	"net/http"

	"github.com/barnbridge/internal-api/response"
	"github.com/gin-gonic/gin"
)

func (a *API) setRoutes() {
	a.engine.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "It works!")
	})

	a.engine.GET("/health", func(ctx *gin.Context) {
		err := a.db.Ping(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status": http.StatusInternalServerError,
				"data":   "NOT OK",
			})

			return
		}

		response.OK(ctx, "OK")
	})
}
