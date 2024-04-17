package services

import (
	"context"

	pb "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/sm888sm/halten-backend/card-service/internal/repositories"
	"github.com/sm888sm/halten-backend/common/constants/contextkeys"
	"github.com/sm888sm/halten-backend/common/errorhandler"
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
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	card := &models.Card{
		Name:    req.Name,
		BoardID: uint(boardID),
		ListID:  uint(req.ListID),
	}

	err := s.cardRepo.CreateCard(repositories.CreateCardParams{Card: card})
	if err != nil {
		return nil, err
	}
	return &pb.CreateCardResponse{
		CardID: uint64(card.ID),
		Name:   card.Name,
		ListID: uint64(card.ListID),
	}, nil
}

func (s *CardService) GetCardByID(ctx context.Context, req *pb.GetCardByIDRequest) (*pb.GetCardByIDResponse, error) {
	card, err := s.cardRepo.GetCardByID(repositories.GetCardByIDParams{
		CardID: req.CardID,
	})
	if err != nil {
		return nil, err
	}

	// Convert slices of Attachment, Label, and User to slices of uint64 for protobuf
	var attachments, labels, members []uint64
	for _, a := range card.Attachments {
		attachments = append(attachments, uint64(a.ID))
	}
	for _, l := range card.Labels {
		labels = append(labels, uint64(l.ID))
	}
	for _, m := range card.Members {
		members = append(members, uint64(m.ID))
	}

	return &pb.GetCardByIDResponse{
		Card: &pb.Card{
			CardID:      uint64(card.ID),
			ListID:      uint64(card.ListID),
			Name:        card.Name,
			Position:    int64(card.Position),
			StartDate:   timestamppb.New(*card.StartDate),
			DueDate:     timestamppb.New(*card.DueDate),
			Attachments: attachments,
			Labels:      labels,
			Members:     members,
			CreatedAt:   timestamppb.New(card.CreatedAt),
			UpdatedAt:   timestamppb.New(card.UpdatedAt),
		},
	}, nil
}

func (s *CardService) GetCardsByList(ctx context.Context, req *pb.GetCardsByListRequest) (*pb.GetCardsByListResponse, error) {
	cards, err := s.cardRepo.GetCardsByList(repositories.GetCardsByListParams{
		ListID: req.ListID,
	})
	if err != nil {
		return nil, err
	}

	var pbCards []*pb.CardMeta
	for _, c := range cards {
		// Convert slice of Label to slice of uint64 for protobuf
		var labels []uint64
		for _, l := range c.Labels {
			labels = append(labels, uint64(l.ID))
		}

		var members []uint64
		for _, l := range c.Members {
			members = append(members, uint64(l.ID))
		}

		pbCards = append(pbCards, &pb.CardMeta{
			CardID:          c.ID,
			ListID:          c.ListID,
			BoardID:         c.BoardID,
			Name:            c.Name,
			Position:        int64(c.Position),
			StartDate:       timestamppb.New(*c.StartDate),
			DueDate:         timestamppb.New(*c.DueDate),
			Labels:          labels,
			Members:         members,
			TotalAttachment: c.TotalAttachment,
			TotalComment:    c.TotalComment,
			CreatedAt:       timestamppb.New(*c.CreatedAt),
			UpdatedAt:       timestamppb.New(*c.UpdatedAt),
		})
	}

	return &pb.GetCardsByListResponse{
		Cards: pbCards,
	}, nil
}

