package services

import (
	"context"

	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/sm888sm/halten-backend/card-service/internal/repositories"
	"github.com/sm888sm/halten-backend/common/constants/contextkeys"
	"github.com/sm888sm/halten-backend/common/errorhandlers"
	"github.com/sm888sm/halten-backend/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CardService struct {
	cardRepo repositories.CardRepository
	pb_card.UnimplementedCardServiceServer
}

func NewCardService(repo repositories.CardRepository) *CardService {
	return &CardService{cardRepo: repo}
}

func (s *CardService) CreateCard(ctx context.Context, req *pb_card.CreateCardRequest) (*pb_card.CreateCardResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	card := &models.Card{
		Name:    req.Name,
		BoardID: boardID,
		ListID:  req.ListID,
	}

	repoRes, err := s.cardRepo.CreateCard(&repositories.CreateCardRequest{Card: card})
	if err != nil {
		return nil, err
	}
	return &pb_card.CreateCardResponse{
		Card: &pb_card.Card{
			CardID:   repoRes.Card.ID,
			Name:     repoRes.Card.Name,
			ListID:   repoRes.Card.ListID,
			Position: repoRes.Card.Position,
		},
	}, nil
}

func (s *CardService) GetCardByID(ctx context.Context, req *pb_card.GetCardByIDRequest) (*pb_card.GetCardByIDResponse, error) {
	repoRes, err := s.cardRepo.GetCardByID(&repositories.GetCardByIDRequest{
		CardID: req.CardID,
	})
	if err != nil {
		return nil, err
	}

	card := repoRes.Card

	return &pb_card.GetCardByIDResponse{
		Card: &pb_card.Card{
			CardID:      uint64(card.ID),
			ListID:      uint64(card.ListID),
			Name:        card.Name,
			Position:    int64(card.Position),
			StartDate:   timestamppb.New(*card.StartDate),
			DueDate:     timestamppb.New(*card.DueDate),
			Attachments: card.Attachments,
			Labels:      card.Labels,
			Members:     card.Members,
			CreatedAt:   timestamppb.New(card.CreatedAt),
			UpdatedAt:   timestamppb.New(card.UpdatedAt),
		},
	}, nil
}

func (s *CardService) GetCardsByList(ctx context.Context, req *pb_card.GetCardsByListRequest) (*pb_card.GetCardsByListResponse, error) {
	repoRes, err := s.cardRepo.GetCardsByList(&repositories.GetCardsByListRequest{
		ListID: req.ListID,
	})
	if err != nil {
		return nil, err
	}

	var pb_cardCards []*pb_card.CardMeta
	for _, c := range repoRes.Cards {
		pb_cardCards = append(pb_cardCards, &pb_card.CardMeta{
			CardID:          c.ID,
			ListID:          c.ListID,
			BoardID:         c.BoardID,
			Name:            c.Name,
			Position:        c.Position,
			StartDate:       timestamppb.New(*c.StartDate),
			DueDate:         timestamppb.New(*c.DueDate),
			Labels:          c.Labels,
			Members:         c.Members,
			TotalAttachment: c.TotalAttachment,
			TotalComment:    c.TotalComment,
			CreatedAt:       timestamppb.New(c.CreatedAt),
			UpdatedAt:       timestamppb.New(c.UpdatedAt),
		})
	}

	return &pb_card.GetCardsByListResponse{
		Cards: pb_cardCards,
	}, nil
}

func (s *CardService) GetCardsByBoard(ctx context.Context, req *pb_card.GetCardsByBoardRequest) (*pb_card.GetCardsByBoardResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoRes, err := s.cardRepo.GetCardsByBoard(&repositories.GetCardsByBoardRequest{
		BoardID: boardID,
	})

	if err != nil {
		return nil, err
	}

	var pb_cardCards []*pb_card.CardMeta
	for _, c := range repoRes.Cards {
		pb_cardCards = append(pb_cardCards, &pb_card.CardMeta{
			CardID:          c.ID,
			ListID:          c.ListID,
			BoardID:         c.BoardID,
			Name:            c.Name,
			Position:        c.Position,
			StartDate:       timestamppb.New(*c.StartDate),
			DueDate:         timestamppb.New(*c.DueDate),
			Labels:          c.Labels,
			Members:         c.Members,
			TotalAttachment: c.TotalAttachment,
			TotalComment:    c.TotalComment,
			CreatedAt:       timestamppb.New(c.CreatedAt),
			UpdatedAt:       timestamppb.New(c.UpdatedAt),
		})
	}

	return &pb_card.GetCardsByBoardResponse{
		Cards: pb_cardCards,
	}, nil
}

