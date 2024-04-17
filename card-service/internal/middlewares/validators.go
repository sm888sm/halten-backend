package middlewares

import (
	"context"
	"errors"

	pb "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/sm888sm/halten-backend/common/constants/contextkeys"
	"github.com/sm888sm/halten-backend/common/constants/fielderrors"
	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/constants/roles"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	"github.com/sm888sm/halten-backend/common/helpers"
	"github.com/sm888sm/halten-backend/models"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	checkVisibility = map[string]bool{
		"/proto.CardService/GetCardByID":     true,
		"/proto.CardService/GetCardsByList":  true,
		"/proto.CardService/GetCardsByBoard": true,
	}

	checkRoles = map[string]string{
		"/proto.CardService/CreateCard":            roles.MemberRole,
		"/proto.CardService/MoveCardPosition":      roles.MemberRole,
		"/proto.CardService/UpdateCardName":        roles.MemberRole,
		"/proto.CardService/UpdateCardDescription": roles.MemberRole,
		"/proto.CardService/AddCardLabel":          roles.MemberRole,
		"/proto.CardService/RemoveCardLabel":       roles.MemberRole,
		"/proto.CardService/SetCardDates":          roles.MemberRole,
		"/proto.CardService/MarkCardComplete":      roles.MemberRole,
		"/proto.CardService/AddCardAttachment":     roles.MemberRole,
		"/proto.CardService/RemoveCardAttachment":  roles.MemberRole,
		"/proto.CardService/AddCardComment":        roles.MemberRole,
		"/proto.CardService/RemoveCardComment":     roles.MemberRole,
		"/proto.CardService/AddCardMembers":        roles.MemberRole,
		"/proto.CardService/RemoveCardMembers":     roles.MemberRole,
		"/proto.CardService/ArchiveCard":           roles.MemberRole,
		"/proto.CardService/RestoreCard":           roles.MemberRole,
		"/proto.CardService/DeleteCard":            roles.AdminRole,
		// Add other methods here...
	}

	roleHierarchy = map[string]int{
		roles.ObserverRole: 1,
		roles.MemberRole:   2,
		roles.AdminRole:    3,
		roles.OwnerRole:    4,
	}
)

type Validator struct {
	db *gorm.DB
}

func NewValidator(db *gorm.DB) *Validator {
	return &Validator{db: db}
}

func (v *Validator) checkBoardUserRole(userID uint64, boardID uint64, requiredRole string) error {
	var permission models.Permission
	if err := v.db.Select("role").Where("board_id = ? AND user_id = ?", boardID, userID).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorhandler.NewAPIError(httpcodes.ErrForbidden, "Permission not found")
		}
		return errorhandler.NewGrpcInternalError()
	}

	// Check if the user's role is sufficient
	if roleHierarchy[permission.Role] < roleHierarchy[requiredRole] {
		return errorhandler.NewAPIError(httpcodes.ErrForbidden, "Insufficient permissions")
	}

	return nil
}

func (v *Validator) checkBoardVisibility(userID, boardID uint64) error {
	var result struct {
		Visibility string
		UserID     uint64
	}

	query := `SELECT boards.visibility, permissions.user_id 
              FROM boards 
              LEFT JOIN permissions ON boards.id = permissions.board_id AND permissions.user_id = ?
              WHERE boards.id = ?`

	if err := v.db.Raw(query, userID, boardID).Scan(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorhandler.NewAPIError(httpcodes.ErrNotFound, "Board not found")
		}
		return errorhandler.NewGrpcInternalError()
	}

	if result.Visibility == "private" && result.UserID != userID {
		return errorhandler.NewAPIError(httpcodes.ErrForbidden, "Permission not found")
	}

	return nil
}

