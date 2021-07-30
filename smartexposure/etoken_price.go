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

func (s *SmartExposure) eTokenPriceChart(ctx *gin.Context) {
	tranche := ctx.Param("eTokenAddress")
	eTokenAddress, err := utils.ValidateAccount(tranche)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	err, exists := s.checkTrancheExists(ctx, eTokenAddress)
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
		response.Error(ctx, err)
		return
	}

	query := fmt.Sprintf(`
		select date_trunc(%s, to_timestamp(block_timestamp)) as point,
			   avg(etoken_price)
		from smart_exposure.tranche_state
		where etoken_address = $1
		and to_timestamp(block_timestamp) > now() - interval %s
		group by point
		order by point;`, dateTrunc, window)

	rows, err := s.db.Connection().Query(ctx, query, eTokenAddress)
	if err != nil && err != sql.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var points []types.ETokenPricePoint
	for rows.Next() {
		var p types.ETokenPricePoint
		err := rows.Scan(&p.Point, &p.ETokenPrice)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		points = append(points, p)
	}

	response.OK(ctx, points)
}
