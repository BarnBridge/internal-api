package governance

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"

	"github.com/barnbridge/internal-api/governance/types"
	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
)

func (g *Governance) HandleVoters(ctx *gin.Context) {
	builder := query.New()

	q, params := builder.WithPaginationFromCtx(ctx).Run(`
	select user_address, bond_staked, locked_until, delegated_power, votes, proposals, voting_power, has_active_delegation
	from governance.voters
	where bond_staked + voting_power > 0
	order by voting_power desc
	$offset$ $limit$
	`)

	rows, err := g.db.Connection().Query(ctx, q, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}
	defer rows.Close()

	var voters []types.Voter
	for rows.Next() {
		var v types.Voter

		err := rows.Scan(&v.Address, &v.BondStaked, &v.LockedUntil, &v.DelegatedPower, &v.Votes, &v.Proposals, &v.VotingPower, &v.HasActiveDelegation)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		voters = append(voters, v)
	}

	q, params = builder.Run(`select count(*) from governance.voters where bond_staked + voting_power > 0`)

	var count int
	err = g.db.Connection().QueryRow(ctx, q, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, g.db, voters, response.Meta().Set("count", count))
}
