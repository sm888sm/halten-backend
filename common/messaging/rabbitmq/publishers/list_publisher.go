package publishers

import (
	"fmt"

	pb_list "github.com/sm888sm/halten-backend/list-service/api/pb"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

const (
	DeleteList MessageType = iota
	// Add other message types here...
)

type ListPublisher struct {
	Channel *amqp.Channel
}

func (p *ListPublisher) Publish(messageType MessageType, message []byte) error {
	switch messageType {
	case DeleteList:
		var msg pb_list.DeleteListRequest
		err := proto.Unmarshal(message, &msg)
		if err != nil {
			return err
		}

		err = p.publishDeleteListMessage(&msg)
		if err != nil {
			return err
		}
	// Add other cases for other message types here...
	default:
		return fmt.Errorf("invalid message type: %v", messageType)
	}

	return nil
}

func NewListPublisher(ch *amqp.Channel) *ListPublisher {
	return &ListPublisher{Channel: ch}
}

func (p *ListPublisher) publishDeleteListMessage(req *pb_list.DeleteListRequest) error {
	ch := p.Channel
	defer ch.Close()

	message, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"halten",
		"list.delete",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/protobuf",
			Body:        message,
		})
	return err
}
