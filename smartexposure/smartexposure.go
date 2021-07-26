package smartexposure

import (
	"github.com/barnbridge/internal-api/db"
	"github.com/gin-gonic/gin"
)

type SmartExposure struct {
	db *db.DB
}

func New(db *db.DB) *SmartExposure {
	return &SmartExposure{db: db}
}

func (s *SmartExposure) SetRoutes(engine *gin.Engine) {
	smartExposure := engine.Group("/api/smartexposure")
	smartExposure.GET("/tranches", s.handleAllSEPoolsTranches)
	smartExposure.GET("/tranches/:eTokenAddress", s.handleTrancheDetails)
}
