package middlewares

import (
	"context"

	pb_card "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/sm888sm/halten-backend/common/constants/fielderrors"
	"github.com/sm888sm/halten-backend/common/errorhandler"
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
		if err := validateCreateCardRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/GetCardByID":
		req := req.(*pb_card.GetCardByIDRequest)
		if err := validateGetCardByIDRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/GetCardsByList":
		req := req.(*pb_card.GetCardsByListRequest)
		if err := validateGetCardsByListRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/GetCardsByBoard":
		req := req.(*pb_card.GetCardsByBoardRequest)
		if err := validateGetCardsByBoardRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/MoveCardPosition":
		req := req.(*pb_card.MoveCardPositionRequest)
		if err := validateMoveCardPositionRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/UpdateCardName":
		req := req.(*pb_card.UpdateCardNameRequest)
		if err := validateUpdateCardNameRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/UpdateCardDescription":
		req := req.(*pb_card.UpdateCardDescriptionRequest)
		if err := validateUpdateCardDescriptionRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/AddCardLabel":
		req := req.(*pb_card.AddCardLabelRequest)
		if err := validateAddCardLabelRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RemoveCardLabel":
		req := req.(*pb_card.RemoveCardLabelRequest)
		if err := validateRemoveCardLabelRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/SetCardDates":
		req := req.(*pb_card.SetCardDatesRequest)
		if err := validateSetCardDatesRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/MarkCardComplete":
		req := req.(*pb_card.MarkCardCompleteRequest)
		if err := validateMarkCardCompleteRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/AddCardAttachment":
		req := req.(*pb_card.AddCardAttachmentRequest)
		if err := validateAddCardAttachmentRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RemoveCardAttachment":
		req := req.(*pb_card.RemoveCardAttachmentRequest)
		if err := validateRemoveCardAttachmentRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/AddCardComment":
		req := req.(*pb_card.AddCardCommentRequest)
		if err := validateAddCardCommentRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RemoveCardComment":
		req := req.(*pb_card.RemoveCardCommentRequest)
		if err := validateRemoveCardCommentRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/AddCardMembers":
		req := req.(*pb_card.AddCardMembersRequest)
		if err := validateAddCardMembersRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RemoveCardMembers":
		req := req.(*pb_card.RemoveCardMembersRequest)
		if err := validateRemoveCardMembersRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/ArchiveCard":
		req := req.(*pb_card.ArchiveCardRequest)
		if err := validateArchiveCardRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RestoreCard":
		req := req.(*pb_card.RestoreCardRequest)
		if err := validateRestoreCardRequest(ctx, req); err != nil {
			return nil, err
		}
	case "/proto.CardService/DeleteCard":
		req := req.(*pb_card.DeleteCardRequest)
		if err := validateDeleteCardRequest(ctx, req); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

func validateCreateCardRequest(ctx context.Context, req *pb_card.CreateCardRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.ListID == 0 {
		fieldErrors["ListID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "ListID is required",
			Field:   "ListID",
		}
	}
	if req.Name == "" {
		fieldErrors["Name"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Name is required",
			Field:   "Name",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateGetCardByIDRequest(ctx context.Context, req *pb_card.GetCardByIDRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)

}

func validateGetCardsByListRequest(ctx context.Context, req *pb_card.GetCardsByListRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.ListID == 0 {
		fieldErrors["ListID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "ListID is required",
			Field:   "ListID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)

}

func validateGetCardsByBoardRequest(ctx context.Context, req *pb_card.GetCardsByBoardRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)

}

func validateMoveCardPositionRequest(ctx context.Context, req *pb_card.MoveCardPositionRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.NewPosition == 0 {
		fieldErrors["Position"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Position is required",
			Field:   "NewPosition",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateUpdateCardNameRequest(ctx context.Context, req *pb_card.UpdateCardNameRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.Name == "" {
		fieldErrors["Name"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Name is required",
			Field:   "Name",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateUpdateCardDescriptionRequest(ctx context.Context, req *pb_card.UpdateCardDescriptionRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}
	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.Description == "" {
		fieldErrors["Description"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Description is required",
			Field:   "Description",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateAddCardLabelRequest(ctx context.Context, req *pb_card.AddCardLabelRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.LabelID == 0 {
		fieldErrors["LabelID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "LabelID is required",
			Field:   "LabelID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateRemoveCardLabelRequest(ctx context.Context, req *pb_card.RemoveCardLabelRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.LabelID == 0 {
		fieldErrors["LabelID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "LabelID is required",
			Field:   "LabelID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateSetCardDatesRequest(ctx context.Context, req *pb_card.SetCardDatesRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateMarkCardCompleteRequest(ctx context.Context, req *pb_card.MarkCardCompleteRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateAddCardAttachmentRequest(ctx context.Context, req *pb_card.AddCardAttachmentRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.AttachmentID == 0 {
		fieldErrors["Attachment"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Attachment is required",
			Field:   "Attachment",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateRemoveCardAttachmentRequest(ctx context.Context, req *pb_card.RemoveCardAttachmentRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.AttachmentID == 0 {
		fieldErrors["AttachmentID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "AttachmentID is required",
			Field:   "AttachmentID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateAddCardCommentRequest(ctx context.Context, req *pb_card.AddCardCommentRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.Content == "" {
		fieldErrors["Content"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Content is required",
			Field:   "Content",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateRemoveCardCommentRequest(ctx context.Context, req *pb_card.RemoveCardCommentRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.CommentID == 0 {
		fieldErrors["CommentID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CommentID is required",
			Field:   "CommentID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateAddCardMembersRequest(ctx context.Context, req *pb_card.AddCardMembersRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if len(req.UserIDs) == 0 {
		fieldErrors["UserIDs"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "At least one UserID is required",
			Field:   "UserIDs",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateRemoveCardMembersRequest(ctx context.Context, req *pb_card.RemoveCardMembersRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if len(req.UserIDs) == 0 {
		fieldErrors["UserIDs"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "At least one UserID is required",
			Field:   "UserIDs",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateArchiveCardRequest(ctx context.Context, req *pb_card.ArchiveCardRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateRestoreCardRequest(ctx context.Context, req *pb_card.RestoreCardRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateDeleteCardRequest(ctx context.Context, req *pb_card.DeleteCardRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}
