package middlewares

import (
	"context"

	"github.com/sm888sm/halten-backend/common/constants/fielderrors"
	"github.com/sm888sm/halten-backend/common/errorhandler"
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
	fieldErrors := make(map[string]errorhandler.FieldError)

	if req.List.Name == "" {
		fieldErrors["Name"] = errorhandler.FieldError{
			Field:   "name",
			Message: "List name cannot be empty",
		}
	}

	return errorhandler.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateGetListByIDRequest(_ *pb_list.GetListByIDRequest) error {
	fieldErrors := make(map[string]errorhandler.FieldError)

	return errorhandler.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateGetListsByBoardRequest(_ *pb_list.GetListsByBoardRequest) error {
	fieldErrors := make(map[string]errorhandler.FieldError)

	return errorhandler.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateUpdateListNameRequest(req *pb_list.UpdateListNameRequest) error {
	fieldErrors := make(map[string]errorhandler.FieldError)

	if req.Name == "" {
		fieldErrors["Name"] = errorhandler.FieldError{
			Code:    fielderrors.ErrRequired,
			Message: "Name is required",
			Field:   "Name",
		}
	}

	return errorhandler.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateMoveListPositionRequest(req *pb_list.MoveListPositionRequest) error {
	fieldErrors := make(map[string]errorhandler.FieldError)

	if req.NewPosition == 0 {
		fieldErrors["NewPosition"] = errorhandler.FieldError{
			Field:   "NewPosition",
			Message: "New position cannot be zero",
		}
	}

	return errorhandler.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateArchiveListRequest(_ *pb_list.ArchiveListRequest) error {
	fieldErrors := make(map[string]errorhandler.FieldError)

	return errorhandler.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateRestoreListRequest(_ *pb_list.RestoreListRequest) error {
	fieldErrors := make(map[string]errorhandler.FieldError)

	return errorhandler.CreateGrpcErrorFromFieldErrors(fieldErrors)
}

func validateDeleteListRequest(_ *pb_list.DeleteListRequest) error {
	fieldErrors := make(map[string]errorhandler.FieldError)

	return errorhandler.CreateGrpcErrorFromFieldErrors(fieldErrors)
}
