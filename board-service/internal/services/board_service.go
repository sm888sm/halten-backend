package services

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	pb_list "github.com/sm888sm/halten-backend/list-service/api/pb"

	external_services "github.com/sm888sm/halten-backend/board-service/internal/services/external"

	"github.com/sm888sm/halten-backend/board-service/internal/repositories"

	"github.com/sm888sm/halten-backend/common"
	"github.com/sm888sm/halten-backend/common/errorhandler"

	"github.com/sm888sm/halten-backend/common/messaging/rabbitmq/publishers"
	"github.com/sm888sm/halten-backend/models"
)

type BoardService struct {
	boardRepo repositories.BoardRepository
	pb_board.UnimplementedBoardServiceServer
	services   *external_services.Services
	publishers publishers.Publishers
}

func NewBoardService(repo repositories.BoardRepository, publishers publishers.Publishers) *BoardService {
	return &BoardService{
		boardRepo:  repo,
		publishers: publishers,
	}
}

// TODO: Implement it with RabbitMQ
func (s *BoardService) CreateBoard(ctx context.Context, req *pb_board.CreateBoardRequest) (*pb_board.CreateBoardResponse, error) {
	board := &models.Board{
		Name:   req.Name,
		UserID: uint(req.UserId),
	}

	board, err := s.boardRepo.CreateBoard(board, uint(req.UserId))
	if err != nil {
		return nil, err
	}

	listService, err := s.services.GetListClient()
	if err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	listReq := &pb_list.CreateListRequest{
		BoardId: uint64(board.ID),
		Name:    "Default List",
	}

	listRes, err := listService.CreateList(ctx, listReq)
	if err != nil {
		return nil, err
	}

	cardService, err := s.services.GetCardClient()
	if err != nil {
		// Delete board through repo
		s.boardRepo.DeleteBoard(uint(req.UserId), uint(listReq.BoardId))
		return nil, errorhandler.NewGrpcInternalError()
	}

	cardReq := &pb_card.CreateCardRequest{
		ListId:  listRes.Id,
		Name:    "Default Card",
		UserId:  uint64(req.UserId),
		BoardId: uint64(board.ID),
	}

	_, err = cardService.CreateCard(ctx, cardReq)
	if err != nil {
		deleteListReq := &pb_list.DeleteListRequest{
			Id:      listRes.Id,
			BoardId: listReq.BoardId,
			UserId:  uint64(req.UserId),
		}

		bytes, err := proto.Marshal(deleteListReq)
		if err != nil {
			return nil, errorhandler.NewGrpcInternalError()
		}

		// Delete board through repo
		s.boardRepo.DeleteBoard(uint(req.UserId), uint(listReq.BoardId))

		// Delete board through publishers
		s.publishers.ListPublisher.Publish(publishers.DeleteList, bytes)
		return nil, err
	}

	return &pb_board.CreateBoardResponse{
		Id:     uint64(board.ID),
		Name:   board.Name,
		UserId: uint64(board.UserID),
	}, nil
}

