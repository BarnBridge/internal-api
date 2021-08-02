package smartyield

import (
	"fmt"
	"strings"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartyield/types"
	"github.com/barnbridge/internal-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (h *SmartYield) PoolAPYTrend(ctx *gin.Context) {
	pool := ctx.Param("address")

	poolAddress, err := utils.ValidateAccount(pool)
	if err != nil {
		response.BadRequest(ctx, errors.New("invalid pool address"))
		return
	}

	window := strings.ToLower(ctx.DefaultQuery("window", "30d"))
	window, dateTrunc, err := validateWindow(window)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	query := fmt.Sprintf(`
		select date_trunc(%s, to_timestamp(block_timestamp)) as scale,
			   avg(senior_apy) as senior_apy,
			   avg(junior_apy) as junior_apy,
		       avg(originator_net_apy) as originator_net_apy
		from smart_yield.pool_state
		where pool_address = $1
		and to_timestamp(block_timestamp) > now() - interval %s
		group by scale
		order by scale;`, dateTrunc, window)

	rows, err := h.db.Connection().Query(ctx, query, poolAddress)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var points []types.APYTrendPoint
	for rows.Next() {
		var p types.APYTrendPoint

		err := rows.Scan(&p.Point, &p.SeniorAPY, &p.JuniorAPY, &p.OriginatorNetAPY)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		points = append(points, p)
	}

	response.OK(ctx, points)
}
