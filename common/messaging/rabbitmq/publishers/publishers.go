package publishers

type MessageType int

type Publishers struct {
	ListPublisher Publisher
	CardPublisher Publisher
	// Add other publishers here...
}

type Publisher interface {
	Publish(messageType MessageType, message []byte) error
}
