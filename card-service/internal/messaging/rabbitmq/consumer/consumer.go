package consumer

import (
	"context"
	"fmt"

	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/sm888sm/halten-backend/card-service/internal/services"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type CardConsumer struct {
	Channel     *amqp.Channel
	CardService *services.CardService
}

func NewCardConsumer(ch *amqp.Channel, boardService *services.CardService) *CardConsumer {
	return &CardConsumer{Channel: ch, CardService: boardService}
}

func (c *CardConsumer) ConsumeCardMessages(ctx context.Context) error {
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
		"card.*", // Use a wildcard to match any routing key that starts with "board."
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
			case "card.delete":
				req := &pb_card.DeleteCardRequest{}
				err := proto.Unmarshal(d.Body, req)
				if err != nil {
					// Log the error or handle it as needed
					continue
				}

				_, err = c.CardService.DeleteCard(ctx, req)
				if err != nil {
					// TODO: Log the error or handle it as needed
				}

				continue

			case "card.update":
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
