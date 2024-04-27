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
		return 0, errorhandler.NewGrpcInternalError()
	}

	boardIDStrs, ok := md["boardID"]
	if !ok || len(boardIDStrs) != 1 {
		return 0, errorhandler.NewGrpcInternalError()
	}

	boardID, err := strconv.ParseUint(boardIDStrs[0], 10, 64)
	if err != nil {
		return 0, errorhandler.NewGrpcInternalError()
	}

	return boardID, nil
}

func ExtractListIDFromContext(ctx context.Context) (uint64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errorhandler.NewGrpcInternalError()
	}

	listIDStrs, ok := md["listID"]
	if !ok || len(listIDStrs) != 1 {
		return 0, errorhandler.NewGrpcInternalError()
	}

	cardID, err := strconv.ParseUint(listIDStrs[0], 10, 64)
	if err != nil {
		return 0, errorhandler.NewGrpcInternalError()
	}

	return cardID, nil
}

func ExtractCardIDFromContext(ctx context.Context) (uint64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errorhandler.NewGrpcInternalError()
	}

	cardIDStrs, ok := md["cardID"]
	if !ok || len(cardIDStrs) != 1 {
		return 0, errorhandler.NewGrpcInternalError()
	}

	cardID, err := strconv.ParseUint(cardIDStrs[0], 10, 64)
	if err != nil {
		return 0, errorhandler.NewGrpcInternalError()
	}

	return cardID, nil
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