func (s *CardService) MoveCardPosition(ctx context.Context, req *pb_card.MoveCardPositionRequest) (*pb_card.MoveCardPositionResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.MoveCardPositionRequest{
		CardID:    req.CardID,
		Position:  req.Position,
		BoardID:   boardID,
		OldListID: req.OldListID,
		NewListID: req.NewListID,
	}

	err := s.cardRepo.MoveCardPosition(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.MoveCardPositionResponse{}, nil
}

func (s *CardService) UpdateCardName(ctx context.Context, req *pb_card.UpdateCardNameRequest) (*pb_card.UpdateCardNameResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.UpdateCardNameRequest{
		CardID:  req.CardID,
		Name:    req.Name,
		BoardID: boardID,
	}
	err := s.cardRepo.UpdateCardName(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.UpdateCardNameResponse{}, nil
}

func (s *CardService) UpdateCardDescription(ctx context.Context, req *pb_card.UpdateCardDescriptionRequest) (*pb_card.UpdateCardDescriptionResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.UpdateCardDescriptionRequest{
		CardID:      req.CardID,
		Description: req.Description,
		BoardID:     boardID,
	}
	err := s.cardRepo.UpdateCardDescription(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.UpdateCardDescriptionResponse{
		Message: "Card description updated",
	}, nil
}

func (s *CardService) AddCardLabel(ctx context.Context, req *pb_card.AddCardLabelRequest) (*pb_card.AddCardLabelResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.AddCardLabelRequest{
		LabelID: req.LabelID,
		CardID:  req.CardID,
		BoardID: boardID,
	}
	err := s.cardRepo.AddCardLabel(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.AddCardLabelResponse{
		Message: "Label added to card",
	}, nil
}

func (s *CardService) RemoveCardLabel(ctx context.Context, req *pb_card.RemoveCardLabelRequest) (*pb_card.RemoveCardLabelResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.RemoveCardLabelRequest{
		LabelID: req.LabelID,
		CardID:  req.CardID,
		BoardID: boardID,
	}
	err := s.cardRepo.RemoveCardLabel(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.RemoveCardLabelResponse{
		Message: "Label removed from card",
	}, nil
}

func (s *CardService) SetCardDates(ctx context.Context, req *pb_card.SetCardDatesRequest) (*pb_card.SetCardDatesResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	startDate := req.StartDate.AsTime()
	dueDate := req.DueDate.AsTime()

	repoReq := &repositories.SetCardDatesRequest{
		StartDate: &startDate,
		DueDate:   &dueDate,
		CardID:    req.CardID,
		BoardID:   boardID,
	}
	err := s.cardRepo.SetCardDates(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.SetCardDatesResponse{
		Message: "Card dates updated",
	}, nil
}

func (s *CardService) ToggleCardCompleted(ctx context.Context, req *pb_card.ToggleCardCompletedRequest) (*pb_card.ToggleCardCompletedResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.ToggleCardCompletedRequest{
		CardID:  req.CardID,
		BoardID: boardID,
	}

	err := s.cardRepo.ToggleCardCompleted(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.ToggleCardCompletedResponse{
		Message: "Card completion status toggled",
	}, nil
}

func (s *CardService) AddCardAttachment(ctx context.Context, req *pb_card.AddCardAttachmentRequest) (*pb_card.AddCardAttachmentResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.AddCardAttachmentRequest{
		AttachmentID: req.AttachmentID,
		CardID:       req.CardID,
		BoardID:      boardID,
	}

	err := s.cardRepo.AddCardAttachment(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.AddCardAttachmentResponse{
		Message: "Attachment added to card",
	}, nil
}

func (s *CardService) RemoveCardAttachment(ctx context.Context, req *pb_card.RemoveCardAttachmentRequest) (*pb_card.RemoveCardAttachmentResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.RemoveCardAttachmentRequest{
		AttachmentID: req.AttachmentID,
		CardID:       req.CardID,
		BoardID:      boardID,
	}

	err := s.cardRepo.RemoveCardAttachment(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.RemoveCardAttachmentResponse{
		Message: "Attachment removed from card",
	}, nil
}

func (s *CardService) AddCardComment(ctx context.Context, req *pb_card.AddCardCommentRequest) (*pb_card.AddCardCommentResponse, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		// Handle error: userID was not a uint64
		return nil, errorhandlers.NewGrpcInternalError()
	}

	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	comment := models.Comment{
		Content: req.Content,
	}

	repoReq := &repositories.AddCardCommentRequest{
		Comment: comment,
		CardID:  req.CardID,
		BoardID: boardID,
		UserID:  userID,
	}

	err := s.cardRepo.AddCardComment(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.AddCardCommentResponse{
		Message: "Comment added to card",
	}, nil
}

func (s *CardService) RemoveCardComment(ctx context.Context, req *pb_card.RemoveCardCommentRequest) (*pb_card.RemoveCardCommentResponse, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.RemoveCardCommentRequest{
		CommentID: req.CommentID,
		CardID:    req.CardID,
		BoardID:   boardID,
		UserID:    userID,
	}

	err := s.cardRepo.RemoveCardComment(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.RemoveCardCommentResponse{
		Message: "Comment removed from card",
	}, nil
}

func (s *CardService) AddCardMembers(ctx context.Context, req *pb_card.AddCardMembersRequest) (*pb_card.AddCardMembersResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	userIDs := append([]uint64(nil), req.UserIDs...)

	repoReq := &repositories.AddCardMembersRequest{
		UserIDs: userIDs,
		CardID:  req.CardID,
		BoardID: boardID,
	}
	err := s.cardRepo.AddCardMembers(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.AddCardMembersResponse{
		Message: "Members added to card",
	}, nil
}

func (s *CardService) RemoveCardMembers(ctx context.Context, req *pb_card.RemoveCardMembersRequest) (*pb_card.RemoveCardMembersResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.RemoveCardMembersRequest{
		UserIDs: req.UserIDs,
		CardID:  req.CardID,
		BoardID: boardID,
	}

	err := s.cardRepo.RemoveCardMembers(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.RemoveCardMembersResponse{
		Message: "Members removed from card",
	}, nil
}

func (s *CardService) ArchiveCard(ctx context.Context, req *pb_card.ArchiveCardRequest) (*pb_card.ArchiveCardResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.ArchiveCardRequest{
		CardID:  req.CardID,
		BoardID: boardID,
	}

	err := s.cardRepo.ArchiveCard(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.ArchiveCardResponse{
		Message: "Card archived",
	}, nil
}

func (s *CardService) RestoreCard(ctx context.Context, req *pb_card.RestoreCardRequest) (*pb_card.RestoreCardResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.RestoreCardRequest{
		CardID:  req.CardID,
		BoardID: boardID,
	}

	err := s.cardRepo.RestoreCard(repoReq)
	if err != nil {
		return nil, err
	}

	return &pb_card.RestoreCardResponse{
		Message: "Card restored",
	}, nil
}

func (s *CardService) DeleteCard(ctx context.Context, req *pb_card.DeleteCardRequest) (*pb_card.DeleteCardResponse, error) {
	boardID, ok := ctx.Value(contextkeys.BoardIDKey{}).(uint64)
	if !ok {
		return nil, errorhandlers.NewGrpcInternalError()
	}

	repoReq := &repositories.DeleteCardRequest{
		CardID:  req.CardID,
		BoardID: boardID,
	}

	err := s.cardRepo.DeleteCard(repoReq)
	if err != nil {
		return nil, err
	}
	return &pb_card.DeleteCardResponse{
		Message: "Card deleted",
	}, nil
}
