package middlewares

import (
	"context"

	external_services "github.com/sm888sm/halten-backend/list-service/external/services"

	pb_auth "github.com/sm888sm/halten-backend/user-service/api/pb"

	"github.com/sm888sm/halten-backend/common/constants/contextkeys"
	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/constants/roles"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	"github.com/sm888sm/halten-backend/common/helpers"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	checkRoleException = map[string]bool{
		"/proto.BoardService/GetListByID":     true,
		"/proto.BoardService/GetListsByBoard": true,

		// Add other methods here...
	}

	checkRole = map[string]string{
		"/proto.BoardService/UpdateListName":  roles.MemberRole,
		"/proto.ListService/CreateList":       roles.MemberRole,
		"/proto.ListService/GetListsByID":     roles.MemberRole,
		"/proto.ListService/GetListsByBoard":  roles.MemberRole,
		"/proto.ListService/MoveListPosition": roles.MemberRole,
		"/proto.ListService/ArchiveList":      roles.MemberRole,
		"/proto.ListService/RestoreList":      roles.MemberRole,
		"/proto.ListService/DeleteList":       roles.AdminRole,
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
			return nil, errorhandler.NewAPIError(httpcodes.ErrForbidden, "Invalid method")
		}

		authService, err := v.svc.GetAuthClient()
		if err != nil {
			return nil, errorhandler.NewGrpcBadRequestError()
		}

		// Extract userID and boardID from meta

		userID, err := helpers.ExtractUserIDFromContext(ctx)
		if err != nil {
			return nil, err
		}

		boardID, err := helpers.ExtractBoardIDFromContext(ctx)
		if err != nil {
			return nil, err
		}

		// Insert userID and boardID to context

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
