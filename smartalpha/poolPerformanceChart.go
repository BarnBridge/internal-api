package smartalpha

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartalpha/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartAlpha) poolPerformanceChart(ctx *gin.Context) {
	poolAddress := ctx.Param("poolAddress")
	if poolAddress == "" {
		response.BadRequest(ctx, errors.New("invalid pool address"))
		return
	}

	var epoch1Start, epochDuration, currentEpoch int64

	poolAddress, err := utils.ValidateAccount(poolAddress)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	err = s.db.Connection().QueryRow(ctx, `
				select p.epoch1_start,
					   p.epoch_duration
				from smart_alpha.pools p
				where p.pool_address = $1`, poolAddress).Scan(&epoch1Start, &epochDuration)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	if err == pgx.ErrNoRows {
		response.NotFound(ctx)
		return
	}

	currentEpoch = getCurrentEpoch(epoch1Start, epochDuration)
	window := strings.ToLower(ctx.DefaultQuery("window", "30d"))
	if window == "last" && currentEpoch == 0 {
		response.OK(ctx, []types.PerformancePoint{})
		return
	}

	startTs, endTs, err := s.poolPerformanceWindow(window, epoch1Start, epochDuration, currentEpoch)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	pointDistance := (endTs - startTs) / ChartNrOfPoints

	query := `	select to_timestamp(ts)               as point,
					   coalesce(senior_without_sa, 0) as senior_without_sa,
					   coalesce(senior_with_sa, 0)    as senior_with_sa,
					   coalesce(junior_without_sa, 0) as junior_without_sa,
					   coalesce(junior_with_sa, 0)    as junior_with_sa,
					   coalesce(token_price_at_ts(( select pool_token_address from smart_alpha.pools where pool_address = $4 ),
												  ( select oracle_asset_symbol from smart_alpha.pools where pool_address = $4 ), ts),
								0)                    as pool_token_price
				from generate_series($1::bigint, $2::bigint, $3::bigint) as ts
						 inner join smart_alpha.performance_at_ts($4, ts) on true
				order by ts;`

	rows, err := s.db.Connection().Query(ctx, query, startTs, endTs, pointDistance, poolAddress)
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

func (s *SmartAlpha) poolPerformanceWindow(window string, epoch1Start int64, epochDuration int64, currentEpoch int64) (int64, int64, error) {
	var startTs, endTs int64
	var duration time.Duration

	switch window {
	case "1w":
		duration = 7 * 24 * time.Hour
		startTs = time.Now().Add(-duration).Unix()
		endTs = time.Now().Unix()

		return startTs, endTs, nil
	case "30d":
		duration = 30 * 24 * time.Hour
		startTs = time.Now().Add(-duration).Unix()
		endTs = time.Now().Unix()

		return startTs, endTs, nil
	case "current":
		startTs = epoch1Start + (currentEpoch-1)*epochDuration
		endTs = epoch1Start + currentEpoch*epochDuration - 1

		return startTs, endTs, nil
	case "last":
		startTs = epoch1Start + (currentEpoch-2)*epochDuration
		endTs = epoch1Start + (currentEpoch-1)*epochDuration - 1

		return startTs, endTs, nil
	default:
		var err error
		duration, err = time.ParseDuration(window)
		startTs = time.Now().Add(-duration).Unix()
		endTs = time.Now().Unix()
		if err != nil {
			return 0, 0, errors.Wrap(err, "invalid window")
		}

		return startTs, endTs, nil
	}
}
