package governance

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/governance/types"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/utils"
)

func (g *Governance) VotesHandler(ctx *gin.Context) {
	proposalIDString := ctx.Param("proposalID")
	proposalID, err := strconv.ParseInt(proposalIDString, 10, 64)
	if err != nil {
		response.Error(ctx, errors.New("invalid proposalID"))
		return
	}

	limit, err := utils.GetQueryLimit(ctx)
	if err != nil {
		response.BadRequest(ctx, err)
		return
	}

	page, err := utils.GetQueryPage(ctx)
	if err != nil {
		response.BadRequest(ctx, err)
		return
	}

	offset := (page - 1) * limit

	filters := utils.NewFilters()

	support := strings.ToLower(ctx.DefaultQuery("support", ""))
	if support != "" {
		if support != "true" && support != "false" {
			response.BadRequest(ctx, errors.New("wrong value for support parameter"))
			return
		}
		filters.Add("support", support)
	}

	query, params := utils.BuildQueryWithFilter(`
	select user_id, support, block_timestamp, power from governance.proposal_votes($param_overwrite$)
	$filters$
	order by power desc
	$offset$ $limit$
	`, filters, &limit, &offset)

	params = append(params, proposalID)
	query = strings.Replace(query, "$param_overwrite$", fmt.Sprintf("$%d", len(params)), 1)

	rows, err := g.db.Connection().Query(ctx, query, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var votes []types.Vote

	for rows.Next() {
		var v types.Vote

		err := rows.Scan(&v.User, &v.Support, &v.BlockTimestamp, &v.Power)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		votes = append(votes, v)
	}

	query, params = utils.BuildQueryWithFilter(`
	select count(*) from governance.proposal_votes($param_overwrite$)
	$filters$
	`, filters, nil, nil)

	params = append(params, proposalID)
	query = strings.Replace(query, "$param_overwrite$", fmt.Sprintf("$%d", len(params)), 1)

	var count int
	err = g.db.Connection().QueryRow(ctx, query, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, g.db, votes, map[string]interface{}{"count": count})
}
