package smartexposure

import (
	"github.com/gin-gonic/gin"
)

func (s *SmartExposure) handleUserPortfolioValue(ctx *gin.Context) {
	/*user, err := getQueryAddress(c, "userAddress")
		if err != nil {
			BadRequest(c, err)
			return
		}

		window := strings.ToLower(c.DefaultQuery("window", "30d"))
		totalPoints := getTotalPoints(window)
		window, _, err = validateWindow(window)
		if err != nil {
			Error(c, err)
			return
		}

		poolAddress := strings.ToLower(c.DefaultQuery("poolAddress", "all"))
		if poolAddress != "all" {
			poolAddress, err = utils.ValidateAccount(poolAddress)
			if err != nil {
				Error(c, err)
				return
			}
			if state.SEPoolByAddress(poolAddress) == nil {
				BadRequest(c, errors.New("invalid pool address"))
				return
			}
		}

		var query string
		var params []interface{}
		params = append(params, user)
		if poolAddress != "all" {
			query = fmt.Sprintf(`select to_timestamp(ts),
	       coalesce (se_user_portfolio_value_by_pool($1,ts,$2),0)
	from generate_series(( select extract(epoch from now() - interval %s)::bigint ),
	                     ( select extract(epoch from now()) )::bigint, %s) as ts order by ts;`, window, totalPoints)
			params = append(params, poolAddress)
		} else {
			query = fmt.Sprintf(`select to_timestamp(ts),
	       coalesce (se_user_portfolio_value($1,ts),0)
	from generate_series(( select extract(epoch from now() - interval %s)::bigint ),
	                     ( select extract(epoch from now()) )::bigint, %s) as ts order by ts;`, window, totalPoints)
		}

		//24h 30*60
		//7d 3*60*60
		//30d 12*60*60
		var points []SEPortfolioValuePoint

		rows, err := a.db.Query(query, params...)

		if err != nil {
			Error(c, err)
			return
		}

		for rows.Next() {
			var p SEPortfolioValuePoint
			err := rows.Scan(&p.Point, &p.PortfolioValueSE)
			if err != nil {
				Error(c, err)
				return
			}
			points = append(points, p)
		}

		OK(c, points)*/

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
