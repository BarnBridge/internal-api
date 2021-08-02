package smartexposure

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartexposure/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartExposure) userPortfolioValueChart(ctx *gin.Context) {
	user := ctx.Param("userAddress")
	if user != "" {
		var err error
		user, err = utils.ValidateAccount(user)
		if err != nil {
			response.BadRequest(ctx, errors.New("invalid accountAddress"))
			return
		}
	}

	window := strings.ToLower(ctx.DefaultQuery("window", "30d"))
	totalPoints := getTotalPoints(window)
	window, _, err := validateWindow(window)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	poolAddress := ctx.DefaultQuery("poolAddress", "all")
	if poolAddress != "all" {
		poolAddress, err = utils.ValidateAccount(poolAddress)
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

	var query string
	var params []interface{}
	params = append(params, user)
	if poolAddress != "all" {
		query = fmt.Sprintf(`select to_timestamp(ts),
       coalesce (smart_exposure.user_portfolio_value_by_pool($1,ts,$2),0)
from generate_series(( select extract(epoch from now() - interval %s)::bigint ),
                     ( select extract(epoch from now()) )::bigint, %s) as ts order by ts;`, window, totalPoints)
		params = append(params, poolAddress)
	} else {
		query = fmt.Sprintf(`select to_timestamp(ts),
       coalesce (smart_exposure.user_portfolio_value($1,ts),0)
from generate_series(( select extract(epoch from now() - interval %s)::bigint ),
                     ( select extract(epoch from now()) )::bigint, %s) as ts order by ts;`, window, totalPoints)
	}

	//24h 30*60
	//7d 3*60*60
	//30d 12*60*60
	var points []types.UserPortfolioValuePoint

	rows, err := s.db.Connection().Query(ctx, query, params...)

	if err != nil {
		response.Error(ctx, err)
		return
	}

	for rows.Next() {
		var p types.UserPortfolioValuePoint
		err := rows.Scan(&p.Point, &p.PortfolioValue)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		points = append(points, p)
	}

	response.OK(ctx, points)

}

func getTotalPoints(window string) string {
	if window == "24h" {
		return "30 * 60"
	}
	if window == "1w" {
		return "3*60*60"
	}
	if window == "30d" {
		return "12*60*60"
	}
	return ""
}
