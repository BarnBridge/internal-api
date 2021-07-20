package governance

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"

	"github.com/barnbridge/internal-api/governance/types"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/utils"
)

func (g *Governance) ProposalDetailsHandler(ctx *gin.Context) {
	proposalID := ctx.Param("proposalID")

	var p types.ProposalFull
	err := g.db.Connection().QueryRow(ctx, `
		select proposal_id,
			   proposer,
			   description,
			   title,
			   create_time,
			   targets,
			   "values",
			   signatures,
			   calldatas,
			   block_timestamp,
			   warm_up_duration,
			   active_duration,
			   queue_duration,
			   grace_period_duration,
			   acceptance_threshold,
			   min_quorum,
		       coalesce(( select sum(power) from proposal_votes(proposal_id) where support = true ), 0) as for_votes,
			   coalesce(( select sum(power) from proposal_votes(proposal_id) where support = false ), 0) as against_votes,
		       coalesce(( select bond_staked_at_ts(to_timestamp(create_time+warm_up_duration)) ), 0) as bond_staked,
			   ( select * from proposal_state(proposal_id) ) as proposal_state
		from governance_proposals
		where proposal_id = $1
	`, proposalID).Scan(
		&p.Id,
		&p.Proposer,
		&p.Description,
		&p.Title,
		&p.CreateTime,
		&p.Targets,
		&p.Values,
		&p.Signatures,
		&p.Calldatas,
		&p.BlockTimestamp,
		&p.WarmUpDuration,
		&p.ActiveDuration,
		&p.QueueDuration,
		&p.GracePeriodDuration,
		&p.AcceptanceThreshold,
		&p.MinQuorum,
		&p.ForVotes,
		&p.AgainstVotes,
		&p.BondStaked,
		&p.State,
	)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	} else if err == pgx.ErrNoRows {
		response.NotFound(ctx)
		return
	}

	p.StateTimeLeft = getTimeLeft(p.State, p.CreateTime, p.WarmUpDuration, p.ActiveDuration, p.QueueDuration, p.GracePeriodDuration)

	history, err := g.history(ctx, p)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	p.History = history

	block, err := utils.GetHighestBlock(ctx, g.db)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OK(ctx, p, map[string]interface{}{
		"block": block,
	})
}
