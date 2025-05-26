package msgs

import (
	"github.com/reiver/go-json"
)

// Leave represents the JSON object with `type` = `Leave` ( [TypeLeave] ) that is sent and received over the LogJam web-socket.
//
// Its is intended to be similar to or the same as ActivityStream / ActivityPub "Leave" activity:
// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-leave
type Leave struct {
	Type  json.Const[string] `json:"type" json.value:"Leave"`
	Actor string             `json:"actor,omitempty"`
}
