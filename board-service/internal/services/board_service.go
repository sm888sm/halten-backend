package services

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"

	external_services "github.com/sm888sm/halten-backend/board-service/external/services"

	"github.com/sm888sm/halten-backend/board-service/internal/repositories"

	"github.com/sm888sm/halten-backend/common/constants/contextkeys"
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
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	board := &models.Board{
		Name:   req.Name,
		UserID: userID,
	}

	// repoReq := &repositories.CreateBoardRequest{
	// 	Board: board,
	// }

	// repoRes, err := s.boardRepo.CreateBoard(repoReq)
	// if err != nil {
	// 	return nil, err
	// }

	// TODO : Create List

	// listService, err := s.services.GetListClient()
	// if err != nil {
	// 	return nil, errorhandler.NewGrpcInternalError()
	// }

	// listReq := &pb_list.CreateListRequest{
	// 	BoardID: repoRes.Board.ID,
	// 	Name:    "Default List",
	// }

	// listRes, err := listService.CreateList(ctx, listReq)
	// if err != nil {
	// 	// Delete board through repo
	// 	err = s.boardRepo.DeleteBoard(&repositories.DeleteBoardRequest{
	// 		BoardID: repoRes.Board.ID,
	// 	})

	// 	if err != nil {
	// 		return nil, errorhandler.NewGrpcInternalError()
	// 	}

	// 	return nil, errorhandler.NewGrpcInternalError()
	// }

	// TODO : Create Card

	// cardService, err := s.services.GetCardClient()
	// if err != nil {
	// 	// Delete board through repo
	// 	err = s.boardRepo.DeleteBoard(&repositories.DeleteBoardRequest{
	// 		BoardID: repoRes.Board.ID,
	// 	})
	// 	if err != nil {
	// 		return nil, errorhandler.NewGrpcInternalError()
	// 	}

	// 	return nil, errorhandler.NewGrpcInternalError()
	// }

	// cardReq := &pb_card.CreateCardRequest{
	// 	ListID: list.ID,
	// 	Name:   "Default Card",
	// }

	// _, err = cardService.CreateCard(ctx, cardReq)
	// if err != nil {
	// 	// Delete board through repo
	// 	err = s.boardRepo.DeleteBoard(&repositories.DeleteBoardRequest{
	// 		BoardID: repoRes.Board.ID,
	// 	})
	// 	if err != nil {
	// 		return nil, errorhandler.NewGrpcInternalError()
	// 	}

	// 	return nil, err
	// }

	return &pb_board.CreateBoardResponse{
		BoardID: board.ID,
		Name:    board.Name,
	}, nil
}

func (s *BoardService) GetBoardByID(ctx context.Context, req *pb_board.GetBoardByIDRequest) (*pb_board.GetBoardByIDResponse, error) {
	// Extract the user ID from the context

	// Call the repository function
	repoRes, err := s.boardRepo.GetBoardByID(&repositories.GetBoardByIDRequest{
		BoardID: req.BoardID,
	})
	if err != nil {
		return nil, err
	}

	// TODO get board members from user service

	// TODO get cards from card service

	// TODO get lists from list service

	// Convert the board model to a protobuf message

	board := repoRes.Board
	boardProto := &pb_board.Board{
		BoardID:    board.ID,
		UserID:     board.ID,
		Name:       board.Name,
		Visibility: board.Visibility,
		// Members:    convertPermissionsToProto(board.Members),
		// Lists:       convertListsToProto(board.Lists),
		// Cards:       convertCardsToProto(board.Cards),
		Labels:    convertLabelsToProto(repoRes.Board.Labels),
		CreatedAt: timestamppb.New(repoRes.Board.CreatedAt),
		UpdatedAt: timestamppb.New(repoRes.Board.UpdatedAt),
	}

	// Return the response
	return &pb_board.GetBoardByIDResponse{Board: boardProto}, nil
}

func (s *BoardService) GetBoardList(ctx context.Context, req *pb_board.GetBoardListRequest) (*pb_board.GetBoardListResponse, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
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
	repoReq := &repositories.GetBoardMembersRequest{
		BoardID: req.BoardID,
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
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
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

func (s *BoardService) UpdateBoardName(ctx context.Context, req *pb_board.UpdateBoardNameRequest) (*pb_board.UpdateBoardNameResponse, error) {

	err := s.boardRepo.UpdateBoardName(&repositories.UpdateBoardNameRequest{BoardID: req.BoardID, Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &pb_board.UpdateBoardNameResponse{
		Message: "Board name updated successfully",
	}, nil
}

func (s *BoardService) AddBoardUsers(ctx context.Context, req *pb_board.AddBoardUsersRequest) (*pb_board.AddBoardUsersResponse, error) {
	// Extract the board ID and user IDs from the request
	boardID := req.BoardID
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
	// Extract the board ID and user IDs from the request
	boardID := req.BoardID
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

func (s *BoardService) AssignBoardUserRole(ctx context.Context, req *pb_board.AssignBoardUserRoleRequest) (*pb_board.AssignBoardUserRoleResponse, error) {
	// Call the repository function to assign the user role to the board
	err := s.boardRepo.AssignBoardUserRole(&repositories.AssignBoardUserRoleRequest{
		BoardID: req.BoardID,
		UserID:  req.UserID,
		Role:    req.Role,
	})
	if err != nil {
		return nil, err
	}

	// Return a successful response
	return &pb_board.AssignBoardUserRoleResponse{
		Message: "User role assigned to the board successfully",
	}, nil
}

func (s *BoardService) ChangeBoardOwner(ctx context.Context, req *pb_board.ChangeBoardOwnerRequest) (*pb_board.ChangeBoardOwnerResponse, error) {

	err := s.boardRepo.ChangeBoardOwner(&repositories.ChangeBoardOwnerRequest{
		BoardID:    req.BoardID,
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
	// Call the repository function to change the visibility of the board
	err := s.boardRepo.ChangeBoardVisibility(&repositories.ChangeBoardVisibilityRequest{
		BoardID:    req.BoardID,
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
	label, err := s.boardRepo.AddLabel(&repositories.AddLabelRequest{
		BoardID: req.BoardID,
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
	err := s.boardRepo.RemoveLabel(&repositories.RemoveLabelRequest{
		BoardID: req.BoardID,
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
	// Call the repository function to archive the board
	err := s.boardRepo.ArchiveBoard(&repositories.ArchiveBoardRequest{
		BoardID: req.BoardID,
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
	// Call the repository function to restore the board
	err := s.boardRepo.RestoreBoard(&repositories.RestoreBoardRequest{
		BoardID: req.BoardID,
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
	// Call the repository function
	err := s.boardRepo.DeleteBoard(&repositories.DeleteBoardRequest{
		BoardID: req.BoardID,
	})
	if err != nil {
		return nil, err
	}

	// Return the response
	return &pb_board.DeleteBoardResponse{
		Message: "Board successfully deleted",
	}, nil
}

// Helpers

func convertLabelsToProto(labels []models.Label) []*pb_board.Label {
	var labelsProto []*pb_board.Label

	for _, label := range labels {
		labelProto := &pb_board.Label{
			LabelID: label.ID,
			Name:    label.Name,
			Color:   label.Color,
			BoardID: label.BoardID,
		}

		labelsProto = append(labelsProto, labelProto)
	}

	return labelsProto
}
