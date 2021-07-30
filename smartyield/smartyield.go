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

	sy.GET("/rewards/pools", h.RewardPools)
	sy.GET("/rewards/pools/:poolAddress/transactions", h.RewardPoolsStakingActions)

	sy.GET("/rewards/v2/pools", h.RewardPoolsV2)
	sy.GET("/rewards/v2/pools/:poolAddress/transactions", h.RewardPoolsStakingActions)

	sy.GET("/pools/:address/apy", h.PoolAPYTrend)
	sy.GET("/pools/:address/liquidity", h.PoolLiquidity)
	sy.GET("/pools/:address/transactions", h.PoolTransactions)
	sy.GET("/pools/:address/senior-bonds", h.PoolSeniorBonds)
	sy.GET("/pools/:address/junior-bonds", h.PoolJuniorBonds)

	sy.GET("/users/:address/history", h.UserTransactionHistory)
	sy.GET("/users/:address/redeems/senior", h.UserSeniorRedeems)
	sy.GET("/users/:address/junior-past-positions", h.JuniorPastPositions)
	sy.GET("/users/:address/portfolio-value", h.UserPortfolioValue)
	sy.GET("/users/:address/portfolio-value/junior", h.UserJuniorPortfolioValue)
	sy.GET("/users/:address/portfolio-value/senior", h.UserSeniorPortfolioValue)
}
