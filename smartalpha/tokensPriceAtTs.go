package smartalpha

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartalpha/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartAlpha) TokensPriceAtTs(ctx *gin.Context) {
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

	timestampString := strings.ToLower(ctx.DefaultQuery("timestamp", "0"))
	ts, err := strconv.ParseInt(timestampString, 10, 64)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	if ts == 0 {
		err, ts = s.getLastPoolStateTimestamp(ctx, poolAddress)
		if err != nil {
			response.Error(ctx, err)
			return
		}

	}

	var p types.PoolTokensPrice

	err = s.db.Connection().QueryRow(ctx, `
			select t.estimated_junior_token_price as junior_price, 
				   t.estimated_senior_token_price as senior_price
			from smart_alpha.pool_state t
			where pool_address = $1
			  and t.block_timestamp <= $2
			order by block_timestamp desc
			limit 1`, poolAddress, ts).Scan(&p.JuniorTokenPrice, &p.SeniorTokenPrice)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	p.JuniorTokenPrice = p.JuniorTokenPrice.Shift(-18)
	p.SeniorTokenPrice = p.SeniorTokenPrice.Shift(-18)
	p.Timestamp = time.Unix(ts, 0)

	response.OK(ctx, p)
}
