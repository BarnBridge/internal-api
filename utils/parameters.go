package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetQueryLimit(c *gin.Context) (int64, error) {
	limit := c.DefaultQuery("limit", "10")

	l, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "invalid 'limit' parameter")
	}

	return l, nil
}

func GetQueryPage(c *gin.Context) (int64, error) {
	page := c.DefaultQuery("page", "1")

	p, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid 'page' parameter")
	}

	return p, nil
}
