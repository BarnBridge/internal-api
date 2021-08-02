package types

import (
	"time"
)

type APYTrendPoint struct {
	Point            time.Time `json:"point"`
	SeniorAPY        float64   `json:"seniorApy"`
	JuniorAPY        float64   `json:"juniorApy"`
	OriginatorNetAPY float64   `json:"originatorNetApy"`
}
