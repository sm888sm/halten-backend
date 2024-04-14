package publishers

import (
	"fmt"

	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type BoardPublisher struct {
	Channel *amqp.Channel
}

func NewBoardPublisher(ch *amqp.Channel) *BoardPublisher {
	return &BoardPublisher{Channel: ch}
}

func (p *BoardPublisher) Publish(messageType MessageType, message []byte) error {
	switch messageType {
	case DeleteCard:
		var msg pb_board.DeleteBoardRequest
		err := proto.Unmarshal(message, &msg)
		if err != nil {
			return err
		}

		err = p.publishDeleteBoardMessage(&msg)
		if err != nil {
			return err
		}
	// Add other cases for other message types here...
	default:
		return fmt.Errorf("invalid message type: %v", messageType)
	}

	return nil
}

func (p *BoardPublisher) publishDeleteBoardMessage(req *pb_board.DeleteBoardRequest) error {
	ch := p.Channel
	defer ch.Close()

	message, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"halten",
		"board.delete",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/protobuf",
			Body:        message,
		})
	return err
}
