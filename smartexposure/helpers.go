package smartexposure

import (
	"context"
	"errors"
	"time"

	"github.com/barnbridge/internal-api/smartexposure/types"
)

func (s *SmartExposure) checkPoolExists(ctx context.Context, addr string) (error, bool) {
	var exists bool

	err := s.db.Connection().QueryRow(ctx, `select (select count(*) from smart_exposure.pools where pool_address = $1) > 0`, addr).Scan(&exists)

	return err, exists
}

func (s *SmartExposure) checkTrancheExists(ctx context.Context, addr string) (error, bool) {
	var exists bool

	err := s.db.Connection().QueryRow(ctx, `select (select count(*) from smart_exposure.tranches where etoken_address = $1) > 0`, addr).Scan(&exists)

	return err, exists
}

func (s *SmartExposure) getPoolByETokenAddress(ctx context.Context, addr string) (*types.Pool, error) {
	var p types.Pool

	err := s.db.Connection().QueryRow(ctx, `
			select pool_address,
				   pool_name,
				   token_a_address,
				   token_a_symbol,
				   token_a_decimals,
				   token_b_address,
				   token_b_symbol,
				   token_b_decimals,
				   start_at_block
			from smart_exposure.pools
			where pool_address = (select pool_address from smart_exposure.tranches where etoken_address = $1)`, addr).Scan(&p.EPoolAddress, &p.ProtocolId, &p.TokenA.TokenAddress, &p.TokenA.TokenSymbol,
		&p.TokenA.TokenDecimals, &p.TokenB.TokenAddress, &p.TokenB.TokenSymbol, &p.TokenB.TokenDecimals, &p.StartAtBlock)

	return &p, err
}

func validateWindow(window string) (string, string, error) {
	if window == "24h" {
		return "'24 hours'", "'minute'", nil
	}

	if window == "1w" {
		return "'7 days'", "'hour'", nil
	}

	if window == "30d" {
		return "'30 days'", "'day'", nil
	}

	return "", "", errors.New("invalid window")
}

func getStartDate(window string) (int64, error) {
	var duration time.Duration

	switch window {
	case "1w":
		duration = 7 * 24 * time.Hour
	case "30d":
		duration = 30 * 24 * time.Hour
	default:
		var err error
		duration, err = time.ParseDuration(window)
		if err != nil {
			return 0, err
		}
	}

	return time.Now().Add(-duration).Unix(), nil
}
