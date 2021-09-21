package smartalpha

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

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
	startTs, endTs, totalPoints, err := s.poolPerformanceWindow(ctx, window, poolAddress)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	query := fmt.Sprintf(`
				select to_timestamp(ts)                                              as point,
					   coalesce(senior_without_sa, 0)                                as senior_without_sa,
					   coalesce(senior_with_sa, 0)                                   as senior_with_sa,
					   coalesce(junior_without_sa, 0)                                as junior_without_sa,
					   coalesce(junior_with_sa, 0)                                   as junior_with_sa,
					   coalesce(token_price_at_ts((select pool_token_address
												   from smart_alpha.pools
												   where pool_address = $3),
												  (select oracle_asset_symbol
												   from smart_alpha.pools
												   where pool_address = $3), ts), 0) as pool_token_price
				from generate_series($1,$2, %s) as ts
						 inner join smart_alpha.performance_at_ts($3, ts) on true
				order by ts;`, totalPoints)

	rows, err := s.db.Connection().Query(ctx, query, startTs, endTs, poolAddress)
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

func (s *SmartAlpha) poolPerformanceWindow(ctx context.Context, window string, poolAddress string) (int64, int64, string, error) {
	var startTs, endTs int64
	var duration time.Duration

	switch window {
	case "1w":
		duration = 7 * 24 * time.Hour
		startTs = time.Now().Add(-duration).Unix()
		endTs = time.Now().Unix()
		return startTs, endTs, "3*60*60", nil
	case "30d":
		duration = 30 * 24 * time.Hour
		startTs = time.Now().Add(-duration).Unix()
		endTs = time.Now().Unix()
		return startTs, endTs, "12*60*60", nil
	case "current":
		err := s.db.Connection().QueryRow(ctx, `select start_ts,end_ts from smart_alpha.get_epoch_ts($1,(select p.epoch_id from smart_alpha.pool_epoch_info p where p.pool_address = $1 order by p.epoch_id desc limit 1))`, poolAddress).Scan(&startTs, &endTs)
		if err != nil {
			return 0, 0, "", errors.Wrap(err, "could not get current epoch timestamps")
		}
		return startTs, endTs, "3*60*60", nil
	case "last":
		err := s.db.Connection().QueryRow(ctx, `select start_ts,end_ts  from smart_alpha.get_epoch_ts($1,((select p.epoch_id from smart_alpha.pool_epoch_info p where p.pool_address = $1 order by p.epoch_id desc limit 1)-1))`, poolAddress).Scan(&startTs, &endTs)
		if err != nil {
			return 0, 0, "", errors.Wrap(err, "could not get last epoch timestamps")
		}
		return startTs, endTs, "3*60*60", nil
	default:
		var err error
		duration, err = time.ParseDuration(window)
		startTs = time.Now().Add(-duration).Unix()
		endTs = time.Now().Unix()
		if err != nil {
			return 0, 0, "", errors.Wrap(err, "invalid window")
		}

		return startTs, endTs, "30*60", nil
	}
}
