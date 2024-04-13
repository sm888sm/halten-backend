package middlewares

import (
	"context"

	pb "github.com/sm888sm/halten-backend/card-service/api/pb"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	"google.golang.org/grpc"
)

func ValidationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	switch info.FullMethod {
	// Card Service
	case "/proto.CardService/CreateCard":
		if err := validateCreateCardRequest(req.(*pb.CreateCardRequest)); err != nil {
			return nil, err
		}
	case "/proto.CardService/GetCardsByList":
		if err := validateGetCardsByListRequest(req.(*pb.GetCardsByListRequest)); err != nil {
			return nil, err
		}
	case "/proto.CardService/UpdateCard":
		if err := validateUpdateCardRequest(req.(*pb.UpdateCardRequest)); err != nil {
			return nil, err
		}
	case "/proto.CardService/DeleteCard":
		if err := validateDeleteCardRequest(req.(*pb.DeleteCardRequest)); err != nil {
			return nil, err
		}
	case "/proto.CardService/MoveCardPosition":
		if err := validateMoveCardPositionRequest(req.(*pb.MoveCardPositionRequest)); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

func validateCreateCardRequest(req *pb.CreateCardRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Name == "" {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "name",
			Message: "Card name cannot be empty",
		})
	}

	if req.ListId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "listId",
			Message: "List ID cannot be zero",
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

func validateGetCardsByListRequest(req *pb.GetCardsByListRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.ListId == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "listId",
			Message: "List ID cannot be zero",
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

func validateUpdateCardRequest(req *pb.UpdateCardRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Id == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "id",
			Message: "Card ID cannot be zero",
		})
	}

	if req.Name == "" {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "name",
			Message: "Card name cannot be empty",
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

func validateDeleteCardRequest(req *pb.DeleteCardRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Id == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "id",
			Message: "Card ID cannot be zero",
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

func validateMoveCardPositionRequest(req *pb.MoveCardPositionRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Id == 0 {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "id",
			Message: "Card ID cannot be zero",
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
