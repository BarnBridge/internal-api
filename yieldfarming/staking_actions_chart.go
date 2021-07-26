package yieldfarming

import (
	"github.com/barnbridge/internal-api/response"
	"github.com/gin-gonic/gin"
)

func (h *YieldFarming) StakingActionsChart(ctx *gin.Context) {


	response.OKWithBlock(ctx, h.db, nil)
}

/*

func (a *API) handleStakinsActionsChart(c *gin.Context) {
	tokensAddress := strings.ToLower(c.DefaultQuery("tokenAddress", ""))

	if tokensAddress == "" {
		BadRequest(c, errors.New("tokenAddress required"))
		return
	}

	tokens := strings.Split(tokensAddress, ",")

	for i, token := range tokens {
		t, err := utils.ValidateAccount(token)
		if err != nil {
			BadRequest(c, err)
			return
		}
		tokens[i] = t
	}

	startTsString := strings.ToLower(c.DefaultQuery("start", "-1"))
	startTs, err := validateTs(startTsString)
	if err != nil {
		BadRequest(c, err)
		return
	}

	endTsString := strings.ToLower(c.DefaultQuery("end", "-1"))
	endTs, err := validateTs(endTsString)
	if err != nil {
		BadRequest(c, err)
		return
	}

	scale := strings.ToLower(c.DefaultQuery("scale", "week"))
	if scale != "week" && scale != "day" {
		BadRequest(c, errors.New("Wrong scale"))
		return
	}

	charts := make(map[string]types.Chart)

	for _, token := range tokens {
		rows, err := a.db.Query(`select * from yf_stats_by_token($1,$2,$3,$4) order by point`, token, startTs.Unix(), endTs.Unix(), scale)
		if err != nil {
			Error(c, err)
			return
		}

		chart, err := getChart(rows)
		if err != nil {
			Error(c, err)
			return
		}
		charts[token] = *chart

	}

	OK(c, charts)
	return
}

func getChart(rows *sql.Rows) (*types.Chart, error) {
	x := make(types.Chart)

	for rows.Next() {
		var t time.Time
		var a types.Aggregate
		err := rows.Scan(&t, &a.SumDeposits, &a.SumWithdrawals)
		if err != nil {
			return nil, err
		}

		x[t] = a
	}
	return &x, nil
}

func validateTs(ts string) (*time.Time, error) {
	timestamp, err := strconv.ParseInt(ts, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "invalid timestamp")
	}

	if timestamp == -1 {
		return nil, errors.New("timestamp required")
	}

	tm := time.Unix(timestamp, 0)

	return &tm, nil
}

 */
