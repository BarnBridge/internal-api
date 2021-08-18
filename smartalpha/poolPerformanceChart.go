package smartalpha

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartalpha/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartAlpha) poolPerformanceChart(ctx *gin.Context) {
	poolAddress := ctx.Param("poolAddress")
	if poolAddress != "" {
		poolAddress, err := utils.ValidateAccount(poolAddress)

		if err != nil {
			response.Error(ctx, err)
			return
		}
		err, exists := s.checkPoolExists(ctx, poolAddress)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		if !exists {
			response.NotFound(ctx)
			return
		}
	}

	window := strings.ToLower(ctx.DefaultQuery("window", "30d"))
	totalPoints := getTotalPoints(window)
	window, _, err := validateWindow(window)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	query := fmt.Sprintf(`
				select to_timestamp(ts)                                                                                          as point,
					   coalesce(senior_without_sa, 0)                                                                            as senior_without_sa,
					   coalesce(senior_with_sa, 0)                                                                               as senior_with_sa,
					   coalesce(junior_without_sa, 0)                                                                            as junior_without_sa,
					   coalesce(junior_with_sa, 0)                                                                               as junior_with_sa,
					   coalesce(token_price_at_ts(( select pool_token_address
													from smart_alpha.pools
													where pool_address = $1 ),
												  ( select oracle_asset_symbol
													from smart_alpha.pools
													where pool_address = $1 ), ts),
								0)                                                                                               as pool_token_price
				from generate_series((select extract(epoch from now() - interval % s)::bigint),
												 (select extract(epoch from now()))::bigint, %s) as ts
						 inner join smart_alpha.performance_at_ts($1, ts) on true
				order by ts;`, window, totalPoints)

	rows, err := s.db.Connection().Query(ctx, query, poolAddress)
	if err != nil && err != sql.ErrNoRows {
		response.Error(ctx, err)
		return
	}
	var points []types.PerformancePoint

	for rows.Next() {
		var p types.PerformancePoint
		err := rows.Scan(&p.Point, &p.SeniorWithoutSA, &p.SeniorWithSA, &p.JuniorWithoutSA, &p.JuniorWithSA, &p.UnderlyingPrice)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		points = append(points, p)
	}

	response.OK(ctx, points)
}
