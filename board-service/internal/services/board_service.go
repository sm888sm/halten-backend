package services

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	pb_list "github.com/sm888sm/halten-backend/list-service/api/pb"

	external_services "github.com/sm888sm/halten-backend/board-service/external/services"

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
	publishers *publishers.Publishers
}

func NewBoardService(repo repositories.BoardRepository, services *external_services.Services, publishers *publishers.Publishers) *BoardService {
	return &BoardService{
		boardRepo:  repo,
		publishers: publishers,
	}
}

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
		// Delete board through repo
		err = s.boardRepo.DeleteBoard(uint(req.UserId), uint(listReq.BoardId))
		if err != nil {
			return nil, errorhandler.NewGrpcInternalError()
		}

		return nil, errorhandler.NewGrpcInternalError()
	}

	cardService, err := s.services.GetCardClient()
	if err != nil {
		// Delete board through repo
		err = s.boardRepo.DeleteBoard(uint(req.UserId), uint(listReq.BoardId))
		if err != nil {
			return nil, errorhandler.NewGrpcInternalError()
		}

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
		err = s.boardRepo.DeleteBoard(uint(req.UserId), uint(listReq.BoardId))
		if err != nil {
			return nil, errorhandler.NewGrpcInternalError()
		}

		// Delete board through publishers
		err = s.publishers.ListPublisher.Publish(publishers.DeleteList, bytes)
		if err != nil {
			return nil, errorhandler.NewGrpcInternalError()
		}

		return nil, err
	}

	return &pb_board.CreateBoardResponse{
		Id:     uint64(board.ID),
		Name:   board.Name,
		UserId: uint64(board.UserID),
	}, nil
}

func (s *BoardService) GetBoardByID(ctx context.Context, req *pb_board.GetBoardByIDRequest) (*pb_board.GetBoardByIDResponse, error) {
	// Extract the user ID from the context

	// Call the repository function
	board, err := s.boardRepo.GetBoardByID(uint(req.Id), uint(req.UserId))
	if err != nil {
		return nil, err

	}

	// Convert the board model to a protobuf message
	boardProto := &pb_board.Board{
		Id:          uint64(board.ID),
		Name:        board.Name,
		Visibility:  board.Visibility,
		UserId:      uint64(board.UserID),
		CreatedAt:   timestamppb.New(board.CreatedAt),
		UpdatedAt:   timestamppb.New(board.UpdatedAt),
		Permissions: convertPermissionsToProto(board.Permissions),
		Lists:       convertListsToProto(board.Lists),
		Cards:       convertCardsToProto(board.Cards),
		Labels:      convertLabelsToProto(board.Labels),
	}

	// Return the response
	return &pb_board.GetBoardByIDResponse{Board: boardProto}, nil
}

