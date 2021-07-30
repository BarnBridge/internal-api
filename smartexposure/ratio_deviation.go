package smartexposure

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartexposure/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartExposure) ratioDeviationChart(ctx *gin.Context) {
	tranche := ctx.Param("eTokenAddress")
	trancheAddress, err := utils.ValidateAccount(tranche)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	err, exists := s.checkTrancheExists(ctx, trancheAddress)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	if !exists {
		response.NotFound(ctx)
		return
	}

	window := strings.ToLower(ctx.DefaultQuery("window", "30d"))
	_, dateTrunc, err := validateWindow(window)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	dateTrunc = strings.Replace(dateTrunc, "'", "", -1)
	startTs, err := getStartDate(window)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	rows, err := s.db.Connection().Query(ctx, `select point, abs(deviation) as deviation from smart_exposure.get_ratio_deviation($1,$2,$3)`, trancheAddress, startTs, dateTrunc)
	if err != nil && err != sql.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var points []types.DeviationPoint
	for rows.Next() {
		var p types.DeviationPoint
		err := rows.Scan(&p.Point, &p.Deviation)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		points = append(points, p)
	}

	response.OK(ctx, points)
}
