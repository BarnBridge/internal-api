package governance

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/governance/types"
	"github.com/barnbridge/internal-api/response"
)

func (g *Governance) HandleProposalDetails(ctx *gin.Context) {
	proposalID, err := getProposalId(ctx)
	if err != nil {
		response.Error(ctx, errors.New("invalid proposalID"))
		return
	}

	var p types.ProposalFull
	err = g.db.Connection().QueryRow(ctx, `
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
		       coalesce(( select sum(power) from governance.proposal_votes(proposal_id) where support = true ), 0) as for_votes,
			   coalesce(( select sum(power) from governance.proposal_votes(proposal_id) where support = false ), 0) as against_votes,
		       coalesce(( select governance.bond_staked_at_ts(create_time+warm_up_duration) ), 0) as bond_staked,
			   ( select * from governance.proposal_state(proposal_id) ) as proposal_state
		from governance.proposals
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

	response.OKWithBlock(ctx, g.db, p)
}
