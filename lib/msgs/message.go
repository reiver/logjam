package msgs

// Message represents the JSON object that is sent and received over the LogJam web-socket.
type Message struct {
	Type   string
	Data   string
	Target string
	Name   string
}
