package smartalpha

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartAlpha) UserQueuePositions(ctx *gin.Context) {
	userAddress := ctx.Param("address")
	userAddress, err := utils.ValidateAccount(userAddress)
	if err != nil {
		response.Error(ctx, errors.Wrap(err, "invalid user address"))
	}
}
