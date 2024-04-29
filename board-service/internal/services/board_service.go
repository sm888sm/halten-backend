package services

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	pb_list "github.com/sm888sm/halten-backend/list-service/api/pb"

	external_services "github.com/sm888sm/halten-backend/board-service/external/services"

	"github.com/sm888sm/halten-backend/board-service/internal/repositories"

	"github.com/sm888sm/halten-backend/common/constants/contextkeys"
	"github.com/sm888sm/halten-backend/common/errorhandlers"
	"github.com/sm888sm/halten-backend/common/helpers"

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
	userID, err := helpers.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	board := &models.Board{
		Name:   req.Name,
		UserID: userID,
	}

	repoReq := &repositories.CreateBoardRequest{
		Board: board,
	}

	repoRes, err := s.boardRepo.CreateBoard(repoReq)
	if err != nil {
		return nil, err
	}

	listService, err := s.services.GetListClient()
	if err != nil {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	listReq := &pb_list.CreateListRequest{
		List: &pb_list.List{
			BoardID: repoRes.Board.ID,
			Name:    "Default List",
		}}

	listRes, err := listService.CreateList(ctx, listReq)
	if err != nil {
		// Delete board through repo
		err = s.boardRepo.DeleteBoard(&repositories.DeleteBoardRequest{
			BoardID: repoRes.Board.ID,
		})

		if err != nil {
			return nil, errorhandlers.NewGrpcInternalError()
		}

		return nil, errorhandlers.NewGrpcInternalError()
	}

	cardService, err := s.services.GetCardClient()
	if err != nil {
		// Delete board through repo
		err = s.boardRepo.DeleteBoard(&repositories.DeleteBoardRequest{
			BoardID: repoRes.Board.ID,
		})
		if err != nil {
			return nil, errorhandlers.NewGrpcInternalError()
		}

		return nil, errorhandlers.NewGrpcInternalError()
	}

	cardReq := &pb_card.CreateCardRequest{
		ListID: listRes.List.ListID,
		Name:   "Default Card",
	}

	_, err = cardService.CreateCard(ctx, cardReq)
	if err != nil {
		// Delete board through repo
		err = s.boardRepo.DeleteBoard(&repositories.DeleteBoardRequest{
			BoardID: repoRes.Board.ID,
		})
		if err != nil {
			return nil, errorhandlers.NewGrpcInternalError()
		}

		return nil, err
	}

	return &pb_board.CreateBoardResponse{
		BoardID: board.ID,
		Name:    board.Name,
	}, nil
}

