package middlewares

import (
	"context"

	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/sm888sm/halten-backend/common/constants/fielderrors"
	"github.com/sm888sm/halten-backend/common/errorhandlers"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type ValidatorInterceptor struct {
	db *gorm.DB
}

func NewValidatorInterceptor(db *gorm.DB) *ValidatorInterceptor {
	return &ValidatorInterceptor{db: db}
}

func (v *ValidatorInterceptor) ValidationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	switch info.FullMethod {
	// Card Service
	case "/proto.CardService/CreateCard":
		req := req.(*pb_card.CreateCardRequest)
		if err := validateCreateCardRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/GetCardByID":
		req := req.(*pb_card.GetCardByIDRequest)
		if err := validateGetCardByIDRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/GetCardsByList":
		req := req.(*pb_card.GetCardsByListRequest)
		if err := validateGetCardsByListRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/GetCardsByBoard":
		req := req.(*pb_card.GetCardsByBoardRequest)
		if err := validateGetCardsByBoardRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/MoveCardPosition":
		req := req.(*pb_card.MoveCardPositionRequest)
		if err := validateMoveCardPositionRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/UpdateCardName":
		req := req.(*pb_card.UpdateCardNameRequest)
		if err := validateUpdateCardNameRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/UpdateCardDescription":
		req := req.(*pb_card.UpdateCardDescriptionRequest)
		if err := validateUpdateCardDescriptionRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/AddCardLabel":
		req := req.(*pb_card.AddCardLabelRequest)
		if err := validateAddCardLabelRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RemoveCardLabel":
		req := req.(*pb_card.RemoveCardLabelRequest)
		if err := validateRemoveCardLabelRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/SetCardDates":
		req := req.(*pb_card.SetCardDatesRequest)
		if err := validateSetCardDatesRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/MarkCardComplete":
		req := req.(*pb_card.MarkCardCompleteRequest)
		if err := validateMarkCardCompleteRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/AddCardAttachment":
		req := req.(*pb_card.AddCardAttachmentRequest)
		if err := validateAddCardAttachmentRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RemoveCardAttachment":
		req := req.(*pb_card.RemoveCardAttachmentRequest)
		if err := validateRemoveCardAttachmentRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/AddCardComment":
		req := req.(*pb_card.AddCardCommentRequest)
		if err := validateAddCardCommentRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RemoveCardComment":
		req := req.(*pb_card.RemoveCardCommentRequest)
		if err := validateRemoveCardCommentRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/AddCardMembers":
		req := req.(*pb_card.AddCardMembersRequest)
		if err := validateAddCardMembersRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RemoveCardMembers":
		req := req.(*pb_card.RemoveCardMembersRequest)
		if err := validateRemoveCardMembersRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/ArchiveCard":
		req := req.(*pb_card.ArchiveCardRequest)
		if err := validateArchiveCardRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RestoreCard":
		req := req.(*pb_card.RestoreCardRequest)
		if err := validateRestoreCardRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/DeleteCard":
		req := req.(*pb_card.DeleteCardRequest)
		if err := validateDeleteCardRequest(req); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

func validateCreateCardRequest(req *pb_card.CreateCardRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.ListID == 0 {
		fieldErrors["ListID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "ListID is required",
			Field:   "ListID",
		}
	}
	if req.Name == "" {
		fieldErrors["Name"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Name is required",
			Field:   "Name",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateGetCardByIDRequest(req *pb_card.GetCardByIDRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)

}

func validateGetCardsByListRequest(req *pb_card.GetCardsByListRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.ListID == 0 {
		fieldErrors["ListID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "ListID is required",
			Field:   "ListID",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)

}

func validateGetCardsByBoardRequest(_ *pb_card.GetCardsByBoardRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)

}

func validateMoveCardPositionRequest(req *pb_card.MoveCardPositionRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.Position == 0 {
		fieldErrors["Position"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Position is required",
			Field:   "Position",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateUpdateCardNameRequest(req *pb_card.UpdateCardNameRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.Name == "" {
		fieldErrors["Name"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Name is required",
			Field:   "Name",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateUpdateCardDescriptionRequest(req *pb_card.UpdateCardDescriptionRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)
	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.Description == "" {
		fieldErrors["Description"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Description is required",
			Field:   "Description",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateAddCardLabelRequest(req *pb_card.AddCardLabelRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.LabelID == 0 {
		fieldErrors["LabelID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "LabelID is required",
			Field:   "LabelID",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateRemoveCardLabelRequest(req *pb_card.RemoveCardLabelRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.LabelID == 0 {
		fieldErrors["LabelID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "LabelID is required",
			Field:   "LabelID",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateSetCardDatesRequest(req *pb_card.SetCardDatesRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateMarkCardCompleteRequest(req *pb_card.MarkCardCompleteRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateAddCardAttachmentRequest(req *pb_card.AddCardAttachmentRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.AttachmentID == 0 {
		fieldErrors["Attachment"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Attachment is required",
			Field:   "Attachment",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateRemoveCardAttachmentRequest(req *pb_card.RemoveCardAttachmentRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.AttachmentID == 0 {
		fieldErrors["AttachmentID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "AttachmentID is required",
			Field:   "AttachmentID",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateAddCardCommentRequest(req *pb_card.AddCardCommentRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.Content == "" {
		fieldErrors["Content"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Content is required",
			Field:   "Content",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateRemoveCardCommentRequest(req *pb_card.RemoveCardCommentRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.CommentID == 0 {
		fieldErrors["CommentID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CommentID is required",
			Field:   "CommentID",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateAddCardMembersRequest(req *pb_card.AddCardMembersRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if len(req.UserIDs) == 0 {
		fieldErrors["UserIDs"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "At least one UserID is required",
			Field:   "UserIDs",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateRemoveCardMembersRequest(req *pb_card.RemoveCardMembersRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if len(req.UserIDs) == 0 {
		fieldErrors["UserIDs"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "At least one UserID is required",
			Field:   "UserIDs",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateArchiveCardRequest(req *pb_card.ArchiveCardRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateRestoreCardRequest(req *pb_card.RestoreCardRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateDeleteCardRequest(req *pb_card.DeleteCardRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}
