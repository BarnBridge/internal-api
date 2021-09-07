package governance

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"

	"github.com/barnbridge/internal-api/config"
	"github.com/barnbridge/internal-api/governance/types"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/utils"
)

func (g *Governance) HandleOverview(ctx *gin.Context) {
	batch := &pgx.Batch{}

	batch.Queue(`select coalesce(avg(locked_until - block_timestamp),0)::bigint from governance.barn_locks;`)
	batch.Queue(`select coalesce(sum(governance.voting_power(user_address)),0) as total_voting_power from governance.barn_users;`)
	batch.Queue(`select count(*) from erc20_users_with_balance($1) where balance > 0;`, utils.NormalizeAddress(config.Store.Addresses.Bond))
	batch.Queue(`select count(*) from erc20_users_with_balance_excluded_transfers($1, $2) where balance > 0;`, utils.NormalizeAddress(config.Store.Addresses.Bond), utils.NormalizeAddresses(config.Store.Addresses.ExcludeTransfers))
	batch.Queue(`
    select coalesce(sum(total),0) 
    from ( select case when action_type = 'INCREASE' then sum(amount)
                       when action_type = 'DECREASE' then -sum(amount) end total
           from governance.barn_delegate_changes
           group by action_type ) x;
	`)
	batch.Queue(`
	select count(*)
	from ( select distinct user_id as address
           from governance.votes
           union
           select distinct user_id
           from governance.abrogation_votes ) x; 
   `)
	batch.Queue(`select count(*) from governance.voters where bond_staked + voting_power > 0`)

	res := g.db.Connection().SendBatch(ctx, batch)
	defer res.Close()

	var overview types.Overview
	err := res.QueryRow().Scan(&overview.AvgLockTimeSeconds)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	err = res.QueryRow().Scan(&overview.TotalVBond)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	err = res.QueryRow().Scan(&overview.Holders)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	err = res.QueryRow().Scan(&overview.HoldersStakingExcluded)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	err = res.QueryRow().Scan(&overview.TotalDelegatedPower)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	err = res.QueryRow().Scan(&overview.Voters)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	err = res.QueryRow().Scan(&overview.BarnUsers)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, g.db, overview)
}
