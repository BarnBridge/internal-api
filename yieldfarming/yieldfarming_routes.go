package yieldfarming

import (
	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/db"
)

type YieldFarming struct {
	db *db.DB
}

func New(db *db.DB) *YieldFarming {
	return &YieldFarming{db: db}
}

func (h *YieldFarming) SetRoutes(engine *gin.Engine) {
	yieldFarming := engine.Group("/api/yieldfarming")
	yieldFarming.GET("/staking-actions/list", h.StakingActionsList)
	yieldFarming.GET("/staking-actions/chart", h.StakingActionsChart)
}
