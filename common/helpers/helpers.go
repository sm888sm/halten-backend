package helpers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/sm888sm/halten-backend/common/errorhandler"
	"google.golang.org/grpc/metadata"
)

func ExtractUserIDFromContext(ctx context.Context) (uint64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errorhandler.NewAPIError(http.StatusInternalServerError, "missing context metadata")
	}

	userIDStrs, ok := md["userID"]
	if !ok || len(userIDStrs) != 1 {
		return 0, errorhandler.NewAPIError(http.StatusInternalServerError, "missing userID in context metadata")
	}

	userID, err := strconv.ParseUint(userIDStrs[0], 10, 64)
	if err != nil {
		return 0, errorhandler.NewAPIError(http.StatusInternalServerError, "invalid userID in context metadata")
	}

	return userID, nil
}

func ExtractBoardIDFromContext(ctx context.Context) (uint64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errorhandler.NewAPIError(http.StatusInternalServerError, "missing context metadata")
	}

	boardIDStrs, ok := md["boardID"]
	if !ok || len(boardIDStrs) != 1 {
		return 0, errorhandler.NewAPIError(http.StatusInternalServerError, "missing boardID in context metadata")
	}

	boardID, err := strconv.ParseUint(boardIDStrs[0], 10, 64)
	if err != nil {
		return 0, errorhandler.NewAPIError(http.StatusInternalServerError, "invalid boardID in context metadata")
	}

	return boardID, nil
}
