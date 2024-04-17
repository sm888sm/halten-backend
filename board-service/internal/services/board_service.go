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

	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/constants/roles"
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
		UserID: uint(req.UserID),
	}

	board, err := s.boardRepo.CreateBoard(board, uint(req.UserID))
	if err != nil {
		return nil, err
	}

	listService, err := s.services.GetListClient()
	if err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	listReq := &pb_list.CreateListRequest{
		BoardID: uint64(board.ID),
		Name:    "Default List",
	}

	listRes, err := listService.CreateList(ctx, listReq)
	if err != nil {
		// Delete board through repo
		err = s.boardRepo.DeleteBoard(uint(req.UserID), uint(listReq.BoardID))
		if err != nil {
			return nil, errorhandler.NewGrpcInternalError()
		}

		return nil, errorhandler.NewGrpcInternalError()
	}

	cardService, err := s.services.GetCardClient()
	if err != nil {
		// Delete board through repo
		err = s.boardRepo.DeleteBoard(uint(req.UserID), uint(listReq.BoardID))
		if err != nil {
			return nil, errorhandler.NewGrpcInternalError()
		}

		return nil, errorhandler.NewGrpcInternalError()
	}

	cardReq := &pb_card.CreateCardRequest{
		ListID:  listRes.ID,
		Name:    "Default Card",
		UserID:  uint64(req.UserID),
		BoardID: uint64(board.ID),
	}

	_, err = cardService.CreateCard(ctx, cardReq)
	if err != nil {
		deleteListReq := &pb_list.DeleteListRequest{
			ID:      listRes.ID,
			BoardID: listReq.BoardID,
			UserID:  uint64(req.UserID),
		}

		bytes, err := proto.Marshal(deleteListReq)
		if err != nil {
			return nil, errorhandler.NewGrpcInternalError()
		}

		// Delete board through repo
		err = s.boardRepo.DeleteBoard(uint(req.UserID), uint(listReq.BoardID))
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
		ID:     uint64(board.ID),
		Name:   board.Name,
		UserID: uint64(board.UserID),
	}, nil
}

func (s *BoardService) GetBoardByID(ctx context.Context, req *pb_board.GetBoardByIDRequest) (*pb_board.GetBoardByIDResponse, error) {
	// Extract the user ID from the context

	// Call the repository function
	board, err := s.boardRepo.GetBoardByID(uint(req.ID), uint(req.UserID))
	if err != nil {
		return nil, err

	}

	// TODO get cards from card service

	// TODO get lists from list service

	// Convert the board model to a protobuf message
	boardProto := &pb_board.Board{
		ID:          uint64(board.ID),
		UserID:      uint64(board.UserID),
		Name:        board.Name,
		Visibility:  board.Visibility,
		Permissions: convertPermissionsToProto(board.Permissions),
		// Lists:       convertListsToProto(board.Lists),
		// Cards:       convertCardsToProto(board.Cards),
		Labels:    convertLabelsToProto(board.Labels),
		CreatedAt: timestamppb.New(board.CreatedAt),
		UpdatedAt: timestamppb.New(board.UpdatedAt),
	}

	// Return the response
	return &pb_board.GetBoardByIDResponse{Board: boardProto}, nil
}

