package governance

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"

	"github.com/barnbridge/internal-api/governance/types"
	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
)

func (g *Governance) HandleProposals(ctx *gin.Context) {
	builder := query.New()

	title := ctx.DefaultQuery("title", "")
	if title != "" {
		builder.Filters.Add("lower(title)", "%"+strings.ToLower(title)+"%", "like")
	}

	proposalState := strings.ToUpper(ctx.DefaultQuery("state", "all"))
	if proposalState != "ALL" {
		var states []string
		if proposalState == "ACTIVE" {
			states = []string{"WARMUP", "ACTIVE", "ACCEPTED", "QUEUED", "GRACE"}
		} else if proposalState == "FAILED" {
			states = []string{"CANCELED", "FAILED", "ABROGATED", "EXPIRED"}
		} else {
			states = []string{proposalState}
		}

		builder.Filters.Add("(select governance.proposal_state(proposal_id) )", states)
	}

	q, params := builder.WithPaginationFromCtx(ctx).Run(`
		select proposal_id,
			   proposer,
			   description,
			   title,
			   create_time,
			   warm_up_duration,
			   active_duration,
			   queue_duration,
			   grace_period_duration,
			   ( select governance.proposal_state(proposal_id) ) as proposal_state,
			   coalesce(( select sum(power) from governance.proposal_votes(proposal_id) where support = true ), 0) as for_votes,
			   coalesce(( select sum(power) from governance.proposal_votes(proposal_id) where support = false ), 0) as against_votes
		from governance.proposals
		$filters$
		order by proposal_id desc
		$offset$ $limit$
	`)

	rows, err := g.db.Connection().Query(ctx, q, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}
	defer rows.Close()

	var proposals []types.ProposalBase

	for rows.Next() {
		var p types.ProposalBase
		var (
			createTime          int64
			warmUpDuration      int64
			activeDuration      int64
			queueDuration       int64
			gracePeriodDuration int64
		)

		err := rows.Scan(&p.Id, &p.Proposer, &p.Description, &p.Title, &createTime, &warmUpDuration, &activeDuration, &queueDuration, &gracePeriodDuration, &p.State, &p.ForVotes, &p.AgainstVotes)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		p.CreateTime = createTime
		p.StateTimeLeft = getTimeLeft(p.State, createTime, warmUpDuration, activeDuration, queueDuration, gracePeriodDuration)

		proposals = append(proposals, p)
	}

	q, params = builder.Run(`
		select count(*) from governance.proposals
		$filters$
	`)

	var count int
	err = g.db.Connection().QueryRow(ctx, q, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, g.db, proposals, response.Meta().Set("count", count))
}
