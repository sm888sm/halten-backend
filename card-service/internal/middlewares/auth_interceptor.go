package middlewares

import (
	"context"

	external_services "github.com/sm888sm/halten-backend/card-service/external/services"

	pb_auth "github.com/sm888sm/halten-backend/user-service/api/pb"

	"github.com/sm888sm/halten-backend/common/constants/contextkeys"
	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/constants/roles"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	"github.com/sm888sm/halten-backend/common/helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	checkRoleException = map[string]bool{
		"/proto.CardService/GetCardByID":     true,
		"/proto.CardService/GetCardsByList":  true,
		"/proto.CardService/GetCardsByBoard": true,
	}

	checkRole = map[string]string{
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
)

type AuthInterceptor struct {
	db  *gorm.DB
	svc *external_services.Services
}

func NewAuthInterceptor(db *gorm.DB, svc *external_services.Services) *AuthInterceptor {
	return &AuthInterceptor{db: db}
}

func (v *AuthInterceptor) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	_, isException := checkRoleException[info.FullMethod]
	if !isException {

		requiredRole, ok := checkRole[info.FullMethod]
		if !ok {
			return nil, status.Errorf(codes.Unavailable, errorhandler.NewAPIError(httpcodes.ErrForbidden, "Invalid method").Error())
		}

		authService, err := v.svc.GetAuthClient()
		if err != nil {
			return nil, errorhandler.NewGrpcBadRequestError()
		}

		// Extract userID, boardID and cardID from meta

		userID, err := helpers.ExtractUserIDFromContext(ctx)
		if err != nil {
			return nil, err
		}

		boardID, err := helpers.ExtractBoardIDFromContext(ctx)
		if err != nil {
			return nil, err
		}

		// Insert userID, boardID and cardID to context

		ctx = context.WithValue(ctx, contextkeys.UserIDKey{}, userID)
		ctx = context.WithValue(ctx, contextkeys.BoardIDKey{}, boardID)

		if _, err := authService.CheckBoardUserRole(ctx, &pb_auth.CheckBoardUserRoleRequest{
			UserID:       userID,
			BoardID:      boardID,
			RequiredRole: requiredRole,
		}); err != nil {
			return nil, err
		}

		// TODO : Move checkVisibility to gateway

		// if checkVisibility[info.FullMethod] {
		// 	// Check the board's visibility
		// 	if _, err := authService.CheckBoardVisibility(ctx, &pb_auth.CheckBoardVisibilityRequest{
		// 		UserID:  userID,
		// 		BoardID: boardID,
		// 	}); err != nil {
		// 		return nil, err
		// 	}
		// }
	}

	return handler(ctx, req)
}
