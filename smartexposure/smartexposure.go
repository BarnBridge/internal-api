package smartexposure

import (
	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/db"
)

type SmartExposure struct {
	db *db.DB
}

func New(db *db.DB) *SmartExposure {
	return &SmartExposure{db: db}
}

func (s *SmartExposure) SetRoutes(engine *gin.Engine) {
	smartExposure := engine.Group("/api/smartexposure")

	smartExposure.GET("/pools", s.sePools)

	smartExposure.GET("/tranches", s.allTranches)
	smartExposure.GET("/tranches/:eTokenAddress", s.trancheDetails)
	smartExposure.GET("/tranches/:eTokenAddress/etoken-price", s.eTokenPriceChart)
	smartExposure.GET("/tranches/:eTokenAddress/price-trend", s.tokensPricesChart)
	smartExposure.GET("/tranches/:eTokenAddress/liquidity", s.trancheLiquidityChart)
	smartExposure.GET("/tranches/:eTokenAddress/ratio-deviation", s.ratioDeviationChart)

	smartExposure.GET("/transactions", s.transactions)
	smartExposure.GET("/users/:userAddress/portfolio-value", s.userPortfolioValueChart)
}
