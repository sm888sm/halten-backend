package middlewares

import (
	"context"

	"github.com/sm888sm/halten-backend/common/constants/fielderrors"
	"github.com/sm888sm/halten-backend/common/errorhandlers"
	pb_list "github.com/sm888sm/halten-backend/list-service/api/pb"
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
	case "/proto.ListService/CreateList":
		if err := validateCreateListRequest(req.(*pb_list.CreateListRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/GetListByID":
		if err := validateGetListByIDRequest(req.(*pb_list.GetListByIDRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/GetListsByBoard":
		if err := validateGetListsByBoardRequest(req.(*pb_list.GetListsByBoardRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/UpdateListName":
		if err := validateUpdateListNameRequest(req.(*pb_list.UpdateListNameRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/MoveListPosition":
		if err := validateMoveListPositionRequest(req.(*pb_list.MoveListPositionRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/ArchiveList":
		if err := validateArchiveListRequest(req.(*pb_list.ArchiveListRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/RestoreList":
		if err := validateRestoreListRequest(req.(*pb_list.RestoreListRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/DeleteList":
		if err := validateDeleteListRequest(req.(*pb_list.DeleteListRequest)); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

func validateCreateListRequest(req *pb_list.CreateListRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.Name == "" {
		fieldErrors["Name"] = errorhandlers.FieldError{
			Field:   "name",
			Message: "List name cannot be empty",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateGetListByIDRequest(_ *pb_list.GetListByIDRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateGetListsByBoardRequest(_ *pb_list.GetListsByBoardRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateUpdateListNameRequest(req *pb_list.UpdateListNameRequest) error {
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

func validateMoveListPositionRequest(req *pb_list.MoveListPositionRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	if req.Position == 0 {
		fieldErrors["Position"] = errorhandlers.FieldError{
			Field:   "Position",
			Message: "New position cannot be zero",
		}
	}

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateArchiveListRequest(_ *pb_list.ArchiveListRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateRestoreListRequest(_ *pb_list.RestoreListRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateDeleteListRequest(_ *pb_list.DeleteListRequest) error {
	fieldErrors := make(map[string]errorhandlers.FieldError)

	return errorhandlers.CreateGrpcErrorFromFieldErrors(fieldErrors)
}
