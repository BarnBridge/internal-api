package smartalpha

import (
	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/db"
)

type SmartAlpha struct {
	db *db.DB
}

func New(db *db.DB) *SmartAlpha {
	return &SmartAlpha{db: db}
}

func (s *SmartAlpha) SetRoutes(engine *gin.Engine) {
	smartalpha := engine.Group("/api/smartalpha")
	smartalpha.GET("/pools", s.Pools)
}
