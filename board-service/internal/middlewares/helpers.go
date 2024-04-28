package middlewares

import (
	"context"

	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/errorhandler"
	"google.golang.org/grpc/metadata"
)

func validateUserAndBoardID(ctx context.Context) (map[string]errorhandler.FieldError, *errorhandler.APIError) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errorhandler.NewAPIError(httpcodes.ErrInternalServerError, "missing context metadata")
	}

	fieldErrors := make(map[string]errorhandler.FieldError)

	userIDs, userIDsExist := md["userID"]
	boardIDs, boardIDsExist := md["boardID"]

	// Check if userID and boardID exist in the metadata
	if !userIDsExist || len(userIDs) == 0 {
		fieldErrors["userID"] = errorhandler.FieldError{
			Field:   "userID",
			Message: "userID cannot be empty",
		}
	}

	if !boardIDsExist || len(boardIDs) == 0 {
		fieldErrors["boardID"] = errorhandler.FieldError{
			Field:   "boardID",
			Message: "boardID cannot be empty",
		}
	}

	return fieldErrors, nil
}
