package services

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/sm888sm/halten-backend/gateway-service/internal/config"

	pbBoard "github.com/sm888sm/halten-backend/board-service/api/pb"
	pbCard "github.com/sm888sm/halten-backend/card-service/api/pb"
	pbList "github.com/sm888sm/halten-backend/list-service/api/pb"
	pbUser "github.com/sm888sm/halten-backend/user-service/api/pb"
)

type Services struct {
	userClient  pbUser.UserServiceClient
	authClient  pbUser.AuthServiceClient
	boardClient pbBoard.BoardServiceClient
	listClient  pbList.ListServiceClient
	cardClient  pbCard.CardServiceClient

	userConn  *grpc.ClientConn
	authConn  *grpc.ClientConn
	boardConn *grpc.ClientConn
	listConn  *grpc.ClientConn
	cardConn  *grpc.ClientConn
}

var services *Services
var once sync.Once

func GetServices(cfg *config.ServiceConfig) *Services {
	once.Do(func() {
		services = &Services{}

		// Function to connect to a service
		connect := func(target string, setConn func(conn *grpc.ClientConn), setClient func(client interface{})) {
			var conn *grpc.ClientConn
			for {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				var err error
				conn, err = grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
				cancel()
				if err != nil {
					log.Printf("Unable to reconnect to service at %s: %v. Will attempt to reconnect in 1 second.", target, err)
					time.Sleep(1 * time.Second) // wait for a second before trying again
				} else {
					setConn(conn)
					break
				}
			}

			// Start a goroutine to monitor the connection
			go func() {
				for {
					time.Sleep(5 * time.Second) // check every 5 seconds
					if conn.GetState() != connectivity.Ready {
						// connection is not ready, attempt to reconnect
						for {
							ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
							newConn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
							cancel()
							if err != nil {
								log.Printf("Unable to reconnect to service at %s: %v. Will attempt to reconnect in 1 second.", target, err)
								time.Sleep(1 * time.Second) // wait for a second before trying again
							} else {
								conn.Close()                                                    // close the old connection
								conn = newConn                                                  // replace the old connection with the new one
								setConn(newConn)                                                // update the connection
								log.Printf("Successfully reconnected to service at %s", target) // log the reconnection
								break
							}
						}
					}
				}
			}()
		}

		// Set up connections to the services
		go connect(cfg.UserServiceAddr, func(conn *grpc.ClientConn) {
			services.userConn = conn
		}, func(client interface{}) {
			services.userClient = client.(pbUser.UserServiceClient)
		})

		go connect(cfg.UserServiceAddr, func(conn *grpc.ClientConn) {
			services.authConn = conn
		}, func(client interface{}) {
			services.authClient = client.(pbUser.AuthServiceClient)
		})

		go connect(cfg.BoardServiceAddr, func(conn *grpc.ClientConn) {
			services.boardConn = conn
		}, func(client interface{}) {
			services.boardClient = client.(pbBoard.BoardServiceClient)
		})

		go connect(cfg.ListServiceAddr, func(conn *grpc.ClientConn) {
			services.listConn = conn
		}, func(client interface{}) {
			services.listClient = client.(pbList.ListServiceClient)
		})

		go connect(cfg.CardServiceAddr, func(conn *grpc.ClientConn) {
			services.cardConn = conn
		}, func(client interface{}) {
			services.cardClient = client.(pbCard.CardServiceClient)
		})
	})

	return services
}

func (s *Services) GetUserClient() (pbUser.UserServiceClient, error) {
	if s.userConn.GetState() != connectivity.Ready {
		return nil, errors.New("user service not available")
	}
	return s.userClient, nil
}

func (s *Services) GetAuthClient() (pbUser.AuthServiceClient, error) {
	if s.authConn.GetState() != connectivity.Ready {
		return nil, errors.New("auth service not available")
	}
	return s.authClient, nil
}

func (s *Services) GetListClient() (pbList.ListServiceClient, error) {
	if s.listConn.GetState() != connectivity.Ready {
		return nil, errors.New("list service not available")
	}
	return s.listClient, nil
}

func (s *Services) GetBoardClient() (pbBoard.BoardServiceClient, error) {
	if s.boardConn.GetState() != connectivity.Ready {
		return nil, errors.New("board service not available")
	}
	return s.boardClient, nil
}

func (s *Services) GetCardClient() (pbCard.CardServiceClient, error) {
	if s.cardConn.GetState() != connectivity.Ready {
		return nil, errors.New("card service not available")
	}
	return s.cardClient, nil
}

func (s *Services) Close() {
	if s.userConn != nil {
		s.userConn.Close()
	}
	if s.authConn != nil {
		s.authConn.Close()
	}
	if s.listConn != nil {
		s.listConn.Close()
	}
	if s.boardConn != nil {
		s.boardConn.Close()
	}
	if s.cardConn != nil {
		s.cardConn.Close()
	}
}
