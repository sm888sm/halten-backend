package middlewares

import (
	"context"

	pb "github.com/sm888sm/halten-backend/user-service/api/pb"

	"github.com/sm888sm/halten-backend/common/constants/httpcodes"
	"github.com/sm888sm/halten-backend/common/errorhandlers"

	"google.golang.org/grpc"
)

func ValidationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	switch info.FullMethod {
	// User Service
	case "/proto.UserService/CreateUser":
		if err := validateCreateUserRequest(req.(*pb.CreateUserRequest)); err != nil {
			return nil, err
		}
	case "/proto.UserService/GetUserByID":
		if err := validateGetUserByIDRequest(req.(*pb.GetUserByIDRequest)); err != nil {
			return nil, err
		}
	case "/proto.UserService/GetUserByUsername":
		if err := validateGetUserByUsernameRequest(req.(*pb.GetUserByUsernameRequest)); err != nil {
			return nil, err
		}
	case "/proto.UserService/UpdateUsername":
		if err := validateUpdateUsernameRequest(req.(*pb.UpdateUsernameRequest)); err != nil {
			return nil, err
		}
	case "/proto.UserService/UpdateEmail":
		if err := validateUpdateEmailRequest(req.(*pb.UpdateEmailRequest)); err != nil {
			return nil, err
		}
	case "/proto.UserService/UpdatePassword":
		if err := validateUpdatePasswordRequest(req.(*pb.UpdatePasswordRequest)); err != nil {
			return nil, err
		}
	case "/proto.UserService/ConfirmEmail":
		if err := validateConfirmEmailRequest(req.(*pb.ConfirmEmailRequest)); err != nil {
			return nil, err
		}
	case "/proto.UserService/ResendConfirmationEmail":
		if err := validateResendConfirmationEmailRequest(req.(*pb.ResendConfirmationEmailRequest)); err != nil {
			return nil, err
		}

	// Auth Service
	case "/proto.AuthService/Login":
		if err := validateLoginRequest(req.(*pb.LoginRequest)); err != nil {
			return nil, err
		}
	case "/proto.AuthService/RefreshToken":
		if err := validateRefreshTokenRequest(req.(*pb.RefreshTokenRequest)); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

// User Service
func validateCreateUserRequest(req *pb.CreateUserRequest) error {
	var fieldErrors []errorhandlers.FieldError

	if req.Username == "" {
		fieldErrors = append(fieldErrors, errorhandlers.FieldError{
			Field:   "username",
			Message: "Username cannot be empty",
		})
	}

	if req.Password == "" {
		fieldErrors = append(fieldErrors, errorhandlers.FieldError{
			Field:   "password",
			Message: "Password cannot be empty",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandlers.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateGetUserByIDRequest(req *pb.GetUserByIDRequest) error {
	if req.UserID == 0 {
		return errorhandlers.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", errorhandlers.FieldError{

			Field:   "userID",
			Message: "ID cannot be empty",
		},
		)

	}

	return nil
}

func validateGetUserByUsernameRequest(req *pb.GetUserByUsernameRequest) error {
	if req.Username == "" {
		return errorhandlers.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", errorhandlers.FieldError{

			Field:   "username",
			Message: "Username cannot be empty",
		})
	}

	return nil
}

func validateUpdateUsernameRequest(req *pb.UpdateUsernameRequest) error {
	if req.Username == "" {
		return errorhandlers.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", errorhandlers.FieldError{

			Field:   "newUsername",
			Message: "New username cannot be empty",
		})
	}

	// Add more field checks as necessary

	return nil
}

func validateUpdateEmailRequest(req *pb.UpdateEmailRequest) error {
	if req.NewEmail == "" {
		return errorhandlers.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", errorhandlers.FieldError{

			Field:   "newEmail",
			Message: "New email cannot be empty",
		})
	}

	return nil
}

func validateUpdatePasswordRequest(req *pb.UpdatePasswordRequest) error {
	if req.NewPassword == "" {
		return errorhandlers.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", errorhandlers.FieldError{

			Field:   "newPassword",
			Message: "New password cannot be empty",
		})
	}

	return nil
}

func validateConfirmEmailRequest(req *pb.ConfirmEmailRequest) error {
	if req.UserID == 0 {
		return errorhandlers.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", errorhandlers.FieldError{

			Field:   "userID",
			Message: "ID cannot be empty",
		},
		)

	}

	return nil
}

func validateResendConfirmationEmailRequest(req *pb.ResendConfirmationEmailRequest) error {
	if req.Username == "" {
		return errorhandlers.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", errorhandlers.FieldError{

			Field:   "email",
			Message: "Email cannot be empty",
		})
	}

	// Add more field checks as necessary

	return nil
}

// Auth Service

func validateLoginRequest(req *pb.LoginRequest) error {
	var fieldErrors []errorhandlers.FieldError

	if req.Username == "" {
		fieldErrors = append(fieldErrors, errorhandlers.FieldError{
			Field:   "username",
			Message: "Username cannot be empty",
		})
	}

	if req.Password == "" {
		fieldErrors = append(fieldErrors, errorhandlers.FieldError{
			Field:   "password",
			Message: "Password cannot be empty",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandlers.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateRefreshTokenRequest(req *pb.RefreshTokenRequest) error {
	if req.RefreshToken == "" {
		return errorhandlers.NewAPIError(httpcodes.ErrBadRequest, "Invalid validation", errorhandlers.FieldError{

			Field:   "refreshToken",
			Message: "Refresh token cannot be empty",
		})
	}

	return nil
}
