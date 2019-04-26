package core

const (
	// KILLPILL - System wide messageType to initiate a shutdown/close of all registered actors and eventually of the actor system
	KILLPILL = "KILLPILL"
)

// DeliveryMode - Different delivery modes of the messages supported by the actor system
type DeliveryMode int

const (
	// Unicast - Single target actor message delivery mode
	Unicast DeliveryMode = 1 + iota
	// Broadcast - Multi-target actor message delivery mode
	Broadcast
)

var deliveryTypes = [...]string{
	"Unicast",
	"Broadcast",
}

// String - Returns the string representing of th DeliveryMode of the message
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

// ActorReference - Simple reference structure to uniquely identify an actor registered in the system
type ActorReference struct {
	ActorType string `json:"ActorType"`
}
