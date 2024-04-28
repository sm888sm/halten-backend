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
		if err := validateCreateListRequest(ctx, req.(*pb_list.CreateListRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/GetListByID":
		if err := validateGetListByIDRequest(ctx, req.(*pb_list.GetListByIDRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/GetListsByBoard":
		if err := validateGetListsByBoardRequest(ctx, req.(*pb_list.GetListsByBoardRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/UpdateListName":
		if err := validateUpdateListNameRequest(ctx, req.(*pb_list.UpdateListNameRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/MoveListPosition":
		if err := validateMoveListPositionRequest(ctx, req.(*pb_list.MoveListPositionRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/ArchiveList":
		if err := validateArchiveListRequest(ctx, req.(*pb_list.ArchiveListRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/RestoreList":
		if err := validateRestoreListRequest(ctx, req.(*pb_list.RestoreListRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/DeleteList":
		if err := validateDeleteListRequest(ctx, req.(*pb_list.DeleteListRequest)); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

func validateCreateListRequest(ctx context.Context, req *pb_list.CreateListRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	if req.List.Name == "" {
		fieldErrors["Name"] = errorhandler.FieldError{
			Field:   "name",
			Message: "List name cannot be empty",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateGetListByIDRequest(ctx context.Context, req *pb_list.GetListByIDRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateGetListsByBoardRequest(ctx context.Context, req *pb_list.GetListsByBoardRequest) *errorhandler.APIError {
	fieldErrors, apiErr := validateUserAndBoardID(ctx)
	if apiErr != nil {
		return apiErr
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateUpdateListNameRequest(ctx context.Context, req *pb_list.UpdateListNameRequest) *errorhandler.APIError {
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

func validateMoveListPositionRequest(ctx context.Context, req *pb_list.MoveListPositionRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	if req.NewPosition == 0 {
		fieldErrors["NewPosition"] = errorhandler.FieldError{
			Field:   "NewPosition",
			Message: "New position cannot be zero",
		}
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateArchiveListRequest(ctx context.Context, req *pb_list.ArchiveListRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateRestoreListRequest(ctx context.Context, req *pb_list.RestoreListRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}

func validateDeleteListRequest(ctx context.Context, req *pb_list.DeleteListRequest) *errorhandler.APIError {
	fieldErrors, err := validateUserAndBoardID(ctx)
	if err != nil {
		return err
	}

	return errorhandler.CreateAPIErrorFromFieldErrors(fieldErrors)
}
