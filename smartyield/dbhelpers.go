package smartyield

import (
	"github.com/gin-gonic/gin"
)

func (h *SmartYield) PoolUnderlyingDecimals(ctx *gin.Context, poolAddress string) (int64, error) {
	var ud int64
	err := h.db.Connection().QueryRow(
		ctx,
		`	select underlying_decimals 
					from smart_yield.pools p
				where p.pool_address = $1`,
		poolAddress,
	).Scan(&ud)
	return ud, err
}