func (v *Validator) ValidationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// Define the methods that don't require user role check

	userID, err := helpers.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	boardID, err := helpers.ExtractBoardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, contextkeys.UserIDKey{}, userID)
	ctx = context.WithValue(ctx, contextkeys.BoardIDKey{}, boardID)

	// Check if the method requires user role check

	if !checkVisibility[info.FullMethod] {
		// Define the required role for each method

		requiredRole, ok := checkRoles[info.FullMethod]
		if !ok {
			return nil, errorhandler.NewAPIError(httpcodes.ErrForbidden, "Invalid method")
		}

		// Check the user's role
		if err := v.checkBoardUserRole(userID, boardID, requiredRole); err != nil {
			return nil, err
		}
	}

	if checkVisibility[info.FullMethod] {
		// Check the board's visibility
		err := v.checkBoardVisibility(userID, boardID)
		if err != nil {
			return nil, err
		}
	}

	switch info.FullMethod {
	// Card Service
	case "/proto.CardService/CreateCard":
		req := req.(*pb.CreateCardRequest)
		if err := validateCreateCardRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/GetCardByID":
		req := req.(*pb.GetCardByIDRequest)
		if err := validateGetCardByIDRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/GetCardsByList":
		req := req.(*pb.GetCardsByListRequest)
		if err := validateGetCardsByListRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/GetCardsByBoard":
		req := req.(*pb.GetCardsByBoardRequest)
		if err := validateGetCardsByBoardRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/MoveCardPosition":
		req := req.(*pb.MoveCardPositionRequest)
		if err := validateMoveCardPositionRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/UpdateCardName":
		req := req.(*pb.UpdateCardNameRequest)
		if err := validateUpdateCardNameRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/UpdateCardDescription":
		req := req.(*pb.UpdateCardDescriptionRequest)
		if err := validateUpdateCardDescriptionRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/AddCardLabel":
		req := req.(*pb.AddCardLabelRequest)
		if err := validateAddCardLabelRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RemoveCardLabel":
		req := req.(*pb.RemoveCardLabelRequest)
		if err := validateRemoveCardLabelRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/SetCardDates":
		req := req.(*pb.SetCardDatesRequest)
		if err := validateSetCardDatesRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/MarkCardComplete":
		req := req.(*pb.MarkCardCompleteRequest)
		if err := validateMarkCardCompleteRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/AddCardAttachment":
		req := req.(*pb.AddCardAttachmentRequest)
		if err := validateAddCardAttachmentRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RemoveCardAttachment":
		req := req.(*pb.RemoveCardAttachmentRequest)
		if err := validateRemoveCardAttachmentRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/AddCardComment":
		req := req.(*pb.AddCardCommentRequest)
		if err := validateAddCardCommentRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RemoveCardComment":
		req := req.(*pb.RemoveCardCommentRequest)
		if err := validateRemoveCardCommentRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/AddCardMembers":
		req := req.(*pb.AddCardMembersRequest)
		if err := validateAddCardMembersRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RemoveCardMembers":
		req := req.(*pb.RemoveCardMembersRequest)
		if err := validateRemoveCardMembersRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/ArchiveCard":
		req := req.(*pb.ArchiveCardRequest)
		if err := validateArchiveCardRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/RestoreCard":
		req := req.(*pb.RestoreCardRequest)
		if err := validateRestoreCardRequest(req); err != nil {
			return nil, err
		}
	case "/proto.CardService/DeleteCard":
		req := req.(*pb.DeleteCardRequest)
		if err := validateDeleteCardRequest(req); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

func validateCreateCardRequest(req *pb.CreateCardRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

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
	if len(req.Name) > 50 {
		if _, exists := fieldErrors["Name"]; !exists {
			fieldErrors["Name"] = errorhandler.FieldError{
				Code:    fielderrors.ErrMaxLength,
				Message: "Name cannot be longer than 50 characters",
				Field:   "Name",
			}
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateGetCardByIDRequest(req *pb.GetCardByIDRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateGetCardsByListRequest(req *pb.GetCardsByListRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

	if req.ListID == 0 {
		fieldErrors["ListID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "ListID is required",
			Field:   "ListID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateGetCardsByBoardRequest(req *pb.GetCardsByBoardRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

	if req.BoardID == 0 {
		fieldErrors["BoardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "BoardID is required",
			Field:   "BoardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateMoveCardPositionRequest(req *pb.MoveCardPositionRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	if req.NewPosition < 0 {
		fieldErrors["Position"] = errorhandler.FieldError{
			Code:    fielderrors.ErrOutOfRange,
			Message: "Position must be greater than or equal to 0",
			Field:   "Position",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateUpdateCardNameRequest(req *pb.UpdateCardNameRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

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

func validateUpdateCardDescriptionRequest(req *pb.UpdateCardDescriptionRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

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

func validateAddCardLabelRequest(req *pb.AddCardLabelRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

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

func validateRemoveCardLabelRequest(req *pb.RemoveCardLabelRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

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

func validateSetCardDatesRequest(req *pb.SetCardDatesRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateMarkCardCompleteRequest(req *pb.MarkCardCompleteRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateAddCardAttachmentRequest(req *pb.AddCardAttachmentRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

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

func validateRemoveCardAttachmentRequest(req *pb.RemoveCardAttachmentRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

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

func validateAddCardCommentRequest(req *pb.AddCardCommentRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

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

func validateRemoveCardCommentRequest(req *pb.RemoveCardCommentRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

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

func validateAddCardMembersRequest(req *pb.AddCardMembersRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

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

func validateRemoveCardMembersRequest(req *pb.RemoveCardMembersRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

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

func validateArchiveCardRequest(req *pb.ArchiveCardRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateRestoreCardRequest(req *pb.RestoreCardRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateDeleteCardRequest(req *pb.DeleteCardRequest) *errorhandler.APIError {
	fieldErrors := make(map[string]errorhandler.FieldError)

	if req.CardID == 0 {
		fieldErrors["CardID"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "CardID is required",
			Field:   "CardID",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}
