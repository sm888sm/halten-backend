package middlewares

import (
	"context"
	"errors"

	pb_board "github.com/sm888sm/halten-backend/board-service/api/pb"
	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
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
	case "/proto.BoardService/AssignBoardUserRole":
		if err := validateAssignBoardUserRoleRequest(req.(*pb_board.AssignBoardUserRoleRequest)); err != nil {
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
	var fieldErrors []errorhandler.FieldError

	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.Name == "" {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "name",
			Message: "Board name cannot be empty",
		})
	}

	// Add more validation rules as needed

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateGetBoardByIDRequest(req *pb_board.GetBoardByIDRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.BoardID == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardID",
			Message: "Board ID cannot be empty",
		})
	}

	// Add more validation rules as needed

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateGetBoardListRequest(req *pb_board.GetBoardListRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req == nil {
		return errors.New("request cannot be nil")
	}

	// Add more validation rules as needed

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateGetBoardMembersRequest(req *pb_board.GetBoardMembersRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.BoardID == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardID",
			Message: "Board ID cannot be zero",
		})
	}

	// Add more validation rules as needed

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateUpdateBoardNameRequest(req *pb_board.UpdateBoardNameRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.BoardID == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardID",
			Message: "Board ID cannot be zero",
		})
	}

	if req.Name == "" {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "name",
			Message: "Board name cannot be empty",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateAddBoardUsersRequest(req *pb_board.AddBoardUsersRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.BoardID == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardID",
			Message: "Board ID cannot be zero",
		})
	}

	if len(req.UserIDs) == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userIDs",
			Message: "User IDs cannot be empty",
		})
	}

	// Add more validation rules as needed

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateRemoveBoardUsersRequest(req *pb_board.RemoveBoardUsersRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.BoardID == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardID",
			Message: "Board ID cannot be zero",
		})
	}

	if len(req.UserIDs) == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userIDs",
			Message: "User IDs cannot be empty",
		})
	}

	// Add more validation rules as needed

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateAssignBoardUserRoleRequest(req *pb_board.AssignBoardUserRoleRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.BoardID == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardID",
			Message: "Board ID cannot be zero",
		})
	}

	if len(req.UserIDs) == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userIDs",
			Message: "User IDs cannot be empty",
		})
	}

	// Add more validation rules as needed

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateChangeBoardOwnerRequest(req *pb_board.ChangeBoardOwnerRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.BoardID == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardID",
			Message: "Board ID cannot be zero",
		})
	}

	if req.NewOwnerID == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "newOwnerID",
			Message: "New owner ID cannot be zero",
		})
	}

	// Add more validation rules as needed

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateGetArchivedBoardListRequest(req *pb_board.GetArchivedBoardListRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req == nil {
		return errors.New("request cannot be nil")
	}

	// Add more validation rules as needed

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateRestoreBoardRequest(req *pb_board.RestoreBoardRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.BoardID == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardID",
			Message: "Board ID cannot be zero",
		})
	}

	// Add more validation rules as needed

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateArchiveBoardRequest(req *pb_board.ArchiveBoardRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.BoardID == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardID",
			Message: "Board ID cannot be zero",
		})
	}

	// Add more validation rules as needed

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateDeleteBoardRequest(req *pb_board.DeleteBoardRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.BoardID == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardID",
			Message: "Board ID cannot be zero",
		})
	}

	// Add more validation rules as needed

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}
