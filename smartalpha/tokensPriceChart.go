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

func (s *SmartAlpha) TokensPriceChart(ctx *gin.Context) {
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
	window, dateTrunc, err := validateWindow(window)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	query := fmt.Sprintf(`
			select date_trunc(%s, to_timestamp(block_timestamp)) as point,
				   avg(t.estimated_junior_token_price)           as junior_price,
				   avg(t.estimated_senior_token_price)           as senior_price
			from smart_alpha.pool_state t
			where pool_address = $1
			  and to_timestamp(block_timestamp) > now() - interval % s
			group by point
			order by point asc;`, dateTrunc, window)

	rows, err := s.db.Connection().Query(ctx, query, poolAddress)
	if err != nil && err != sql.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var points []types.TokensPricePoint
	for rows.Next() {
		var p types.TokensPricePoint
		err := rows.Scan(&p.Point, &p.JuniorTokenPrice, &p.SeniorTokenPrice)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		p.SeniorTokenPrice = p.SeniorTokenPrice.Shift(-18)
		p.JuniorTokenPrice = p.JuniorTokenPrice.Shift(-18)
		points = append(points, p)
	}

	response.OK(ctx, points)
}
