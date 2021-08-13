package notifications

import (
	"github.com/barnbridge/internal-api/db"
	"github.com/gin-gonic/gin"
)

type Notifications struct {
	db *db.DB
}

func New(db *db.DB) *Notifications {
	return &Notifications{db: db}
}

func (h *Notifications) SetRoutes(engine *gin.Engine) {
	notifs := engine.Group("/api/notifications")
	notifs.GET("/list", h.NotificationsList)
}
