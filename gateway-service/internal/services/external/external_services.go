package external

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
	var err error
	once.Do(func() {
		services = &Services{}

		// Set up a connection to each microservice server.
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		services.userConn, err = grpc.DialContext(ctx, cfg.UserServiceAddr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("Error connecting to UserService: %v", err)
		} else {
			services.userClient = pbUser.NewUserServiceClient(services.userConn)
			log.Printf("Successfully connected to UserService")
		}

		// Set up a connection to the AuthService.
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		services.authConn, err = grpc.DialContext(ctx, cfg.UserServiceAddr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("Error connecting to AuthService: %v", err)
		} else {
			services.authClient = pbUser.NewAuthServiceClient(services.authConn)
			log.Printf("Successfully connected to AuthService")
		}

		// Set up a connection to the BoardService.
		// ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		// defer cancel()
		// services.boardConn, err = grpc.DialContext(ctx, cfg.BoardServiceAddr, grpc.WithInsecure(), grpc.WithBlock())
		// if err != nil {
		// 	log.Printf("Error connecting to BoardService: %v", err)
		// } else {
		// 	services.boardClient = pbBoard.NewBoardServiceClient(services.boardConn)
		// 	log.Printf("Successfully connected to BoardService")
		// }

		// // Set up a connection to the ListService.
		// ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		// defer cancel()
		// services.listConn, err = grpc.DialContext(ctx, cfg.ListServiceAddr, grpc.WithInsecure(), grpc.WithBlock())
		// if err != nil {
		// 	log.Printf("Error connecting to ListService: %v", err)
		// } else {
		// 	services.listClient = pbList.NewListServiceClient(services.listConn)
		// 	log.Printf("Successfully connected to ListService")
		// }

		// // Set up a connection to the CardService.
		// ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		// defer cancel()
		// services.cardConn, err = grpc.DialContext(ctx, cfg.CardServiceAddr, grpc.WithInsecure(), grpc.WithBlock())
		// if err != nil {
		// 	log.Printf("Error connecting to CardService: %v", err)
		// } else {
		// 	services.cardClient = pbCard.NewCardServiceClient(services.cardConn)
		// 	log.Printf("Successfully connected to CardService")
		// }
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
	s.userConn.Close()
	s.authConn.Close()
	s.listConn.Close()
	s.boardConn.Close()
	s.cardConn.Close()
}
