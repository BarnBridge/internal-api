package smartyield

import (
	"fmt"
	"strings"

	"github.com/barnbridge/internal-api/smartyield/types"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/utils"
)

func (h *SmartYield) PoolLiquidity(ctx *gin.Context) {
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

	underlyingDecimals, err := h.PoolUnderlyingDecimals(ctx, poolAddress)
	if err != nil {
		response.Error(ctx, errors.Wrap(err, "could not find smartyield pool"))
		return
	}
	tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(underlyingDecimals))

	query := fmt.Sprintf(`
		select date_trunc(%s, to_timestamp(block_timestamp)) as scale,
			   avg(senior_liquidity) as senior_liquidity,
			   avg(junior_liquidity) as junior_liquidity
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

	var points []types.LiquidityPoint
	for rows.Next() {
		var p types.LiquidityPoint
		err := rows.Scan(&p.Point, &p.SeniorLiquidity, &p.JuniorLiquidity)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		p.JuniorLiquidity = p.JuniorLiquidity.DivRound(tenPowDec, int32(underlyingDecimals))
		p.SeniorLiquidity = p.SeniorLiquidity.DivRound(tenPowDec, int32(underlyingDecimals))

		points = append(points, p)
	}

	response.OK(ctx, points)
}
