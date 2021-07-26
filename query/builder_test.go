package query

import (
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBuildQuery1(t *testing.T) {
	qb := New()
	qb.Filters.Add("user_address", "0xdeadbeef")
	qb.Filters.Add("protocol_id", []string{"compound/v2", "aave/v2"})

	query, params := qb.WithPagination(0, 1).Run(`
		select *
		from smart_yield_transaction_history
		$filters$
		order by included_in_block desc, tx_index desc, log_index desc
		$offset$ $limit$;
	`)

	assert.True(t, strings.Contains(query, "protocol_id = ANY($2)"))

	_, ok := params[1].([]string)
	assert.True(t, ok)
}

func TestBuildQuery(t *testing.T) {
	qb := New()

	qb.Filters.Add("user_address", "0xdeadbeef")
	qb.Filters.Add("protocol_id", "compound/v2")

	query, params := qb.WithPagination(0, 2).Run(`
		select protocol_id,
			   sy_address,
			   underlying_token_address,
			   amount,
			   tranche,
			   transaction_type,
			   tx_hash,
			   block_timestamp,
			   included_in_block
		from smart_yield_transaction_history
		$filters$
		order by included_in_block desc, tx_index desc, log_index desc
		$offset$ $limit$;
	`)

	assert.True(t, strings.Contains(query, "user_address = $1"))
	assert.True(t, strings.Contains(query, "protocol_id = $2"))
	assert.True(t, strings.Contains(query, "offset $3"))
	assert.True(t, strings.Contains(query, "limit $4"))
	assert.Len(t, params, 4)
	assert.Equal(t, int64(0), params[2])
	assert.Equal(t, int64(2), params[3])


	qb.Filters = new(Filters).Add("user_address", "0xdeadbeef")

	query, params = qb.Run(`
		select count(*)
		from smart_yield_transaction_history
		$filters$
	`)

	assert.True(t, strings.Contains(query, "user_address = $1"))
	assert.Len(t, params, 1)
}

func TestBuildQueryWithPaginationFromCtx(t *testing.T) {
	qb := New()

	qb.Filters.Add("user_address", "0xdeadbeef")
	qb.Filters.Add("protocol_id", "compound/v2")

	query, params := qb.WithPaginationFromCtx(&gin.Context{}).Run(`
		select protocol_id,
			   sy_address,
			   underlying_token_address,
			   amount,
			   tranche,
			   transaction_type,
			   tx_hash,
			   block_timestamp,
			   included_in_block
		from smart_yield_transaction_history
		$filters$
		order by included_in_block desc, tx_index desc, log_index desc
		$offset$ $limit$;
	`)

	assert.True(t, strings.Contains(query, "user_address = $1"))
	assert.True(t, strings.Contains(query, "protocol_id = $2"))
	assert.True(t, strings.Contains(query, "offset $3"))
	assert.True(t, strings.Contains(query, "limit $4"))
	assert.Len(t, params, 4)
	// defaults
	assert.Equal(t, int64(0), params[2])
	assert.Equal(t, int64(10), params[3])
}
