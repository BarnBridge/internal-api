package smartalpha

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartalpha/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartAlpha) transactions(ctx *gin.Context) {
	builder := query.New()
	poolAddress := ctx.DefaultQuery("poolAddress", "")
	if poolAddress != "" {
		poolAddress, err := utils.ValidateAccount(poolAddress)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		builder.Filters.Add("t.pool_address", poolAddress)
	}

	accountAddress := strings.ToLower(ctx.DefaultQuery("userAddress", ""))
	if accountAddress != "" {
		accountAddress, err := utils.ValidateAccount(accountAddress)
		if err != nil {
			response.BadRequest(ctx, errors.New("invalid accountAddress"))
			return
		}
		builder.Filters.Add("user_address", accountAddress)
	}

	transactionType := strings.ToUpper(ctx.DefaultQuery("transactionType", "ALL"))
	if transactionType != "ALL" {
		if !checkTxType(transactionType) {
			response.BadRequest(ctx, errors.New("invalid transaction type"))
			return
		}
		builder.Filters.Add("transaction_type", transactionType)
	}

	q, params := builder.WithPaginationFromCtx(ctx).Run(`
		select t.pool_address,
			   t.user_address,
			   t.tranche,
			   t.transaction_type,
			   t.amount,
			   (select token_price_at_ts((select pool_token_address
										  from smart_alpha.pools p
										  where p.pool_address = t.pool_address
										  limit 1), (select oracle_asset_symbol
													 from smart_alpha.pools p
													 where p.pool_address = t.pool_address
													 limit 1), t.block_timestamp)),
			   coalesce((select junior_token_price_start
				from smart_alpha.pool_epoch_info pi
				where pi.pool_address = t.pool_address and pi.block_timestamp <= t.block_timestamp
				order by epoch_id desc limit 1), '1e18'),
			   coalesce((select senior_token_price_start
				from smart_alpha.pool_epoch_info pi
				where pi.pool_address = t.pool_address and pi.block_timestamp <= t.block_timestamp
				order by epoch_id desc limit 1), '1e18'),
			   (select token_usd_price_at_ts((select pool_token_address
											  from smart_alpha.pools p
											  where p.pool_address = t.pool_address
											  limit 1), t.block_timestamp)),
			   t.block_timestamp,
			   t.tx_hash,
			   (select pool_token_decimals from smart_alpha.pools p where p.pool_address = t.pool_address limit 1) as decimals,
			   p.oracle_asset_symbol,
			   p.pool_token_symbol,
			   p.junior_token_symbol,
			   p.senior_token_symbol,
	           p.pool_token_address,
	           p.junior_token_address,
	           p.senior_token_address
		from smart_alpha.transaction_history t
				 inner join smart_alpha.pools p on p.pool_address = t.pool_address
		$filters$
		order by block_timestamp desc, tx_index desc, log_index desc
		$offset$ $limit$
	`)
	rows, err := s.db.Connection().Query(ctx, q, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}
	defer rows.Close()

	var history = make([]types.Transaction, 0)
	for rows.Next() {
		var h types.Transaction
		var decimals int32
		var juniorTokenSymbol, seniorTokenSymbol string
		var juniorTokenAddress, seniorTokenAddress string
		var poolTokenPrice, poolTokenPriceUSD, juniorTokenPrice, seniorTokenPrice decimal.Decimal
		err := rows.Scan(&h.PoolAddress, &h.UserAddress, &h.Tranche, &h.TransactionType, &h.Amount, &poolTokenPrice, &juniorTokenPrice, &seniorTokenPrice, &poolTokenPriceUSD, &h.BlockTimestamp, &h.TransactionHash, &decimals,
			&h.OracleAssetSymbol, &h.PoolTokenSymbol, &juniorTokenSymbol, &seniorTokenSymbol, &h.PoolTokenAddress, &juniorTokenAddress, &seniorTokenAddress)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		h.Amount = h.Amount.Shift(-decimals)

		juniorTokenPrice = juniorTokenPrice.Shift(-18)
		seniorTokenPrice = seniorTokenPrice.Shift(-18)

		amountInAsset := getAmountInAsset(h.TransactionType, h.Amount, juniorTokenPrice, seniorTokenPrice)

		h.AmountInQuoteAsset = amountInAsset.Mul(poolTokenPrice)
		h.AmountInUSD = amountInAsset.Mul(poolTokenPriceUSD)

		h.TokenSymbol = getTxTokenSymbol(h.TransactionType, h.PoolTokenSymbol, juniorTokenSymbol, seniorTokenSymbol)
		switch h.TokenSymbol {
		case h.PoolTokenSymbol:
			h.TokenAddress = h.PoolTokenAddress
		case juniorTokenSymbol:
			h.TokenAddress = juniorTokenAddress
		case seniorTokenSymbol:
			h.TokenAddress = seniorTokenAddress
		}

		history = append(history, h)
	}

	q, params = builder.Run(`select count(*) from smart_alpha.transaction_history t $filters$`)
	var count int64

	err = s.db.Connection().QueryRow(ctx, q, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, s.db, history, response.Meta().Set("count", count))
}
