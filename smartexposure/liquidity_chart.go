package smartexposure

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartexposure/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartExposure) trancheLiquidityChart(ctx *gin.Context) {
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

	window := strings.ToLower(ctx.DefaultQuery("window", "30d"))
	window, dateTrunc, err := validateWindow(window)
	if err != nil {
		response.NotFound(ctx)
		return
	}

	query := fmt.Sprintf(`
			select date_trunc(%s, to_timestamp(block_timestamp)) as point,
				   avg(token_a_liquidity) ,
				   avg(token_b_liquidity)
			from smart_exposure.tranche_state
			where etoken_address = $1
			  and to_timestamp(block_timestamp) > now() - interval %s
			group by point
			order by point;`, dateTrunc, window)

	rows, err := s.db.Connection().Query(ctx, query, trancheAddress)
	if err != nil && err != sql.ErrNoRows {
		response.NotFound(ctx)
		return
	}

	var points []types.LiquidityPoint
	for rows.Next() {
		var p types.LiquidityPoint
		err := rows.Scan(&p.Point, &p.TokenALiquidity, &p.TokenBLiquidity)
		if err != nil {
			response.NotFound(ctx)
			return
		}
		points = append(points, p)
	}

	response.OK(ctx, points)
}
