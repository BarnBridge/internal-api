package smartexposure

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartexposure/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartExposure) tokensPricesChart(ctx *gin.Context) {
	tranche := ctx.Param("eTokenAddress")
	trancheAddress, err := utils.ValidateAccount(tranche)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	err, exists := s.checkTrancheExists(ctx, trancheAddress)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	if !exists {
		response.NotFound(ctx)
		return
	}

	pool, err := s.getPoolByETokenAddress(ctx, trancheAddress)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	window := strings.ToLower(ctx.DefaultQuery("window", "30d"))
	_, dateTrunc, err := validateWindow(window)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	dateTrunc = strings.Replace(dateTrunc, "'", "", -1)
	startTs, err := getStartDate(window)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	rows, err := s.db.Connection().Query(ctx, `select a.point, a.token_price as token_a_price,b.token_price as token_b_price from smart_exposure.get_token_price_chart($1,$3,$4) a
									inner join smart_exposure.get_token_price_chart($2,$3,$4) b on a.point = b.point`, pool.TokenA.Address, pool.TokenB.Address, startTs, dateTrunc)
	if err != nil && err != sql.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var points []types.PriceTrendPoint
	for rows.Next() {
		var p types.PriceTrendPoint
		err := rows.Scan(&p.Point, &p.TokenAPrice, &p.TokenBPrice)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		points = append(points, p)
	}

	response.OK(ctx, points)
}
