package utils

import (
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestBuildQuery1(t *testing.T) {
	filters := new(Filters)
	filters.Add("user_address", "0xdeadbeef")
	filters.Add("protocol_id", []string{"compound/v2", "aave/v2"})

	query, params := BuildQueryWithFilter(`
		select *
		from smart_yield_transaction_history
		where %s
		order by included_in_block desc, tx_index desc, log_index desc
		%s %s;
	`,
		filters,
		nil,
		nil,
	)

	assert.True(t, strings.Contains(query, "protocol_id = ANY($2)"))

	_, ok := params[1].(*pq.StringArray)
	assert.True(t, ok)
}

func TestBuildQuery(t *testing.T) {
	var limit int64 = 10
	var offset int64 = 0

	filters := new(Filters)
	filters.Add("user_address", "0xdeadbeef")
	filters.Add("protocol_id", "compound/v2")

	query, params := BuildQueryWithFilter(`
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
		where %s
		order by included_in_block desc, tx_index desc, log_index desc
		%s %s;
	`,
		filters,
		&limit,
		&offset,
	)

	assert.True(t, strings.Contains(query, "user_address = $1"))
	assert.True(t, strings.Contains(query, "protocol_id = $2"))
	assert.True(t, strings.Contains(query, "offset $3"))
	assert.True(t, strings.Contains(query, "limit $4"))
	assert.Len(t, params, 4)

	query, params = BuildQueryWithFilter(`
		select count(*)
		from smart_yield_transaction_history
		where %s
		order by included_in_block desc, tx_index desc, log_index desc
		%s %s;
	`,
		new(Filters).Add("user_address", "0xdeadbeef"),
		nil,
		nil,
	)

	assert.True(t, strings.Contains(query, "user_address = $1"))
	assert.Len(t, params, 1)
}
