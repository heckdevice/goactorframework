package core

const (
	KILLPILL = "KILLPILL"
)

type DeliveryMode int

const (
	Unicast DeliveryMode = 1 + iota
	Broadcast
)

var deliveryTypes = [...]string{
	"Unicast",
	"Broadcast",
}

func (dm DeliveryMode) String() string { return deliveryTypes[dm-1] }

// Message - Simple message payload
type Message struct {
	MessageType string
	Mode        DeliveryMode
	Payload     interface{}
	Sender      *ActorReference
	UnicastTo   *ActorReference
	BroadcastTo []*ActorReference
}

type ActorReference struct {
	ActorType string `json:"ActorType"`
}
