package middlewares

import (
	"context"
	"fmt"

	pb "github.com/sm888sm/halten-backend/board-service/api/pb"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	"google.golang.org/grpc"
)

func ValidationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	switch info.FullMethod {
	// Board Service
	case "/proto.BoardService/CreateBoard":
		if err := validateCreateBoardRequest(req.(*pb.CreateBoardRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/GetBoardByID":
		if err := validateGetBoardByIDRequest(req.(*pb.GetBoardByIDRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/GetBoardList":
		if err := validateGetBoardListRequest(req.(*pb.GetBoardListRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/GetBoardUsers":
		if err := validateGetBoardUsersRequest(req.(*pb.GetBoardUsersRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/UpdateBoard":
		if err := validateUpdateBoardRequest(req.(*pb.UpdateBoardRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/DeleteBoard":
		if err := validateDeleteBoardRequest(req.(*pb.DeleteBoardRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/AddUser":
		if err := validateAddUsersRequest(req.(*pb.AddUsersRequest)); err != nil {
			return nil, err
		}
	case "/proto.BoardService/RemoveUser":
		if err := validateRemoveUsersRequest(req.(*pb.RemoveUsersRequest)); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

// Board Service
func validateCreateBoardRequest(req *pb.CreateBoardRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Name == "" {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "name",
			Message: "Board name cannot be empty",
		})
	}

	if req.UserId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userId",
			Message: "User ID cannot be zero",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateGetBoardByIDRequest(req *pb.GetBoardByIDRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Id == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "id",
			Message: "Board ID cannot be zero",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateGetBoardListRequest(req *pb.GetBoardListRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.UserId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userId",
			Message: "User ID cannot be zero",
		})
	}

	if req.PageNumber == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "page",
			Message: "Page cannot be zero",
		})
	}

	if req.PageSize == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "pageSize",
			Message: "Page size cannot be zero",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateGetBoardUsersRequest(req *pb.GetBoardUsersRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Id == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardId",
			Message: "Board ID cannot be zero",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateUpdateBoardRequest(req *pb.UpdateBoardRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Id == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardId",
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
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateDeleteBoardRequest(req *pb.DeleteBoardRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Id == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardId",
			Message: "Board ID cannot be zero",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateAddUsersRequest(req *pb.AddUsersRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Id == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardId",
			Message: "Board ID cannot be zero",
		})
	}

	if len(req.UserIds) == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userId",
			Message: "User ID cannot be zero",
		})
	}

	if len(req.UserIds) > 10 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userId",
			Message: "Cannot process more than 10 users at a time",
		})
	}

	for _, userId := range req.UserIds {
		if !isValidAccountNumber(userId) {
			return errorhandler.NewAPIError(errorhandler.ErrBadRequest, fmt.Sprintf("User ID %d is not a valid account number", userId))
		}
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateRemoveUsersRequest(req *pb.RemoveUsersRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Id == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "boardId",
			Message: "Board ID cannot be zero",
		})
	}

	if len(req.UserIds) == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userId",
			Message: "User ID cannot be zero",
		})
	}

	if len(req.UserIds) > 10 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "userId",
			Message: "Cannot process more than 10 users at a time",
		})
	}

	for _, userId := range req.UserIds {
		if !isValidAccountNumber(userId) {
			return errorhandler.NewAPIError(errorhandler.ErrBadRequest, fmt.Sprintf("User ID %d is not a valid account number", userId))
		}
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func isValidAccountNumber(userId uint64) bool {
	// Check if the user ID is a positive number
	return userId > 0
}
