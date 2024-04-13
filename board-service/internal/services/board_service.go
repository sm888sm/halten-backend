package services

import (
	"context"

	pb "github.com/sm888sm/halten-backend/board-service/api/pb"
	"github.com/sm888sm/halten-backend/board-service/internal/repositories"
	"github.com/sm888sm/halten-backend/common"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	"github.com/sm888sm/halten-backend/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BoardService struct {
	boardRepo repositories.BoardRepository
	pb.UnimplementedBoardServiceServer
}

func NewBoardService(repo repositories.BoardRepository) *BoardService {
	return &BoardService{boardRepo: repo}
}

func (s *BoardService) CreateBoard(ctx context.Context, req *pb.CreateBoardRequest) (*pb.CreateBoardResponse, error) {
	board := &models.Board{
		Name:   req.Name,
		UserID: uint(req.UserId),
	}
	board, err := s.boardRepo.CreateBoard(board, uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.CreateBoardResponse{
		Id:     uint64(board.ID),
		Name:   board.Name,
		UserId: uint64(board.UserID),
	}, nil
}

func (s *BoardService) GetBoardByID(ctx context.Context, req *pb.GetBoardByIDRequest) (*pb.GetBoardByIDResponse, error) {
	board, err := s.boardRepo.GetBoardByID(uint(req.Id), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	res := &pb.Board{
		Id:         uint64(board.ID),
		Name:       board.Name,
		UserId:     uint64(board.UserID),
		Visibility: board.Visibility,
	}
	return &pb.GetBoardByIDResponse{Board: res}, nil
}

func (s *BoardService) GetBoards(ctx context.Context, req *pb.GetBoardsRequest) (*pb.GetBoardsResponse, error) {
	boardList, err := s.boardRepo.GetBoards(uint(req.UserId), int(req.PageNumber), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	var boards []*pb.Board
	for _, b := range boardList.Boards {
		boards = append(boards, &pb.Board{
			Id:         uint64(b.ID),
			Name:       b.Name,
			UserId:     uint64(b.UserID),
			Visibility: b.Visibility,
		})
	}
	return &pb.GetBoardsResponse{
		Boards: boards,
		Pagination: &pb.Pagination{
			CurrentPage:  uint64(boardList.Pagination.CurrentPage),
			TotalPages:   uint64(boardList.Pagination.TotalPages),
			ItemsPerPage: uint64(boardList.Pagination.ItemsPerPage),
			TotalItems:   uint64(boardList.Pagination.TotalItems),
			HasMore:      boardList.Pagination.HasMore,
		},
	}, nil
}

func (s *BoardService) UpdateBoard(ctx context.Context, req *pb.UpdateBoardRequest) (*pb.UpdateBoardResponse, error) {
	err := s.boardRepo.UpdateBoard(uint(req.Id), uint(req.Id), req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateBoardResponse{
		Id:     req.Id,
		UserId: req.Id,
		Name:   req.Name,
	}, nil
}

func (s *BoardService) DeleteBoard(ctx context.Context, req *pb.DeleteBoardRequest) (*pb.DeleteBoardResponse, error) {
	err := s.boardRepo.DeleteBoard(uint(req.UserId), uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteBoardResponse{
		Message: "Board deleted successfully",
	}, nil
}

func (s *BoardService) AddBoardUsers(ctx context.Context, req *pb.AddUsersRequest) (*pb.AddUsersResponse, error) {
	if req.Role != common.MemberRole && req.Role != common.ObserverRole {
		return nil, status.Error(codes.InvalidArgument, errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid role. Role should be either 'member' or 'observer'").Error())
	}

	var userIds []uint
	for _, id := range req.UserIds {
		userIds = append(userIds, uint(id))
	}
	err := s.boardRepo.AddBoardUsers(uint(req.UserId), uint(req.Id), userIds, req.Role)
	if err != nil {
		return nil, err
	}
	return &pb.AddUsersResponse{
		Message: "Users added successfully",
	}, nil
}

func (s *BoardService) RemoveBoardUsers(ctx context.Context, req *pb.RemoveUsersRequest) (*pb.RemoveUsersResponse, error) {
	var userIds []uint
	for _, id := range req.UserIds {
		userIds = append(userIds, uint(id))
	}
	err := s.boardRepo.RemoveBoardUsers(uint(req.UserId), uint(req.Id), userIds)
	if err != nil {
		return nil, err
	}
	return &pb.RemoveUsersResponse{
		Message: "Users removed successfully",
	}, nil
}

func (s *BoardService) GetBoardUsers(ctx context.Context, req *pb.GetBoardUsersRequest) (*pb.GetBoardUsersResponse, error) {
	users, err := s.boardRepo.GetBoardUsers(uint(req.Id))
	if err != nil {
		return nil, err
	}
	var pbUsers []*pb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.User{
			UserId:   uint64(u.ID),
			UserName: u.Username,
			Role:     u.Email,
		})
	}
	return &pb.GetBoardUsersResponse{
		Users: pbUsers,
	}, nil
}

func (s *BoardService) AssignUserRoleBoard(ctx context.Context, req *pb.AssignUserRoleRequest) (*pb.AssignUserRoleResponse, error) {
	err := s.boardRepo.AssignUserRoleBoard(uint(req.UserId), uint(req.Id), uint(req.AssignUserId), req.Role)
	if err != nil {
		return nil, err
	}
	return &pb.AssignUserRoleResponse{
		Message: "Role assigned successfully",
	}, nil
}

func (s *BoardService) ChangeBoardOwner(ctx context.Context, req *pb.ChangeBoardOwnerRequest) (*pb.ChangeBoardOwnerResponse, error) {
	err := s.boardRepo.ChangeBoardOwner(uint(req.Id), uint(req.CurrentOwnerId), uint(req.NewOwnerId))
	if err != nil {
		return nil, err
	}
	return &pb.ChangeBoardOwnerResponse{
		Message: "Owner changed successfully",
	}, nil
}
