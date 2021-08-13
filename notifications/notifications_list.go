package notifications

import (
	"strconv"
	"strings"

	"github.com/barnbridge/internal-api/notifications/types"
	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (h *Notifications) NotificationsList(ctx *gin.Context) {
	builder := query.New()

	timestamp := ctx.DefaultQuery("timestamp", "-1")
	_, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		response.BadRequest(ctx, errors.Wrap(err, "invalid 'timestamp' parameter"))
		return
	}

	builder.Filters.Add("starts_on", timestamp, ">")
	builder.Filters.AddRaw( "starts_on < extract(epoch from now())::bigint")
	builder.Filters.AddRaw( "extract(epoch from now())::bigint < expires_on")


	target := strings.ToLower(ctx.DefaultQuery("target", ""))
	if target != "" {
		builder.Filters.Add("target", []string{"system", target})
	} else {
		builder.Filters.Add("target", "system")
	}


	query, params := builder.WithPaginationFromCtx(ctx).Run(`
		select "id", "target", "type", "starts_on", "expires_on", "message", "metadata"
		from public."notifications"
		$filters$
		order by "starts_on"
		$offset$ $limit$
		;
	`)

	rows, err := h.db.Connection().Query(ctx, query, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var notifications []types.Notification
	for rows.Next() {
		var n types.Notification
		err := rows.Scan(&n.Id, &n.Target, &n.NotificationType, &n.StartsOn, &n.ExpiresOn, &n.Message, &n.Metadata)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		notifications = append(notifications, n)
	}

	response.OK(ctx, notifications)
}
