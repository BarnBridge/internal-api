package smartalpha

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartalpha/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartAlpha) UserPortfolioValue(ctx *gin.Context) {
	userAddress := ctx.Param("address")
	userAddress, err := utils.ValidateAccount(userAddress)
	if err != nil {
		response.Error(ctx, errors.Wrap(err, "invalid user address"))
	}

	window := strings.ToLower(ctx.DefaultQuery("window", "30d"))
	totalPoints := getTotalPoints(window)
	window, _, err = validateWindow(window)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	query := fmt.Sprintf(`
			select to_timestamp(ts),
				   coalesce((select smart_alpha.junior_portfolio_value_at_ts($1, ts)), 0) as junior_value,
				   coalesce((select smart_alpha.senior_portfolio_value_at_ts($1, ts)), 0) as senior_value,
	               coalesce((select smart_alpha.entry_queue_portfolio_value_at_ts($1, ts)), 0) as entry_queue_value,
	               coalesce((select smart_alpha.exit_queue_portfolio_value_at_ts($1, ts)), 0) as exit_queue_value
			from generate_series((select extract(epoch from now() - interval % s)::bigint),
								 (select extract(epoch from now()))::bigint, %s) as ts
			order by ts;`, window, totalPoints)

	rows, err := s.db.Connection().Query(ctx, query, userAddress)
	if err != nil && err != sql.ErrNoRows {
		response.Error(ctx, err)
		return
	}
	var points []types.UserPortfolioPoint

	for rows.Next() {
		var p types.UserPortfolioPoint
		err := rows.Scan(&p.Point, &p.JuniorValue, &p.SeniorValue, &p.EntryQueueValue, &p.ExitQueueValue)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		points = append(points, p)
	}

	response.OK(ctx, points)
}