func (s *BoardService) GetBoardList(ctx context.Context, req *pb_board.GetBoardListRequest) (*pb_board.GetBoardListResponse, error) {
	boardList, err := s.boardRepo.GetBoardList(uint(req.UserId), int(req.PageNumber), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	var boards []*pb_board.BoardMeta
	for _, b := range boardList.Boards {
		boards = append(boards, &pb_board.BoardMeta{
			Id:         uint64(b.ID),
			Name:       b.Name,
			UserId:     uint64(b.UserID),
			Visibility: b.Visibility,
			CreatedAt:  timestamppb.New(b.CreatedAt),
			UpdatedAt:  timestamppb.New(b.UpdatedAt),
		})
	}
	return &pb_board.GetBoardListResponse{
		Boards: boards,
		Pagination: &pb_board.Pagination{
			CurrentPage:  boardList.Pagination.CurrentPage,
			TotalPages:   boardList.Pagination.TotalPages,
			ItemsPerPage: boardList.Pagination.ItemsPerPage,
			TotalItems:   boardList.Pagination.TotalItems,
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

func (s *BoardService) AssignBoardUserRole(ctx context.Context, req *pb_board.AssignUserRoleRequest) (*pb_board.AssignUserRoleResponse, error) {
	err := s.boardRepo.AssignBoardUserRole(uint(req.UserId), uint(req.Id), uint(req.AssignUserId), req.Role)
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

func (s *BoardService) GetArchivedBoardList(ctx context.Context, req *pb_board.GetBoardListRequest) (*pb_board.GetBoardListResponse, error) {
	boardList, err := s.boardRepo.GetArchivedBoardList(uint(req.UserId), int(req.PageNumber), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	var boards []*pb_board.BoardMeta
	for _, b := range boardList.Boards {
		boards = append(boards, &pb_board.BoardMeta{
			Id:         uint64(b.ID),
			Name:       b.Name,
			UserId:     uint64(b.UserID),
			Visibility: b.Visibility,
			CreatedAt:  timestamppb.New(b.CreatedAt),
			UpdatedAt:  timestamppb.New(b.UpdatedAt),
		})
	}
	return &pb_board.GetBoardListResponse{
		Boards: boards,
		Pagination: &pb_board.Pagination{
			CurrentPage:  boardList.Pagination.CurrentPage,
			TotalPages:   boardList.Pagination.TotalPages,
			ItemsPerPage: boardList.Pagination.ItemsPerPage,
			TotalItems:   boardList.Pagination.TotalItems,
			HasMore:      boardList.Pagination.HasMore,
		},
	}, nil
}

func (s *BoardService) RestoreBoard(ctx context.Context, req *pb_board.RestoreBoardRequest) (*pb_board.RestoreBoardResponse, error) {
	err := s.boardRepo.RestoreBoard(uint(req.UserId), uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb_board.RestoreBoardResponse{
		Message: "Board restored successfully",
	}, nil
}

func (s *BoardService) ArchiveBoard(ctx context.Context, req *pb_board.ArchiveBoardRequest) (*pb_board.ArchiveBoardResponse, error) {
	err := s.boardRepo.ArchiveBoard(uint(req.UserId), uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb_board.ArchiveBoardResponse{
		Message: "Board archived successfully",
	}, nil
}

func convertPermissionsToProto(permissions []models.Permission) []*pb_board.Permission {
	var permissionsProto []*pb_board.Permission

	for _, permission := range permissions {
		permissionProto := &pb_board.Permission{
			Id:      uint64(permission.ID),
			BoardId: uint64(permission.BoardID),
			UserId:  uint64(permission.UserID),
			Role:    permission.Role,
		}

		permissionsProto = append(permissionsProto, permissionProto)
	}

	return permissionsProto
}

func convertListsToProto(lists []models.List) []*pb_board.List {
	var listsProto []*pb_board.List

	for _, list := range lists {
		listProto := &pb_board.List{
			Id:       uint64(list.ID),
			BoardId:  uint64(list.BoardID),
			Name:     list.Name,
			Position: int32(list.Position),
		}

		listsProto = append(listsProto, listProto)
	}

	return listsProto
}

func convertCardsToProto(cards []models.Card) []*pb_board.Card {
	var cardsProto []*pb_board.Card

	for _, card := range cards {
		cardProto := &pb_board.Card{
			Id:          uint64(card.ID),
			BoardId:     uint64(card.BoardID),
			ListId:      uint64(card.ListID),
			Name:        card.Name,
			Description: card.Description,
			Position:    int32(card.Position),
			StartDate:   timestamppb.New(*card.StartDate),
			DueDate:     timestamppb.New(*card.DueDate),
		}

		cardsProto = append(cardsProto, cardProto)
	}

	return cardsProto
}

func convertLabelsToProto(labels []models.Label) []*pb_board.Label {
	var labelsProto []*pb_board.Label

	for _, label := range labels {
		labelProto := &pb_board.Label{
			Id:      uint64(label.ID),
			Name:    label.Name,
			Color:   label.Color,
			BoardId: uint64(label.BoardID),
		}

		labelsProto = append(labelsProto, labelProto)
	}

	return labelsProto
}