func (s *CardService) GetCardsByBoard(ctx context.Context, req *pb.GetCardsByBoardRequest) (*pb.GetCardsByBoardResponse, error) {

	cards, err := s.cardRepo.GetCardsByBoard(repositories.GetCardsByBoardParams{
		BoardID: req.BoardID,
	})

	if err != nil {
		return nil, err
	}

	var pbCards []*pb.CardMeta
	for _, c := range cards {
		// Convert slice of Label to slice of uint64 for protobuf
		var labels []uint64
		for _, l := range c.Labels {
			labels = append(labels, uint64(l.ID))
		}

		var members []uint64
		for _, l := range c.Members {
			members = append(members, uint64(l.ID))
		}

		pbCards = append(pbCards, &pb.CardMeta{
			CardID:          c.ID,
			ListID:          c.ListID,
			BoardID:         c.BoardID,
			Name:            c.Name,
			Position:        int64(c.Position),
			StartDate:       timestamppb.New(*c.StartDate),
			DueDate:         timestamppb.New(*c.DueDate),
			Labels:          labels,
			Members:         members,
			TotalAttachment: c.TotalAttachment,
			TotalComment:    c.TotalComment,
		})
	}

	return &pb.GetCardsByBoardResponse{
		Cards: pbCards,
	}, nil
}

func (s *CardService) MoveCardPosition(ctx context.Context, req *pb.MoveCardPositionRequest) (*pb.MoveCardPositionResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	params := repositories.MoveCardPositionParams{
		CardID:      req.CardID,
		NewPosition: int(req.NewPosition),
		BoardID:     boardID,
		OldListID:   req.OldListID,
		NewListID:   req.NewListID,
	}
	err := s.cardRepo.MoveCardPosition(params)
	if err != nil {
		return nil, err
	}
	return &pb.MoveCardPositionResponse{}, nil
}

func (s *CardService) UpdateCardName(ctx context.Context, req *pb.UpdateCardNameRequest) (*pb.UpdateCardNameResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	params := repositories.UpdateCardNameParams{
		CardID:  req.CardID,
		Name:    req.Name,
		BoardID: boardID,
	}
	err := s.cardRepo.UpdateCardName(params)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateCardNameResponse{}, nil
}

func (s *CardService) UpdateCardDescription(ctx context.Context, req *pb.UpdateCardDescriptionRequest) (*pb.UpdateCardDescriptionResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	params := repositories.UpdateCardDescriptionParams{
		CardID:         req.CardID,
		NewDescription: req.Description,
		BoardID:        boardID,
	}
	err := s.cardRepo.UpdateCardDescription(params)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateCardDescriptionResponse{}, nil
}

func (s *CardService) AddCardLabel(ctx context.Context, req *pb.AddCardLabelRequest) (*pb.AddCardLabelResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	params := repositories.AddCardLabelParams{
		LabelID: req.LabelID,
		CardID:  req.CardID,
		BoardID: boardID,
	}
	err := s.cardRepo.AddCardLabel(params)
	if err != nil {
		return nil, err
	}
	return &pb.AddCardLabelResponse{}, nil
}

func (s *CardService) RemoveCardLabel(ctx context.Context, req *pb.RemoveCardLabelRequest) (*pb.RemoveCardLabelResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	params := repositories.RemoveCardLabelParams{
		LabelID: req.LabelID,
		CardID:  req.CardID,
		BoardID: boardID,
	}
	err := s.cardRepo.RemoveCardLabel(params)
	if err != nil {
		return nil, err
	}
	return &pb.RemoveCardLabelResponse{}, nil
}

func (s *CardService) SetCardDates(ctx context.Context, req *pb.SetCardDatesRequest) (*pb.SetCardDatesResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	startDate := req.StartDate.AsTime()
	dueDate := req.DueDate.AsTime()

	params := repositories.SetCardDatesParams{
		StartDate: &startDate,
		DueDate:   &dueDate,
		CardID:    req.CardID,
		BoardID:   boardID,
	}
	err := s.cardRepo.SetCardDates(params)
	if err != nil {
		return nil, err
	}
	return &pb.SetCardDatesResponse{}, nil
}

func (s *CardService) MarkCardComplete(ctx context.Context, req *pb.MarkCardCompleteRequest) (*pb.MarkCardCompleteResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	params := repositories.MarkCardCompleteParams{
		CardID:  req.CardID,
		BoardID: boardID,
	}

	err := s.cardRepo.MarkCardComplete(params)
	if err != nil {
		return nil, err
	}
	return &pb.MarkCardCompleteResponse{}, nil
}

