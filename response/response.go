package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/db"
	"github.com/barnbridge/internal-api/utils"
)

func Error(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"status": http.StatusInternalServerError,
		"data":   err.Error(),
	})
}

func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, map[string]interface{}{
		"status": http.StatusBadRequest,
		"data":   err.Error(),
	})
}

func OK(c *gin.Context, data interface{}, meta ...interface{}) {
	resp := map[string]interface{}{
		"status": http.StatusOK,
		"data":   data,
	}

	if len(meta) > 0 {
		resp["meta"] = meta[0]
	}

	c.JSON(http.StatusOK, resp)
}

func OKWithBlock(c *gin.Context, db *db.DB, data interface{}, meta ...interface{}) {
	block, err := utils.GetHighestBlock(c, db)
	if err != nil {
		Error(c, err)
		return
	}

	if meta == nil || len(meta) == 0 {
		OK(c, data, Meta().Set("block", block))
	} else {

		if m, ok := meta[0].(map[string]interface{}); ok {
			m["block"] = block
			OK(c, data, meta...)
			return
		}

		if m, ok := meta[0].(ResponseMeta); ok {
			m.Set("block", block)
			OK(c, data, meta...)
			return
		}
	}
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusNotFound,
		"data":   nil,
	})
}
