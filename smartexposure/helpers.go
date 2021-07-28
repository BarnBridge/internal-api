package smartexposure

import (
	"context"
	"errors"
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
