package services

import (
	"context"

	pb "github.com/sm888sm/halten-backend/list-service/api/pb"
	"github.com/sm888sm/halten-backend/list-service/internal/repositories"
	models "github.com/sm888sm/halten-backend/models"
)

type ListService struct {
	listRepo repositories.ListRepository
	pb.UnimplementedListServiceServer
}

func NewListService(repo repositories.ListRepository) *ListService {
	return &ListService{listRepo: repo}
}

func (s *ListService) CreateList(ctx context.Context, req *pb.CreateListRequest) (*pb.CreateListResponse, error) {
	list := &models.List{
		Name:    req.Name,
		BoardID: uint(req.BoardId),
	}
	err := s.listRepo.CreateList(list, uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.CreateListResponse{
		Id:      uint64(list.ID),
		UserId:  req.UserId,
		BoardId: req.BoardId,
		Name:    list.Name,
	}, nil
}

func (s *ListService) GetListsByBoard(ctx context.Context, req *pb.GetListsByBoardRequest) (*pb.GetListsByBoardResponse, error) {
	lists, err := s.listRepo.GetListsByBoard(uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	var pb_lists []*pb.List
	for _, l := range lists {
		pb_lists = append(pb_lists, &pb.List{
			Id:       uint64(l.ID),
			BoardId:  uint64(l.BoardID),
			Name:     l.Name,
			Position: int32(l.Position),
		})
	}
	return &pb.GetListsByBoardResponse{
		Lists: pb_lists,
	}, nil
}

func (s *ListService) UpdateList(ctx context.Context, req *pb.UpdateListRequest) (*pb.UpdateListResponse, error) {
	err := s.listRepo.UpdateList(uint(req.Id), req.Name, uint(req.BoardId), uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.UpdateListResponse{
		Id:     req.Id,
		UserId: req.UserId,
		Name:   req.Name,
	}, nil
}

func (s *ListService) DeleteList(ctx context.Context, req *pb.DeleteListRequest) (*pb.DeleteListResponse, error) {
	err := s.listRepo.DeleteList(uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteListResponse{
		Message: true,
	}, nil
}

func (s *ListService) MoveListPosition(ctx context.Context, req *pb.MoveListPositionRequest) (*pb.MoveListPositionResponse, error) {
	err := s.listRepo.MoveListPosition(uint(req.Id), int(req.NewPosition), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.MoveListPositionResponse{
		Message: true,
	}, nil
}
