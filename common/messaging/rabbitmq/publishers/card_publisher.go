package publishers

import (
	"fmt"

	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

const (
	DeleteCard MessageType = iota
	// Add other message types here...
)

type CardPublisher struct {
	Channel *amqp.Channel
}

func (p *CardPublisher) Publish(messageType MessageType, message []byte) error {
	switch messageType {
	case DeleteCard:
		var msg pb_card.DeleteCardRequest
		err := proto.Unmarshal(message, &msg)
		if err != nil {
			return err
		}

		err = p.publishDeleteCardMessage(&msg)
		if err != nil {
			return err
		}
	// Add other cases for other message types here...
	default:
		return fmt.Errorf("invalid message type: %v", messageType)
	}

	return nil
}

func NewCardPublisher(ch *amqp.Channel) *CardPublisher {
	return &CardPublisher{Channel: ch}
}

func (p *CardPublisher) publishDeleteCardMessage(req *pb_card.DeleteCardRequest) error {
	ch := p.Channel
	defer ch.Close()

	message, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	err = ch.Publish("halten",
		"card.delete",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/protobuf",
			Body:        message,
		})
	return err
}
