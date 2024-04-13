package services

import (
	"context"

	pb "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/sm888sm/halten-backend/card-service/internal/repositories"
	"github.com/sm888sm/halten-backend/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CardService struct {
	cardRepo repositories.CardRepository
	pb.UnimplementedCardServiceServer
}

func NewCardService(repo repositories.CardRepository) *CardService {
	return &CardService{cardRepo: repo}
}

func (s *CardService) CreateCard(ctx context.Context, req *pb.CreateCardRequest) (*pb.CreateCardResponse, error) {
	card := &models.Card{
		Name:   req.Name,
		ListID: uint(req.ListId),
	}
	err := s.cardRepo.CreateCard(card, uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.CreateCardResponse{
		Id: uint64(card.ID),
	}, nil
}

func (s *CardService) GetCardsByList(ctx context.Context, req *pb.GetCardsByListRequest) (*pb.GetCardsByListResponse, error) {
	cards, err := s.cardRepo.GetCardsByList(uint(req.ListId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	var pbCards []*pb.Card
	for _, c := range cards {
		pbCards = append(pbCards, &pb.Card{
			Id:          uint64(c.ID),
			ListId:      uint64(c.ListID),
			Name:        c.Name,
			Description: c.Description,
			Position:    int32(c.Position),
			StartDate:   timestamppb.New(*c.StartDate),
			DueDate:     timestamppb.New(*c.DueDate),
		})
	}
	return &pb.GetCardsByListResponse{
		Cards: pbCards,
	}, nil
}

func (s *CardService) UpdateCard(ctx context.Context, req *pb.UpdateCardRequest) (*pb.UpdateCardResponse, error) {
	err := s.cardRepo.UpdateCard(uint(req.Id), req.Name, uint(req.ListId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.UpdateCardResponse{}, nil
}

func (s *CardService) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
	err := s.cardRepo.DeleteCard(uint(req.Id), uint(req.ListId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCardResponse{}, nil
}

func (s *CardService) MoveCardPosition(ctx context.Context, req *pb.MoveCardPositionRequest) (*pb.MoveCardPositionResponse, error) {
	err := s.cardRepo.MoveCardPosition(uint(req.Id), int(req.NewPosition), uint(req.ListId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.MoveCardPositionResponse{}, nil
}
