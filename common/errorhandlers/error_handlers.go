package errorhandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type APIError struct {
	Meta   Meta         `json:"meta"`
	Errors []FieldError `json:"errors,omitempty"`
}

type Meta struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type FieldError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func (v *APIError) Error() string {
	jsonError, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("Validation error: %s", v.Meta.Message)
	}
	return string(jsonError)
}

func HandleError(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		c.JSON(http.StatusInternalServerError, NewHttpInternalError())
		return
	}

	var apiError APIError
	if err := json.Unmarshal([]byte(st.Message()), &apiError); err != nil {
		c.JSON(http.StatusInternalServerError, NewHttpInternalError())
		return
	}
	c.JSON(apiError.Meta.Status, apiError)
}

func NewAPIError(status int, message string, errors ...FieldError) *APIError {
	return &APIError{
		Meta: Meta{
			Status:  status,
			Message: message,
		},
		Errors: errors,
	}
}

func CreateGrpcErrorFromFieldErrors(fieldErrors map[string]FieldError) error {
	if len(fieldErrors) > 0 {
		errorsSlice := make([]FieldError, 0, len(fieldErrors))
		for _, err := range fieldErrors {
			errorsSlice = append(errorsSlice, err)
		}
		return status.Errorf(codes.InvalidArgument, NewAPIError(http.StatusBadRequest, "Invalid request parameters", errorsSlice...).Error())
	}
	return nil
}

func NewGrpcInternalError() error {
	return status.Errorf(codes.Internal, NewAPIError(http.StatusInternalServerError, "Internal server error").Error())
}

func NewGrpcBadRequestError(message string) error {
	return status.Errorf(codes.Internal, NewAPIError(http.StatusBadRequest, message).Error())
}

func NewGrpcNotFoundError(message string) error {
	return status.Errorf(codes.NotFound, NewAPIError(http.StatusNotFound, message).Error())
}

func NewHttpInternalError() *APIError {
	return NewAPIError(http.StatusInternalServerError, "Internal server error")
}

func NewHttpBadRequestError() *APIError {
	return NewAPIError(http.StatusBadRequest, "Bad request")
}
