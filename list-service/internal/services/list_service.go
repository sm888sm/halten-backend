package services

import (
	"context"

	"github.com/sm888sm/halten-backend/common/constants/contextkeys"
	"github.com/sm888sm/halten-backend/common/errorhandler"
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
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	list := &models.List{
		BoardID: boardID,
		Name:    req.List.Name,
	}

	createListReq := &repositories.CreateListRequest{
		List:   list,
		UserID: userID,
	}
	resp, err := s.listRepo.CreateList(createListReq)
	if err != nil {
		return nil, err
	}
	return &pb.CreateListResponse{
		List: &pb.List{
			ListID:  resp.List.ID,
			BoardID: resp.List.BoardID,
			Name:    resp.List.Name,
		},
	}, nil
}

func (s *ListService) GetListsByID(ctx context.Context, req *pb.GetListByIDRequest) (*pb.GetListByIDResponse, error) {
	// Extract UserID and BoardID from the context
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	// Create a GetListsByBoardRequest object
	getListsReq := &repositories.GetListsByBoardRequest{
		UserID:  userID,
		BoardID: boardID,
	}

	// Call the GetListsByBoard method of the ListRepository interface
	getListsRes, err := s.listRepo.GetListsByBoard(getListsReq)
	if err != nil {
		return nil, err
	}

	// Convert the lists from the repository response to protobuf lists
	pbLists := make([]*pb.List, len(getListsRes.Lists))
	for i, list := range getListsRes.Lists {
		pbLists[i] = &pb.List{
			ListID:   list.ID,
			BoardID:  list.BoardID,
			Name:     list.Name,
			Position: int32(list.Position),
		}
	}

	// Return the response
	return &pb.GetListByIDResponse{
		Lists: pbLists,
	}, nil
}

func (s *ListService) GetListsByBoard(ctx context.Context, req *pb.GetListsByBoardRequest) (*pb.GetListsByBoardResponse, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	repoReq := &repositories.GetListsByBoardRequest{
		BoardID: boardID,
		UserID:  userID,
	}

	repoRes, err := s.listRepo.GetListsByBoard(repoReq)
	if err != nil {
		return nil, err
	}

	res := &pb.GetListsByBoardResponse{
		Lists: make([]*pb.List, len(repoRes.Lists)),
	}

	for i, list := range repoRes.Lists {
		res.Lists[i] = &pb.List{
			ListID:   list.ID,
			BoardID:  list.BoardID,
			Name:     list.Name,
			Position: int32(list.Position),
		}
	}

	return res, nil
}

func (s *ListService) UpdateListName(ctx context.Context, req *pb.UpdateListNameRequest) (*pb.UpdateListNameResponse, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	updateReq := &repositories.UpdateListNameRequest{
		ID:      req.ListID,
		Name:    req.Name,
		BoardID: boardID,
		UserID:  userID,
	}

	if err := s.listRepo.UpdateListName(updateReq); err != nil {
		return nil, err
	}

	return &pb.UpdateListNameResponse{
		Message: "List name updated successfully",
	}, nil
}

func (s *ListService) MoveListPosition(ctx context.Context, req *pb.MoveListPositionRequest) (*pb.MoveListPositionResponse, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	moveListPositionReq := &repositories.MoveListPositionRequest{
		ID:          req.ListID,
		NewPosition: int64(req.NewPosition),
		BoardID:     boardID,
		UserID:      userID,
	}

	if err := s.listRepo.MoveListPosition(moveListPositionReq); err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	return &pb.MoveListPositionResponse{
		Message: "List position updated successfully",
	}, nil
}

func (s *ListService) ArchiveList(ctx context.Context, req *pb.ArchiveListRequest) (*pb.ArchiveListResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	archiveListReq := &repositories.ArchiveListRequest{
		ListID:  req.ListID,
		BoardID: boardID,
	}

	if err := s.listRepo.ArchiveList(archiveListReq); err != nil {
		return nil, err
	}

	return &pb.ArchiveListResponse{
		Message: "List archived successfully",
	}, nil
}

func (s *ListService) RestoreList(ctx context.Context, req *pb.RestoreListRequest) (*pb.RestoreListResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	restoreReq := &repositories.RestoreListRequest{
		ListID:  req.ListID,
		BoardID: boardID,
	}

	if err := s.listRepo.RestoreList(restoreReq); err != nil {
		return nil, err
	}

	return &pb.RestoreListResponse{
		Message: "List restored successfully",
	}, nil
}

func (s *ListService) DeleteList(ctx context.Context, req *pb.DeleteListRequest) (*pb.DeleteListResponse, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	deleteReq := &repositories.DeleteListRequest{
		ID:      req.ListID,
		BoardID: boardID,
		UserID:  userID,
	}

	err := s.listRepo.DeleteList(deleteReq)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteListResponse{Message: "List deleted successfully"}, nil
}
