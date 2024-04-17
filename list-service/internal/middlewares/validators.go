package middlewares

import (
	"context"

	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	pb "github.com/sm888sm/halten-backend/list-service/api/pb"
	"google.golang.org/grpc"
)

func ValidationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	switch info.FullMethod {
	// List Service
	case "/proto.ListService/CreateList":
		if err := validateCreateListRequest(req.(*pb.CreateListRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/GetListsByBoard":
		if err := validateGetListsByBoardRequest(req.(*pb.GetListsByBoardRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/UpdateList":
		if err := validateUpdateListRequest(req.(*pb.UpdateListRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/DeleteList":
		if err := validateDeleteListRequest(req.(*pb.DeleteListRequest)); err != nil {
			return nil, err
		}
	case "/proto.ListService/MoveListPosition":
		if err := validateMoveListPositionRequest(req.(*pb.MoveListPositionRequest)); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

func validateCreateListRequest(req *pb.CreateListRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Name == "" {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "name",
			Message: "List name cannot be empty",
		})
	}

	if req.BoardId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardId",
			Message: "Board ID cannot be zero",
		})
	}

	if req.UserId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userId",
			Message: "User ID cannot be zero",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateGetListsByBoardRequest(req *pb.GetListsByBoardRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.BoardId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardId",
			Message: "Board ID cannot be zero",
		})
	}

	if req.UserId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userId",
			Message: "User ID cannot be zero",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateUpdateListRequest(req *pb.UpdateListRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Id == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "listId",
			Message: "List ID cannot be zero",
		})
	}

	if req.Name == "" {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "name",
			Message: "List name cannot be empty",
		})
	}

	if req.UserId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userId",
			Message: "User ID cannot be zero",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateDeleteListRequest(req *pb.DeleteListRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Id == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "listId",
			Message: "List ID cannot be zero",
		})
	}

	if req.BoardId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardId",
			Message: "Board ID cannot be zero",
		})
	}

	if req.UserId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userId",
			Message: "User ID cannot be zero",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateMoveListPositionRequest(req *pb.MoveListPositionRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Id == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "listId",
			Message: "List ID cannot be zero",
		})
	}

	if req.NewPosition < 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "newPosition",
			Message: "New position cannot be negative",
		})
	}

	if req.BoardId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardId",
			Message: "Board ID cannot be zero",
		})
	}

	if req.UserId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userId",
			Message: "User ID cannot be zero",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}
