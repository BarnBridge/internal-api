package types

type Notification struct {
	Id               int64            `json:"id"`
	Target           string           `json:"target"`
	NotificationType string           `json:"notificationType"`
	StartsOn         int64            `json:"startsOn"`
	ExpiresOn        int64            `json:"expiresOn"`
	Message          string           `json:"message"`
	Metadata         JSONObject `json:"metadata"`
}
