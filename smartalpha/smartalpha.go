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
	smartalpha.GET("/pools/:poolAddress/tokens-price-chart", s.TokensPriceChart)
	smartalpha.GET("/pools/:poolAddress/pool-performance-chart", s.poolPerformanceChart)
	smartalpha.GET("/pools/:poolAddress/previous-epochs", s.poolPreviousEpochs)

	smartalpha.GET("/users/:address/portfolio-value", s.UserPortfolioValue)
	smartalpha.GET("/users/:address/queue-positions", s.UserQueuePositions)

	smartalpha.GET("/transactions", s.transactions)

	smartalpha.GET("/rewards/pools", s.RewardPools)
	smartalpha.GET("/rewards/pools/:poolAddress/transactions", s.RewardPoolTransactions)
}