func (s *BoardService) GetBoardByID(ctx context.Context, req *pb_board.GetBoardByIDRequest) (*pb_board.GetBoardByIDResponse, error) {
	// Extract the board ID from the context
	boardID, err := helpers.ExtractBoardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Call the repository function
	repoRes, err := s.boardRepo.GetBoardByID(&repositories.GetBoardByIDRequest{
		BoardID: boardID,
	})
	if err != nil {
		return nil, err
	}

	listService, err := s.services.GetListClient()
	if err != nil {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	cardService, err := s.services.GetCardClient()
	if err != nil {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	listReq := &pb_list.GetListsByBoardRequest{}

	cardReq := &pb_card.GetCardsByBoardRequest{}

	listRes, err := listService.GetListsByBoard(ctx, listReq)
	if err != nil {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	cardRes, err := cardService.GetCardsByBoard(ctx, cardReq)
	if err != nil {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	// Convert the board model to a protobuf message

	board := repoRes.Board
	boardProto := &pb_board.Board{
		BoardID:    board.ID,
		UserID:     board.ID,
		Name:       board.Name,
		Visibility: board.Visibility,
		Members:    convertMembersToProto(board.Members),
		Lists:      convertListsToProto(listRes.Lists),
		Cards:      convertCardsToProto(cardRes.Cards),
		Labels:     convertLabelsToProto(board.Labels),
		CreatedAt:  timestamppb.New(repoRes.Board.CreatedAt),
		UpdatedAt:  timestamppb.New(repoRes.Board.UpdatedAt),
	}

	// Return the response
	return &pb_board.GetBoardByIDResponse{Board: boardProto}, nil
}

func (s *BoardService) GetBoardList(ctx context.Context, req *pb_board.GetBoardListRequest) (*pb_board.GetBoardListResponse, error) {
	userID, err := helpers.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	boardList, err := s.boardRepo.GetBoardList(&repositories.GetBoardListRequest{
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
		UserID:     userID,
	})

	if err != nil {
		return nil, err
	}

	var boards []*pb_board.BoardMeta
	for _, board := range boardList.Boards {
		boards = append(boards, &pb_board.BoardMeta{
			BoardID:    board.ID,
			Name:       board.Name,
			Visibility: board.Visibility,
			CreatedAt:  timestamppb.New(board.CreatedAt),
			UpdatedAt:  timestamppb.New(board.UpdatedAt),
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
func (s *BoardService) GetBoardMembers(ctx context.Context, req *pb_board.GetBoardMembersRequest) (*pb_board.GetBoardMembersResponse, error) {
	boardID, err := helpers.ExtractBoardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	repoReq := &repositories.GetBoardMembersRequest{
		BoardID: boardID,
	}

	repoRes, err := s.boardRepo.GetBoardMembers(repoReq)
	if err != nil {
		return nil, err
	}

	res := &pb_board.GetBoardMembersResponse{
		Members: make([]*pb_board.BoardMember, len(repoRes.Members)),
	}

	for i, member := range repoRes.Members {
		res.Members[i] = &pb_board.BoardMember{
			UserID:   member.ID,
			Username: member.Username,
			Role:     member.Role,
			Fullname: member.Fullname,
		}
	}

	return res, nil
}

func (s *BoardService) GetArchivedBoardList(ctx context.Context, req *pb_board.GetArchivedBoardListRequest) (*pb_board.GetArchivedBoardListResponse, error) {
	userID, err := helpers.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	boardList, err := s.boardRepo.GetArchivedBoardList(&repositories.GetArchivedBoardListRequest{
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
		UserID:     userID,
	})

	if err != nil {
		return nil, err
	}

	var boards []*pb_board.BoardMeta
	for _, board := range boardList.Boards {
		boards = append(boards, &pb_board.BoardMeta{
			BoardID:    board.ID,
			Name:       board.Name,
			Visibility: board.Visibility,
			CreatedAt:  timestamppb.New(board.CreatedAt),
			UpdatedAt:  timestamppb.New(board.UpdatedAt),
		})
	}
	return &pb_board.GetArchivedBoardListResponse{
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

// Authorization Required

func (s *BoardService) UpdateBoardName(ctx context.Context, req *pb_board.UpdateBoardNameRequest) (*pb_board.UpdateBoardNameResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	err := s.boardRepo.UpdateBoardName(&repositories.UpdateBoardNameRequest{BoardID: boardID, Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &pb_board.UpdateBoardNameResponse{
		Message: "Board name updated successfully",
	}, nil
}

func (s *BoardService) AddBoardUsers(ctx context.Context, req *pb_board.AddBoardUsersRequest) (*pb_board.AddBoardUsersResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	// Extract the board ID and user IDs from the request
	userIDs := req.UserIDs

	// Call the repository function to add the users to the board
	err := s.boardRepo.AddBoardUsers(&repositories.AddBoardUsersRequest{
		BoardID: boardID,
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}

	// Return a successful response
	return &pb_board.AddBoardUsersResponse{
		Message: "Users added to the board successfully",
	}, nil
}

func (s *BoardService) RemoveBoardUsers(ctx context.Context, req *pb_board.RemoveBoardUsersRequest) (*pb_board.RemoveBoardUsersResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	// Extract the board ID and user IDs from the request
	userIDs := req.UserIDs

	// Call the repository function to remove the users from the board
	err := s.boardRepo.RemoveBoardUsers(&repositories.RemoveBoardUsersRequest{
		BoardID: boardID,
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}

	// Return a successful response
	return &pb_board.RemoveBoardUsersResponse{
		Message: "Users removed from the board successfully",
	}, nil
}

func (s *BoardService) AssignBoardUsersRole(ctx context.Context, req *pb_board.AssignBoardUsersRoleRequest) (*pb_board.AssignBoardUsersRoleResponse, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	// Call the repository function to assign the user role to the board
	err := s.boardRepo.AssignBoardUsersRole(&repositories.AssignBoardUsersRoleRequest{
		BoardID: boardID,
		UserID:  userID,
		UserIDs: req.UserIDs,
		Role:    req.Role,
	})
	if err != nil {
		return nil, err
	}

	// Return a successful response
	return &pb_board.AssignBoardUsersRoleResponse{
		Message: "User role assigned to the board successfully",
	}, nil
}

func (s *BoardService) ChangeBoardOwner(ctx context.Context, req *pb_board.ChangeBoardOwnerRequest) (*pb_board.ChangeBoardOwnerResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	err := s.boardRepo.ChangeBoardOwner(&repositories.ChangeBoardOwnerRequest{
		BoardID:    boardID,
		NewOwnerID: req.NewOwnerID,
	})
	if err != nil {
		return nil, err
	}

	// Return a successful response
	return &pb_board.ChangeBoardOwnerResponse{
		Message: "Board owner changed successfully",
	}, nil
}

func (s *BoardService) ChangeBoardVisibility(ctx context.Context, req *pb_board.ChangeBoardVisibilityRequest) (*pb_board.ChangeBoardVisibilityResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	// Call the repository function to change the visibility of the board
	err := s.boardRepo.ChangeBoardVisibility(&repositories.ChangeBoardVisibilityRequest{
		BoardID:    boardID,
		Visibility: req.Visibility,
	})
	if err != nil {
		return nil, err
	}

	// Return a successful response
	return &pb_board.ChangeBoardVisibilityResponse{
		Message: "Board visibility changed successfully",
	}, nil
}

func (s *BoardService) AddLabel(ctx context.Context, req *pb_board.AddLabelRequest) (*pb_board.AddLabelResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	label, err := s.boardRepo.AddLabel(&repositories.AddLabelRequest{
		BoardID: boardID,
		Color:   req.Color,
		Name:    req.Name,
	})
	if err != nil {
		return nil, err
	}

	// Return a successful response
	return &pb_board.AddLabelResponse{
		Label: &pb_board.Label{
			LabelID: label.ID,
			Color:   label.Color,
			Name:    label.Name,
		},
	}, nil
}

func (s *BoardService) RemoveLabel(ctx context.Context, req *pb_board.RemoveLabelRequest) (*pb_board.RemoveLabelResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	err := s.boardRepo.RemoveLabel(&repositories.RemoveLabelRequest{
		BoardID: boardID,
		LabelID: req.LabelID,
	})
	if err != nil {
		return nil, err
	}

	// Return a successful response
	return &pb_board.RemoveLabelResponse{
		Message: "Label removed successfully",
	}, nil
}

func (s *BoardService) ArchiveBoard(ctx context.Context, req *pb_board.ArchiveBoardRequest) (*pb_board.ArchiveBoardResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	// Call the repository function to archive the board
	err := s.boardRepo.ArchiveBoard(&repositories.ArchiveBoardRequest{
		BoardID: boardID,
	})
	if err != nil {
		return nil, err
	}

	// Return a successful response
	return &pb_board.ArchiveBoardResponse{
		Message: "Board archived successfully",
	}, nil
}

func (s *BoardService) RestoreBoard(ctx context.Context, req *pb_board.RestoreBoardRequest) (*pb_board.RestoreBoardResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	// Call the repository function to restore the board
	err := s.boardRepo.RestoreBoard(&repositories.RestoreBoardRequest{
		BoardID: boardID,
	})
	if err != nil {
		return nil, err
	}

	// Return a successful response
	return &pb_board.RestoreBoardResponse{
		Message: "Board restored successfully",
	}, nil
}

func (s *BoardService) DeleteBoard(ctx context.Context, req *pb_board.DeleteBoardRequest) (*pb_board.DeleteBoardResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	// Call the repository function
	err := s.boardRepo.DeleteBoard(&repositories.DeleteBoardRequest{
		BoardID: boardID,
	})
	if err != nil {
		return nil, err
	}

	// Return the response
	return &pb_board.DeleteBoardResponse{
		Message: "Board successfully deleted",
	}, nil
}

func (s *BoardService) GetBoardIDByList(ctx context.Context, req *pb_board.GetBoardIDByListRequest) (*pb_board.GetBoardIDByListResponse, error) {
	boardID, err := s.boardRepo.GetBoardIDByList(&repositories.GetBoardIDByListRequest{ListID: req.ListID})
	if err != nil {
		return nil, err
	}

	return &pb_board.GetBoardIDByListResponse{
		BoardID: boardID,
	}, nil
}

func (s *BoardService) GetBoardIDByCard(ctx context.Context, req *pb_board.GetBoardIDByCardRequest) (*pb_board.GetBoardIDByCardResponse, error) {
	boardID, err := s.boardRepo.GetBoardIDByCard(&repositories.GetBoardIDByCardRequest{CardID: req.CardID})
	if err != nil {
		return nil, err
	}

	return &pb_board.GetBoardIDByCardResponse{
		BoardID: boardID,
	}, nil
}
