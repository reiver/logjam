package msgs

// MessageContract represents the JSON object that is sent and received over the LogJam web-socket.
type MessageContract struct {
	Type   string
	Data   string
	Target string
	Name   string
}