func (s *CardService) AddCardAttachment(ctx context.Context, req *pb.AddCardAttachmentRequest) (*pb.AddCardAttachmentResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	params := repositories.AddCardAttachmentParams{
		AttachmentID: req.AttachmentID,
		CardID:       req.CardID,
		BoardID:      boardID,
	}

	err := s.cardRepo.AddCardAttachment(params)
	if err != nil {
		return nil, err
	}
	return &pb.AddCardAttachmentResponse{}, nil
}

func (s *CardService) RemoveCardAttachment(ctx context.Context, req *pb.RemoveCardAttachmentRequest) (*pb.RemoveCardAttachmentResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	params := repositories.RemoveCardAttachmentParams{
		AttachmentID: req.AttachmentID,
		CardID:       req.CardID,
		BoardID:      boardID,
	}

	err := s.cardRepo.RemoveCardAttachment(params)
	if err != nil {
		return nil, err
	}
	return &pb.RemoveCardAttachmentResponse{}, nil
}

func (s *CardService) AddCardComment(ctx context.Context, req *pb.AddCardCommentRequest) (*pb.AddCardCommentResponse, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		// Handle error: userID was not a uint64
		return nil, errorhandler.NewGrpcInternalError()
	}

	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	comment := models.Comment{
		Content: req.Content,
	}

	params := repositories.AddCardCommentParams{
		Comment: comment,
		CardID:  req.CardID,
		BoardID: boardID,
		UserID:  userID,
	}

	err := s.cardRepo.AddCardComment(params)
	if err != nil {
		return nil, err
	}
	return &pb.AddCardCommentResponse{}, nil
}

func (s *CardService) RemoveCardComment(ctx context.Context, req *pb.RemoveCardCommentRequest) (*pb.RemoveCardCommentResponse, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	params := repositories.RemoveCardCommentParams{
		CommentID: req.CommentID,
		CardID:    req.CardID,
		BoardID:   boardID,
		UserID:    userID,
	}

	err := s.cardRepo.RemoveCardComment(params)
	if err != nil {
		return nil, err
	}
	return &pb.RemoveCardCommentResponse{}, nil
}

func (s *CardService) AddCardMembers(ctx context.Context, req *pb.AddCardMembersRequest) (*pb.AddCardMembersResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	userIDs := append([]uint64(nil), req.UserIDs...)

	params := repositories.AddCardMembersParams{
		UserIDs: userIDs,
		CardID:  req.CardID,
		BoardID: boardID,
	}
	err := s.cardRepo.AddCardMembers(params)
	if err != nil {
		return nil, err
	}
	return &pb.AddCardMembersResponse{}, nil
}

func (s *CardService) RemoveCardMembers(ctx context.Context, req *pb.RemoveCardMembersRequest) (*pb.RemoveCardMembersResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	params := repositories.RemoveCardMembersParams{
		UserIDs: req.UserIDs,
		CardID:  req.CardID,
		BoardID: boardID,
	}

	err := s.cardRepo.RemoveCardMembers(params)
	if err != nil {
		return nil, err
	}
	return &pb.RemoveCardMembersResponse{}, nil
}

func (s *CardService) ArchiveCard(ctx context.Context, req *pb.ArchiveCardRequest) (*pb.ArchiveCardResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	params := repositories.ArchiveCardParams{
		CardID:  req.CardID,
		BoardID: boardID,
	}

	err := s.cardRepo.ArchiveCard(params)
	if err != nil {
		return nil, err
	}
	return &pb.ArchiveCardResponse{}, nil
}

func (s *CardService) RestoreCard(ctx context.Context, req *pb.RestoreCardRequest) (*pb.RestoreCardResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	params := repositories.RestoreCardParams{
		CardID:  req.CardID,
		BoardID: boardID,
	}

	err := s.cardRepo.RestoreCard(params)
	if err != nil {
		return nil, err
	}

	return &pb.RestoreCardResponse{}, nil
}

func (s *CardService) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandler.NewGrpcInternalError()
	}

	err := s.cardRepo.DeleteCard(repositories.DeleteCardParams{
		CardID:  req.CardID,
		BoardID: boardID,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCardResponse{}, nil
}
