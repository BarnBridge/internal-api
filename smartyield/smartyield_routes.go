package smartyield

import (
	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/db"
)

type SmartYield struct {
	db *db.DB
}

func New(db *db.DB) *SmartYield {
	return &SmartYield{db: db}
}

func (h *SmartYield) SetRoutes(engine *gin.Engine) {
	sy := engine.Group("/api/smartyield")
	sy.GET("/pools", h.Pools)
	sy.GET("/pools/:address", h.PoolDetails)

	// smartYield.GET("/pools/:address/apy", a.handlePoolAPYTrend)
	// smartYield.GET("/pools/:address/liquidity", a.handlePoolLiquidity)
	// smartYield.GET("/pools/:address/transactions", a.handlePoolTransactions)
	// smartYield.GET("/pools/:address/senior-bonds", a.handlePoolSBonds)
	// smartYield.GET("/pools/:address/junior-bonds", a.handlePoolJBonds)
	// smartYield.GET("/users/:address/history", a.handleSYUserTransactionHistory)
	// smartYield.GET("/users/:address/redeems/senior", a.handleSeniorRedeems)
	// smartYield.GET("/users/:address/junior-past-positions", a.handleJuniorPastPositions)
	// smartYield.GET("/users/:address/portfolio-value", a.handleSYUserPortfolioValue)
	// smartYield.GET("/users/:address/portfolio-value/junior", a.handleSYUserJuniorPortfolioValue)
	// smartYield.GET("/users/:address/portfolio-value/senior", a.handleSYUserSeniorPortfolioValue)
	// smartYield.GET("/rewards/pools", a.handleRewardPools)
	// smartYield.GET("/rewards/pools/:poolAddress/transactions", a.handleRewardPoolsStakingActions)
	// smartYield.GET("/rewards/v2/pools", a.handleRewardPoolsV2)
	// smartYield.GET("/rewards/v2/pools/:poolAddress/transactions", a.handleRewardPoolsStakingActions)

}


