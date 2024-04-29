package middlewares

import (
	"context"

	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
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
	// Board Service
	case "/proto.BoardService/CreateBoard":
		if err := validateCreateBoardRequest(req.(*pb_board.CreateBoardRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/GetBoardByID":
		if err := validateGetBoardByIDRequest(req.(*pb_board.GetBoardByIDRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/GetBoardList":
		if err := validateGetBoardListRequest(req.(*pb_board.GetBoardListRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/GetBoardMembers":
		if err := validateGetBoardMembersRequest(req.(*pb_board.GetBoardMembersRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/UpdateBoardName":
		if err := validateUpdateBoardNameRequest(req.(*pb_board.UpdateBoardNameRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/AddBoardUsers":
		if err := validateAddBoardUsersRequest(req.(*pb_board.AddBoardUsersRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/RemoveBoardUsers":
		if err := validateRemoveBoardUsersRequest(req.(*pb_board.RemoveBoardUsersRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/AssignBoardUsersRole":
		if err := validateAssignBoardUsersRoleRequest(req.(*pb_board.AssignBoardUsersRoleRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/ChangeBoardOwner":
		if err := validateChangeBoardOwnerRequest(req.(*pb_board.ChangeBoardOwnerRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/GetArchivedBoardList":
		if err := validateGetArchivedBoardListRequest(req.(*pb_board.GetArchivedBoardListRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/RestoreBoard":
		if err := validateRestoreBoardRequest(req.(*pb_board.RestoreBoardRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/ArchiveBoard":
		if err := validateArchiveBoardRequest(req.(*pb_board.ArchiveBoardRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/DeleteBoard":
		if err := validateDeleteBoardRequest(req.(*pb_board.DeleteBoardRequest)); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

// Board Service

func validateCreateBoardRequest(req *pb_board.CreateBoardRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.Name == "" {
		fieldErrors["Name"] = errorhandlers.FieldError{
			Field:   "name",
			Message: "Board name cannot be empty",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateGetBoardByIDRequest(_ *pb_board.GetBoardByIDRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateGetBoardListRequest(_ *pb_board.GetBoardListRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateGetBoardMembersRequest(_ *pb_board.GetBoardMembersRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateUpdateBoardNameRequest(req *pb_board.UpdateBoardNameRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.Name == "" {
		fieldErrors["Name"] = errorhandlers.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Name is required",
			Field:   "Name",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateAddBoardUsersRequest(req *pb_board.AddBoardUsersRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if len(req.UserIDs) == 0 {
		fieldErrors["Name"] = errorhandlers.FieldError{
			Field:   "UserIDs",
			Message: "User IDs cannot be empty",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateRemoveBoardUsersRequest(req *pb_board.RemoveBoardUsersRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if len(req.UserIDs) == 0 {
		fieldErrors["Name"] = errorhandlers.FieldError{
			Field:   "UserIDs",
			Message: "User IDs cannot be empty",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateAssignBoardUsersRoleRequest(req *pb_board.AssignBoardUsersRoleRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if len(req.UserIDs) == 0 {
		fieldErrors["Name"] = errorhandlers.FieldError{
			Field:   "userIDs",
			Message: "User IDs cannot be empty",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateChangeBoardOwnerRequest(req *pb_board.ChangeBoardOwnerRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.NewOwnerID == 0 {
		fieldErrors["NewOwnerID"] = errorhandlers.FieldError{
			Field:   "NewOwnerID",
			Message: "New owner ID cannot be zero",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateGetArchivedBoardListRequest(_ *pb_board.GetArchivedBoardListRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateRestoreBoardRequest(_ *pb_board.RestoreBoardRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateArchiveBoardRequest(_ *pb_board.ArchiveBoardRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)

}

func validateDeleteBoardRequest(_ *pb_board.DeleteBoardRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)

}