func (s *BoardService) GetBoardList(ctx context.Context, req *pb_board.GetBoardListRequest) (*pb_board.GetBoardListResponse, error) {
	boardList, err := s.boardRepo.GetBoardList(uint(req.UserID), int(req.PageNumber), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	var boards []*pb_board.BoardMeta
	for _, b := range boardList.Boards {
		boards = append(boards, &pb_board.BoardMeta{
			ID:         uint64(b.ID),
			Name:       b.Name,
			UserID:     uint64(b.UserID),
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
	err := s.boardRepo.UpdateBoard(uint(req.ID), uint(req.ID), req.Name)
	if err != nil {
		return nil, err
	}
	return &pb_board.UpdateBoardResponse{
		ID:     req.ID,
		UserID: req.ID,
		Name:   req.Name,
	}, nil
}

func (s *BoardService) DeleteBoard(ctx context.Context, req *pb_board.DeleteBoardRequest) (*pb_board.DeleteBoardResponse, error) {
	err := s.boardRepo.DeleteBoard(uint(req.UserID), uint(req.ID))
	if err != nil {
		return nil, err
	}
	return &pb_board.DeleteBoardResponse{
		Message: "Board deleted successfully",
	}, nil
}

func (s *BoardService) AddBoardUsers(ctx context.Context, req *pb_board.AddUsersRequest) (*pb_board.AddUsersResponse, error) {
	if req.Role != roles.MemberRole && req.Role != roles.ObserverRole {
		return nil, status.Error(codes.InvalidArgument, errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid role. Role should be either 'member' or 'observer'").Error())
	}

	var userIDs []uint
	for _, id := range req.UserIDs {
		userIDs = append(userIDs, uint(id))
	}
	err := s.boardRepo.AddBoardUsers(uint(req.UserID), uint(req.ID), userIDs, req.Role)
	if err != nil {
		return nil, err
	}
	return &pb_board.AddUsersResponse{
		Message: "Users added successfully",
	}, nil
}

func (s *BoardService) RemoveBoardUsers(ctx context.Context, req *pb_board.RemoveUsersRequest) (*pb_board.RemoveUsersResponse, error) {
	var userIDs []uint
	for _, id := range req.UserIDs {
		userIDs = append(userIDs, uint(id))
	}
	err := s.boardRepo.RemoveBoardUsers(uint(req.UserID), uint(req.ID), userIDs)
	if err != nil {
		return nil, err
	}
	return &pb_board.RemoveUsersResponse{
		Message: "Users removed successfully",
	}, nil
}

func (s *BoardService) GetBoardUsers(ctx context.Context, req *pb_board.GetBoardUsersRequest) (*pb_board.GetBoardUsersResponse, error) {
	users, err := s.boardRepo.GetBoardUsers(uint(req.ID))
	if err != nil {
		return nil, err
	}
	var pbUsers []*pb_board.User
	for _, u := range users {
		pbUsers = append(pbUsers, &pb_board.User{
			UserID:   uint64(u.ID),
			UserName: u.Username,
			Role:     u.Email,
		})
	}
	return &pb_board.GetBoardUsersResponse{
		Users: pbUsers,
	}, nil
}

func (s *BoardService) AssignBoardUserRole(ctx context.Context, req *pb_board.AssignUserRoleRequest) (*pb_board.AssignUserRoleResponse, error) {
	err := s.boardRepo.AssignBoardUserRole(uint(req.UserID), uint(req.ID), uint(req.AssignUserID), req.Role)
	if err != nil {
		return nil, err
	}
	return &pb_board.AssignUserRoleResponse{
		Message: "Role assigned successfully",
	}, nil
}

func (s *BoardService) ChangeBoardOwner(ctx context.Context, req *pb_board.ChangeBoardOwnerRequest) (*pb_board.ChangeBoardOwnerResponse, error) {
	err := s.boardRepo.ChangeBoardOwner(uint(req.ID), uint(req.CurrentOwnerID), uint(req.NewOwnerID))
	if err != nil {
		return nil, err
	}
	return &pb_board.ChangeBoardOwnerResponse{
		Message: "Owner changed successfully",
	}, nil
}

func (s *BoardService) GetArchivedBoardList(ctx context.Context, req *pb_board.GetBoardListRequest) (*pb_board.GetBoardListResponse, error) {
	boardList, err := s.boardRepo.GetArchivedBoardList(uint(req.UserID), int(req.PageNumber), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	var boards []*pb_board.BoardMeta
	for _, b := range boardList.Boards {
		boards = append(boards, &pb_board.BoardMeta{
			ID:         uint64(b.ID),
			Name:       b.Name,
			UserID:     uint64(b.UserID),
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
	err := s.boardRepo.RestoreBoard(uint(req.UserID), uint(req.ID))
	if err != nil {
		return nil, err
	}
	return &pb_board.RestoreBoardResponse{
		Message: "Board restored successfully",
	}, nil
}

func (s *BoardService) ArchiveBoard(ctx context.Context, req *pb_board.ArchiveBoardRequest) (*pb_board.ArchiveBoardResponse, error) {
	err := s.boardRepo.ArchiveBoard(uint(req.UserID), uint(req.ID))
	if err != nil {
		return nil, err
	}
	return &pb_board.ArchiveBoardResponse{
		Message: "Board archived successfully",
	}, nil
}

func (s *BoardService) GetBoardPermissionByCard(ctx context.Context, req *pb_board.GetBoardPermissionByCardRequest) (*pb_board.GetBoardPermissionByCardResponse, error) {
	result, err := s.boardRepo.GetBoardPermissionByCard(repositories.GetBoardPermissionByCardParams{CardID: uint(req.CardID), UserID: uint(req.UserID)})
	if err != nil {
		return nil, err
	}

	var permission *pb_board.Permission
	if result.Permission != nil {
		permission = &pb_board.Permission{
			PermissionID: uint64(result.Permission.ID),
			BoardID:      uint64(result.Permission.BoardID),
			UserID:       uint64(result.Permission.UserID),
			Role:         result.Permission.Role,
		}
	}

	return &pb_board.GetBoardPermissionByCardResponse{
		BoardID:    uint64(result.BoardID),
		Visibility: result.Visibility,
		Permission: permission,
	}, nil
}

func (s *BoardService) GetBoardPermissionByList(ctx context.Context, req *pb_board.GetBoardPermissionByListRequest) (*pb_board.GetBoardPermissionByListResponse, error) {
	result, err := s.boardRepo.GetBoardPermissionByList(repositories.GetBoardPermissionByListParams{ListID: uint(req.ListID), UserID: uint(req.UserID)})
	if err != nil {
		return nil, err
	}

	var permission *pb_board.Permission
	if result.Permission != nil {
		permission = &pb_board.Permission{
			PermissionID: uint64(result.Permission.ID),
			BoardID:      uint64(result.Permission.BoardID),
			UserID:       uint64(result.Permission.UserID),
			Role:         result.Permission.Role,
		}
	}

	return &pb_board.GetBoardPermissionByListResponse{
		BoardID:    uint64(result.BoardID),
		Visibility: result.Visibility,
		Permission: permission,
	}, nil
}

// Helpers

func convertPermissionsToProto(permissions []models.Permission) []*pb_board.Permission {
	var permissionsProto []*pb_board.Permission

	for _, permission := range permissions {
		permissionProto := &pb_board.Permission{
			PermissionID: uint64(permission.ID),
			BoardID:      uint64(permission.BoardID),
			UserID:       uint64(permission.UserID),
			Role:         permission.Role,
		}

		permissionsProto = append(permissionsProto, permissionProto)
	}

	return permissionsProto
}

func convertLabelsToProto(labels []models.Label) []*pb_board.Label {
	var labelsProto []*pb_board.Label

	for _, label := range labels {
		labelProto := &pb_board.Label{
			LabelID: uint64(label.ID),
			Name:    label.Name,
			Color:   label.Color,
			BoardID: uint64(label.BoardID),
		}

		labelsProto = append(labelsProto, labelProto)
	}

	return labelsProto
}
