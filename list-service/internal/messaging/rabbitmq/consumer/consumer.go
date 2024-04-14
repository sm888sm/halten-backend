package consumer

import (
	"context"
	"fmt"

	pb_list "github.com/sm888sm/halten-backend/list-service/api/pb"
	"github.com/sm888sm/halten-backend/list-service/internal/services"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type ListConsumer struct {
	Channel     *amqp.Channel
	ListService *services.ListService
}

func NewListConsumer(ch *amqp.Channel, listService *services.ListService) *ListConsumer {
	return &ListConsumer{Channel: ch, ListService: listService}
}

func (c *ListConsumer) ConsumeListMessages(ctx context.Context) error {
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
		"list.*", // Use a wildlist to match any routing key that starts with "list."
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
			case "list.delete":
				req := &pb_list.DeleteListRequest{}
				err := proto.Unmarshal(d.Body, req)
				if err != nil {
					// Log the error or handle it as needed
					continue
				}

				_, err = c.ListService.DeleteList(ctx, req)
				if err != nil {
					// TODO: Log the error or handle it as needed
				}

				continue

			case "list.update":
				// Handle list.update messages
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
