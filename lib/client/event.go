package client

// Taken from Mattermost Server
type WebsocketBroadcast struct {
	OmitUsers map[string]bool `json:"omit_users"` // broadcast is omitted for users listed here
	UserId    string          `json:"user_id"`    // broadcast only occurs for this user
	ChannelId string          `json:"channel_id"` // broadcast only occurs for users in this channel
	TeamId    string          `json:"team_id"`    // broadcast only occurs for users in this team
}

type WebSocketEvent struct {
	Event     string                 `json:"event"`
	Data      map[string]interface{} `json:"data"`
	Broadcast *WebsocketBroadcast    `json:"broadcast"`
	Sequence  int64                  `json:"seq"`
}
