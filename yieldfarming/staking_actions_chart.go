package yieldfarming

import (
	"strings"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/utils"
	"github.com/barnbridge/internal-api/yieldfarming/types"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (h *YieldFarming) StakingActionsChart(ctx *gin.Context) {
	tokensAddress := strings.ToLower(ctx.DefaultQuery("tokenAddress", ""))
	if tokensAddress == "" {
		response.BadRequest(ctx, errors.New("tokenAddress required"))
		return
	}

	tokens := strings.Split(tokensAddress, ",")

	for i, token := range tokens {
		t, err := utils.ValidateAccount(token)
		if err != nil {
			response.BadRequest(ctx, err)
			return
		}
		tokens[i] = t
	}

	startTsString := strings.ToLower(ctx.DefaultQuery("start", "-1"))
	startTs, err := validateTs(startTsString)
	if err != nil {
		response.BadRequest(ctx, err)
		return
	}

	endTsString := strings.ToLower(ctx.DefaultQuery("end", "-1"))
	endTs, err := validateTs(endTsString)
	if err != nil {
		response.BadRequest(ctx, err)
		return
	}

	scale := strings.ToLower(ctx.DefaultQuery("scale", "week"))
	if scale != "week" && scale != "day" {
		response.BadRequest(ctx, errors.New("Wrong scale"))
		return
	}

	charts := make(map[string]types.Chart)

	query := `select * from yield_farming.yf_stats_by_token($1,$2,$3,$4) order by point`
	for _, token := range tokens {
		rows, err := h.db.Connection().Query(ctx, query, token, startTs.Unix(), endTs.Unix(), scale)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		chart, err := getChart(rows)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		charts[token] = *chart

	}

	response.OK(ctx, charts)
}
