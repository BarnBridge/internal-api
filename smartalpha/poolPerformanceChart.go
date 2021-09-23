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

	poolAddress, err := utils.ValidateAccount(poolAddress)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	rows, err := s.db.Connection().Query(ctx, `
				select p.epoch1_start,
					   e.block_timestamp
				from smart_alpha.pools p
					left join smart_alpha.epoch_end_events e on e.pool_address = p.pool_address
				where p.pool_address = $1
				order by e.epoch_id desc limit 2`, poolAddress)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	if err == pgx.ErrNoRows {
		response.NotFound(ctx)
		return
	}

	var epochTs []int64
	var epoch1Start *int64
	for rows.Next() {
		var ts *int64
		err = rows.Scan(&epoch1Start, &ts)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		if ts != nil {
			epochTs = append(epochTs, *ts)
		}
	}

	if epoch1Start == nil {
		response.NotFound(ctx)
		return
	} else if len(epochTs) == 0 {
		response.OK(ctx, []types.PerformancePoint{})
		return
	}

	window := strings.ToLower(ctx.DefaultQuery("window", "30d"))
	if window == "last" && len(epochTs) == 1 {
		response.OK(ctx, []types.PerformancePoint{})
		return
	}

	startTs, endTs, err := s.poolPerformanceWindow(window, epochTs)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	if startTs < *epoch1Start {
		startTs = *epoch1Start
	}

	if endTs < startTs {
		endTs = startTs
	}

	pointDistance := (endTs - startTs) / ChartNrOfPoints

	if pointDistance < 1 {
		pointDistance = 1
	}

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

	rows, err = s.db.Connection().Query(ctx, query, startTs, endTs, pointDistance, poolAddress)
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

func (s *SmartAlpha) poolPerformanceWindow(window string, epochTs []int64) (int64, int64, error) {
	var startTs, endTs int64
	var duration time.Duration

	switch window {
	case "1w":
		duration = 7 * 24 * time.Hour
		startTs = time.Now().Add(-duration).Unix()
		endTs = time.Now().Unix()

	case "30d":
		duration = 30 * 24 * time.Hour
		startTs = time.Now().Add(-duration).Unix()
		endTs = time.Now().Unix()

	case "current":
		startTs = epochTs[0]
		endTs = time.Now().Unix()

	case "last":
		startTs = epochTs[1]
		endTs = epochTs[0] - 1

	default:
		var err error
		duration, err = time.ParseDuration(window)
		startTs = time.Now().Add(-duration).Unix()
		endTs = time.Now().Unix()
		if err != nil {
			return 0, 0, errors.Wrap(err, "invalid window")
		}

	}

	return startTs, endTs, nil
}
