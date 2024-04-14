package consumer

import (
	"context"
	"fmt"

	"github.com/sm888sm/halten-backend/user-service/internal/services"
	"github.com/streadway/amqp"
)

type UserConsumer struct {
	Channel     *amqp.Channel
	UserService *services.UserService
}

func NewUserConsumer(ch *amqp.Channel, userService *services.UserService) *UserConsumer {
	return &UserConsumer{Channel: ch, UserService: userService}
}

func (c *UserConsumer) ConsumeUserMessages(ctx context.Context) error {
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
		"user.*", // Use a wilduser to match any routing key that starts with "user."
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
			// case "user.delete":
			// 	req := &pb_user.DeleteUserRequest{}
			// 	err := proto.Unmarshal(d.Body, req)
			// 	if err != nil {
			// 		// Log the error or handle it as needed
			// 		continue
			// 	}

			// 	_, err = c.UserService.DeleteUser(ctx, req)
			// 	if err != nil {
			// 		// TODO: Log the error or handle it as needed
			// 	}

			// 	continue

			// case "user.update":
			// 	// Handle user.update messages
			// 	// ...

			// 	// Add more cases as needed
			default:
				// Log or handle unknown routing keys
				// ...
			}
		}
	}()

	return nil
}
