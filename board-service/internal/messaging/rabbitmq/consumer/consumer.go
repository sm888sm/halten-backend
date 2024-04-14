package consumer

import (
	"context"
	"fmt"

	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
	"github.com/sm888sm/halten-backend/board-service/internal/services"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type BoardConsumer struct {
	Channel      *amqp.Channel
	BoardService *services.BoardService
}

func NewBoardConsumer(ch *amqp.Channel, boardService *services.BoardService) *BoardConsumer {
	return &BoardConsumer{Channel: ch, BoardService: boardService}
}

func (c *BoardConsumer) ConsumeBoardMessages(ctx context.Context) error {
	ch := c.Channel

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	if err != nil {
		return err
	}

	// Bind the queue to multiple routing keys
	err = ch.QueueBind(
		q.Name,
		"board.*", // Use a wildboard to match any routing key that starts with "board."
		"halten",
		false,
		nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			fmt.Println("Received a message: ", d.RoutingKey)
			switch d.RoutingKey {
			case "board.delete":
				req := &pb_board.DeleteBoardRequest{}
				err := proto.Unmarshal(d.Body, req)
				if err != nil {
					// Log the error or handle it as needed
					continue
				}

				_, err = c.BoardService.DeleteBoard(ctx, req)
				if err != nil {
					// TODO: Log the error or handle it as needed
				}

				continue

			case "board.update":
				// Handle board.update messages
				// ...

				// Add more cases as needed
			default:
				// Log or handle unknown routing keys
				// ...
			}
		}
	}()

	return nil
}
