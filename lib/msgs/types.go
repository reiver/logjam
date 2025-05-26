package msgs

// lib/msgs.Message.Type
const (
	TypeAddAudience                  = "add_audience"
	TypeAltBroadcast                 = "alt-broadcast"
	TypeBroadcasting                 = "broadcasting"
	TypeError                        = "error"
	TypeEventBroadcasterDisconnected = "event-broadcaster-disconnected"
	TypeEventParentDC                = "event-parent-dc"
	TypeGoldGorillaJoined            = "goldgorilla-joined"
	TypeMetaDataGet                  = "metadata-get"
	TypeMetaDataSet                  = "metadata-set"
	TypeMuted                        = "muted"
	TypeNewIceCandidate              = "new-ice-candidate"
	TypeNewMessage                   = "new-message"
	TypePing                         = "ping"
	TypePong                         = "pong"
	TypeReconnect                    = "reconnect"
	TypeReconnectChildren            = "reconnect-children"
	TypeRole                         = "role"
	TypeSendMessage                  = "send-message"
	TypeStart                        = "start"
	TypeStream                       = "stream"
	TypeTree                         = "tree"
	TypeTurnStatus                   = "turn_status"
	TypeUpdateStreamID               = "updateStreamId"
	TypeUserByStream                 = "user-by-stream"
	TypeVideoAnswer                  = "video-answer"
	TypeVideoOffer                   = "video-offer"
)

// lib/msgs.Message.Type
//
// Event
const (
	TypeLeave                        = "Leave"
)

// lib/msgs.Message.Type
//
// Query
const (
	TypeGetLatestUserList            = "get-latest-user-list"
)

// lib/msgs.Message.Type
//
// Response
const (
	TypeUserEvent                    = "user-event"
)
