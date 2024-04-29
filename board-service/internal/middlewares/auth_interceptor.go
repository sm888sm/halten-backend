package middlewares

import (
	"context"

	external_services "github.com/sm888sm/halten-backend/board-service/external/services"

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
		"/proto.BoardService/CreateBoard":          true,
		"/proto.BoardService/GetBoardByID":         true,
		"/proto.BoardService/GetBoardList":         true,
		"/proto.BoardService/GetArchivedBoardList": true,
		"/proto.BoardService/GetBoardMembers":      true,

		// Add other methods here...
	}

	checkRole = map[string]string{
		"/proto.BoardService/UpdateBoardName":       roles.AdminRole,
		"/proto.BoardService/AddBoardUsers":         roles.AdminRole,
		"/proto.BoardService/RemoveBoardUsers":      roles.AdminRole,
		"/proto.BoardService/AssignBoardUsersRole":  roles.AdminRole,
		"/proto.BoardService/ChangeBoardOwner":      roles.OwnerRole,
		"/proto.BoardService/ChangeBoardVisibility": roles.AdminRole,
		"/proto.BoardService/AddLabel":              roles.MemberRole,
		"/proto.BoardService/RemoveLabel":           roles.MemberRole,
		"/proto.BoardService/RestoreBoard":          roles.AdminRole,
		"/proto.BoardService/ArchiveBoard":          roles.AdminRole,
		"/proto.BoardService/DeleteBoard":           roles.OwnerRole,
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
	}

	return handler(ctx, req)
}
