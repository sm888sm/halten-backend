package middlewares

import (
	"context"

	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
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
	// Board Service
	case "/proto.BoardService/CreateBoard":
		if err := validateCreateBoardRequest(ctx, req.(*pb_board.CreateBoardRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/GetBoardByID":
		if err := validateGetBoardByIDRequest(ctx, req.(*pb_board.GetBoardByIDRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/GetBoardList":
		if err := validateGetBoardListRequest(ctx, req.(*pb_board.GetBoardListRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/GetBoardMembers":
		if err := validateGetBoardMembersRequest(ctx, req.(*pb_board.GetBoardMembersRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/UpdateBoardName":
		if err := validateUpdateBoardNameRequest(ctx, req.(*pb_board.UpdateBoardNameRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/AddBoardUsers":
		if err := validateAddBoardUsersRequest(ctx, req.(*pb_board.AddBoardUsersRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/RemoveBoardUsers":
		if err := validateRemoveBoardUsersRequest(ctx, req.(*pb_board.RemoveBoardUsersRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/AssignBoardUserRole":
		if err := validateAssignBoardUserRoleRequest(ctx, req.(*pb_board.AssignBoardUserRoleRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/ChangeBoardOwner":
		if err := validateChangeBoardOwnerRequest(ctx, req.(*pb_board.ChangeBoardOwnerRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/GetArchivedBoardList":
		if err := validateGetArchivedBoardListRequest(ctx, req.(*pb_board.GetArchivedBoardListRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/RestoreBoard":
		if err := validateRestoreBoardRequest(ctx, req.(*pb_board.RestoreBoardRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/ArchiveBoard":
		if err := validateArchiveBoardRequest(ctx, req.(*pb_board.ArchiveBoardRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/DeleteBoard":
		if err := validateDeleteBoardRequest(ctx, req.(*pb_board.DeleteBoardRequest)); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

// Board Service

func validateCreateBoardRequest(ctx context.Context, req *pb_board.CreateBoardRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.Name == "" {
		fieldErrors["Name"] = errorhandler.FieldError{
			Field:   "name",
			Message: "Board name cannot be empty",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateGetBoardByIDRequest(ctx context.Context, req *pb_board.GetBoardByIDRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateGetBoardListRequest(ctx context.Context, req *pb_board.GetBoardListRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateGetBoardMembersRequest(ctx context.Context, req *pb_board.GetBoardMembersRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateUpdateBoardNameRequest(ctx context.Context, req *pb_board.UpdateBoardNameRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
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

func validateAddBoardUsersRequest(ctx context.Context, req *pb_board.AddBoardUsersRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	if len(req.UserIDs) == 0 {
		fieldErrors["Name"] = errorhandler.FieldError{
			Field:   "UserIDs",
			Message: "User IDs cannot be empty",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateRemoveBoardUsersRequest(ctx context.Context, req *pb_board.RemoveBoardUsersRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	if len(req.UserIDs) == 0 {
		fieldErrors["Name"] = errorhandler.FieldError{
			Field:   "UserIDs",
			Message: "User IDs cannot be empty",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateAssignBoardUserRoleRequest(ctx context.Context, req *pb_board.AssignBoardUserRoleRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	if len(req.UserIDs) == 0 {
		fieldErrors["Name"] = errorhandler.FieldError{
			Field:   "userIDs",
			Message: "User IDs cannot be empty",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateChangeBoardOwnerRequest(ctx context.Context, req *pb_board.ChangeBoardOwnerRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	if req.NewOwnerID == 0 {
		fieldErrors["NewOwnerID"] = errorhandler.FieldError{
			Field:   "NewOwnerID",
			Message: "New owner ID cannot be zero",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateGetArchivedBoardListRequest(ctx context.Context, req *pb_board.GetArchivedBoardListRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateRestoreBoardRequest(ctx context.Context, req *pb_board.RestoreBoardRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateArchiveBoardRequest(ctx context.Context, req *pb_board.ArchiveBoardRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)

}

func validateDeleteBoardRequest(ctx context.Context, req *pb_board.DeleteBoardRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)

}
