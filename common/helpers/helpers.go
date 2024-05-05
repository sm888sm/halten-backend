package helpers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/sm888sm/halten-backend/common/errorhandlers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func extractIDFromContext(ctx context.Context, key string, errMsg string) (uint64, error) {
	errInvalid := status.Errorf(codes.InvalidArgument, errorhandlers.NewAPIError(http.StatusBadRequest, errMsg).Error())

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errInvalid
	}

	idStrs, ok := md[key]
	if !ok || len(idStrs) != 1 {
		return 0, errInvalid
	}

	id, err := strconv.ParseUint(idStrs[0], 10, 64)
	if err != nil {
		return 0, errInvalid
	}

	return id, nil
}

func ExtractUserIDFromContext(ctx context.Context) (uint64, error) {
	return extractIDFromContext(ctx, "userID", "Invalid userID")
}

func ExtractBoardIDFromContext(ctx context.Context) (uint64, error) {
	return extractIDFromContext(ctx, "boardID", "Invalid boardID")
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
