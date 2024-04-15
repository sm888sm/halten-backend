package services

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/sm888sm/halten-backend/card-service/internal/repositories"
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
	card := &models.Card{
		Name:    req.Name,
		BoardID: uint(req.BoardId),
		ListID:  uint(req.ListId),
	}
	err := s.cardRepo.CreateCard(card, uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.CreateCardResponse{
		Id: uint64(card.ID),
	}, nil
}

func (s *CardService) GetCardByID(ctx context.Context, req *pb.GetCardByIDRequest) (*pb.GetCardByIDResponse, error) {
	card, err := s.cardRepo.GetCardByID(uint(req.Id), uint(req.BoardId), uint(req.UserId))
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
			Id:          uint64(card.ID),
			ListId:      uint64(card.ListID),
			Name:        card.Name,
			Position:    int64(card.Position),
			StartDate:   timestamppb.New(*card.StartDate),
			DueDate:     timestamppb.New(*card.DueDate),
			Attachments: attachments,
			Labels:      labels,
			Members:     members,
		},
	}, nil
}

func (s *CardService) GetCardsByList(ctx context.Context, req *pb.GetCardsByListRequest) (*pb.GetCardsByListResponse, error) {
	cards, err := s.cardRepo.GetCardsByList(uint(req.ListId), uint(req.BoardId), uint(req.UserId))
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
			Id:              uint64(c.ID),
			ListId:          uint64(c.ListID),
			BoardId:         uint64(c.BoardID),
			Name:            c.Name,
			Position:        int32(c.Position),
			StartDate:       timestamppb.New(*c.StartDate),
			DueDate:         timestamppb.New(*c.DueDate),
			Labels:          labels,
			Members:         members,
			TotalAttachment: c.TotalAttachment,
			TotalComment:    c.TotalComment,
		})
	}

	return &pb.GetCardsByListResponse{
		Cards: pbCards,
	}, nil
}

func (s *CardService) GetCardsByBoard(ctx context.Context, req *pb.GetCardsByBoardRequest) (*pb.GetCardsByBoardResponse, error) {
	cards, err := s.cardRepo.GetCardsByBoard(uint(req.BoardId), uint(req.UserId))
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
			Id:              uint64(c.ID),
			ListId:          uint64(c.ListID),
			BoardId:         uint64(c.BoardID),
			Name:            c.Name,
			Position:        int32(c.Position),
			StartDate:       timestamppb.New(*c.StartDate),
			DueDate:         timestamppb.New(*c.DueDate),
			Labels:          labels,
			Members:         members,
			TotalAttachment: c.TotalAttachment,
			TotalComment:    c.TotalComment,
		})
	}

	return &pb.GetCardsByListResponse{
		Cards: pbCards,
	}, nil
}

func (s *CardService) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
	err := s.cardRepo.DeleteCard(uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCardResponse{}, nil
}

func (s *CardService) MoveCardPosition(ctx context.Context, req *pb.MoveCardPositionRequest) (*pb.MoveCardPositionResponse, error) {
	err := s.cardRepo.MoveCardPosition(uint(req.Id), int(req.NewPosition), uint(req.BoardId), uint(req.OldListId), uint(req.NewListId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.MoveCardPositionResponse{}, nil
}

func (s *CardService) UpdateCardName(ctx context.Context, req *pb.UpdateCardNameRequest) (*pb.UpdateCardNameResponse, error) {
	err := s.cardRepo.UpdateCardName(uint(req.Id), req.NewName, uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.UpdateCardNameResponse{}, nil
}

func (s *CardService) UpdateCardDescription(ctx context.Context, req *pb.UpdateCardDescriptionRequest) (*pb.UpdateCardDescriptionResponse, error) {
	err := s.cardRepo.UpdateCardDescription(uint(req.Id), req.NewDescription, uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.UpdateCardDescriptionResponse{}, nil
}

func (s *CardService) AddCardLabel(ctx context.Context, req *pb.AddCardLabelRequest) (*pb.AddCardLabelResponse, error) {
	err := s.cardRepo.AddCardLabel(req.Label, uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.AddCardLabelResponse{}, nil
}

func (s *CardService) RemoveCardLabel(ctx context.Context, req *pb.RemoveCardLabelRequest) (*pb.RemoveCardLabelResponse, error) {
	err := s.cardRepo.RemoveCardLabel(uint(req.LabelId), uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.RemoveCardLabelResponse{}, nil
}

func (s *CardService) SetCardDates(ctx context.Context, req *pb.SetCardDatesRequest) (*pb.SetCardDatesResponse, error) {
	startDate, _ := ptypes.Timestamp(req.StartDate)
	dueDate, _ := ptypes.Timestamp(req.DueDate)
	err := s.cardRepo.SetCardDates(&startDate, &dueDate, uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.SetCardDatesResponse{}, nil
}

func (s *CardService) MarkCardComplete(ctx context.Context, req *pb.MarkCardCompleteRequest) (*pb.MarkCardCompleteResponse, error) {
	err := s.cardRepo.MarkCardComplete(uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.MarkCardCompleteResponse{}, nil
}

func (s *CardService) AddCardAttachment(ctx context.Context, req *pb.AddCardAttachmentRequest) (*pb.AddCardAttachmentResponse, error) {
	err := s.cardRepo.AddCardAttachment(uint(req.AttachmentId), uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.AddCardAttachmentResponse{}, nil
}

func (s *CardService) RemoveCardAttachment(ctx context.Context, req *pb.RemoveCardAttachmentRequest) (*pb.RemoveCardAttachmentResponse, error) {
	err := s.cardRepo.RemoveCardAttachment(uint(req.AttachmentId), uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.RemoveCardAttachmentResponse{}, nil
}

func (s *CardService) AddCardComment(ctx context.Context, req *pb.AddCardCommentRequest) (*pb.AddCardCommentResponse, error) {
	comment := models.Comment{
		Content: req.Content,
	}
	err := s.cardRepo.AddCardComment(comment, uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.AddCardCommentResponse{}, nil
}

func (s *CardService) RemoveCardComment(ctx context.Context, req *pb.RemoveCardCommentRequest) (*pb.RemoveCardCommentResponse, error) {
	err := s.cardRepo.RemoveCardComment(uint(req.CommentId), uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.RemoveCardCommentResponse{}, nil
}

func (s *CardService) AddCardMembers(ctx context.Context, req *pb.AddCardMembersRequest) (*pb.AddCardMembersResponse, error) {
	err := s.cardRepo.AddCardMembers(req.UserIds, uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.AddCardMembersResponse{}, nil
}

func (s *CardService) RemoveCardMembers(ctx context.Context, req *pb.RemoveCardMembersRequest) (*pb.RemoveCardMembersResponse, error) {
	err := s.cardRepo.RemoveCardMembers(req.UserIds, uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.RemoveCardMembersResponse{}, nil
}

func (s *CardService) ArchiveCard(ctx context.Context, req *pb.ArchiveCardRequest) (*pb.ArchiveCardResponse, error) {
	err := s.cardRepo.ArchiveCard(uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.ArchiveCardResponse{}, nil
}

func (s *CardService) RestoreCard(ctx context.Context, req *pb.RestoreCardRequest) (*pb.RestoreCardResponse, error) {
	err := s.cardRepo.RestoreCard(uint(req.Id), uint(req.BoardId), uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.RestoreCardResponse{}, nil
}