func (s *BoardService) GetBoardByID(ctx context.Context, req *pb_board.GetBoardByIDRequest) (*pb_board.GetBoardByIDResponse, error) {
	board, err := s.boardRepo.GetBoardByID(uint(req.Id), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	res := &pb_board.Board{
		Id:         uint64(board.ID),
		Name:       board.Name,
		UserId:     uint64(board.UserID),
		Visibility: board.Visibility,
	}
	return &pb_board.GetBoardByIDResponse{Board: res}, nil
}

func (s *BoardService) GetBoards(ctx context.Context, req *pb_board.GetBoardsRequest) (*pb_board.GetBoardsResponse, error) {
	boardList, err := s.boardRepo.GetBoards(uint(req.UserId), int(req.PageNumber), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	var boards []*pb_board.Board
	for _, b := range boardList.Boards {
		boards = append(boards, &pb_board.Board{
			Id:         uint64(b.ID),
			Name:       b.Name,
			UserId:     uint64(b.UserID),
			Visibility: b.Visibility,
		})
	}
	return &pb_board.GetBoardsResponse{
		Boards: boards,
		Pagination: &pb_board.Pagination{
			CurrentPage:  uint64(boardList.Pagination.CurrentPage),
			TotalPages:   uint64(boardList.Pagination.TotalPages),
			ItemsPerPage: uint64(boardList.Pagination.ItemsPerPage),
			TotalItems:   uint64(boardList.Pagination.TotalItems),
			HasMore:      boardList.Pagination.HasMore,
		},
	}, nil
}

func (s *BoardService) UpdateBoard(ctx context.Context, req *pb_board.UpdateBoardRequest) (*pb_board.UpdateBoardResponse, error) {
	err := s.boardRepo.UpdateBoard(uint(req.Id), uint(req.Id), req.Name)
	if err != nil {
		return nil, err
	}
	return &pb_board.UpdateBoardResponse{
		Id:     req.Id,
		UserId: req.Id,
		Name:   req.Name,
	}, nil
}

func (s *BoardService) DeleteBoard(ctx context.Context, req *pb_board.DeleteBoardRequest) (*pb_board.DeleteBoardResponse, error) {
	err := s.boardRepo.DeleteBoard(uint(req.UserId), uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb_board.DeleteBoardResponse{
		Message: "Board deleted successfully",
	}, nil
}

func (s *BoardService) AddBoardUsers(ctx context.Context, req *pb_board.AddUsersRequest) (*pb_board.AddUsersResponse, error) {
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
	return &pb_board.AddUsersResponse{
		Message: "Users added successfully",
	}, nil
}

func (s *BoardService) RemoveBoardUsers(ctx context.Context, req *pb_board.RemoveUsersRequest) (*pb_board.RemoveUsersResponse, error) {
	var userIds []uint
	for _, id := range req.UserIds {
		userIds = append(userIds, uint(id))
	}
	err := s.boardRepo.RemoveBoardUsers(uint(req.UserId), uint(req.Id), userIds)
	if err != nil {
		return nil, err
	}
	return &pb_board.RemoveUsersResponse{
		Message: "Users removed successfully",
	}, nil
}

func (s *BoardService) GetBoardUsers(ctx context.Context, req *pb_board.GetBoardUsersRequest) (*pb_board.GetBoardUsersResponse, error) {
	users, err := s.boardRepo.GetBoardUsers(uint(req.Id))
	if err != nil {
		return nil, err
	}
	var pbUsers []*pb_board.User
	for _, u := range users {
		pbUsers = append(pbUsers, &pb_board.User{
			UserId:   uint64(u.ID),
			UserName: u.Username,
			Role:     u.Email,
		})
	}
	return &pb_board.GetBoardUsersResponse{
		Users: pbUsers,
	}, nil
}

func (s *BoardService) AssignUserRoleBoard(ctx context.Context, req *pb_board.AssignUserRoleRequest) (*pb_board.AssignUserRoleResponse, error) {
	err := s.boardRepo.AssignUserRoleBoard(uint(req.UserId), uint(req.Id), uint(req.AssignUserId), req.Role)
	if err != nil {
		return nil, err
	}
	return &pb_board.AssignUserRoleResponse{
		Message: "Role assigned successfully",
	}, nil
}

func (s *BoardService) ChangeBoardOwner(ctx context.Context, req *pb_board.ChangeBoardOwnerRequest) (*pb_board.ChangeBoardOwnerResponse, error) {
	err := s.boardRepo.ChangeBoardOwner(uint(req.Id), uint(req.CurrentOwnerId), uint(req.NewOwnerId))
	if err != nil {
		return nil, err
	}
	return &pb_board.ChangeBoardOwnerResponse{
		Message: "Owner changed successfully",
	}, nil
}
