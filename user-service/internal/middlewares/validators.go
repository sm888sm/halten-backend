package middlewares

import (
	"context"

	pb "github.com/sm888sm/halten-backend/user-service/api/pb"

	"github.com/sm888sm/halten-backend/common/errorhandler"

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
	case "/proto.UserService/ConfirmNewEmail":
		if err := validateConfirmNewEmailRequest(req.(*pb.ConfirmNewEmailRequest)); err != nil {
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
	var fieldErrors []errorhandler.FieldError

	if req.Username == "" {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "username",
			Message: "Username cannot be empty",
		})
	}

	if req.Password == "" {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "password",
			Message: "Password cannot be empty",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateGetUserByIDRequest(req *pb.GetUserByIDRequest) error {
	if req.Id == 0 {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", errorhandler.FieldError{

			Field:   "id",
			Message: "ID cannot be empty",
		},
		)

	}

	return nil
}

func validateGetUserByUsernameRequest(req *pb.GetUserByUsernameRequest) error {
	if req.Username == "" {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", errorhandler.FieldError{

			Field:   "username",
			Message: "Username cannot be empty",
		})
	}

	return nil
}

func validateUpdateUsernameRequest(req *pb.UpdateUsernameRequest) error {
	if req.NewUsername == "" {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", errorhandler.FieldError{

			Field:   "newUsername",
			Message: "New username cannot be empty",
		})
	}

	// Add more field checks as necessary

	return nil
}

func validateUpdateEmailRequest(req *pb.UpdateEmailRequest) error {
	if req.NewEmail == "" {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", errorhandler.FieldError{

			Field:   "newEmail",
			Message: "New email cannot be empty",
		})
	}

	return nil
}

func validateUpdatePasswordRequest(req *pb.UpdatePasswordRequest) error {
	if req.NewPassword == "" {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", errorhandler.FieldError{

			Field:   "newPassword",
			Message: "New password cannot be empty",
		})
	}

	return nil
}

func validateConfirmNewEmailRequest(req *pb.ConfirmNewEmailRequest) error {
	if req.Username == "" {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", errorhandler.FieldError{

			Field:   "username",
			Message: "username cannot be empty",
		})
	}

	return nil
}

func validateResendConfirmationEmailRequest(req *pb.ResendConfirmationEmailRequest) error {
	if req.Username == "" {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", errorhandler.FieldError{

			Field:   "email",
			Message: "Email cannot be empty",
		})
	}

	// Add more field checks as necessary

	return nil
}

// Auth Service

func validateLoginRequest(req *pb.LoginRequest) error {
	var fieldErrors []errorhandler.FieldError

	if req.Username == "" {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "username",
			Message: "Username cannot be empty",
		})
	}

	if req.Password == "" {
		fieldErrors = append(fieldErrors, errorhandler.FieldError{
			Field:   "password",
			Message: "Password cannot be empty",
		})
	}

	if len(fieldErrors) > 0 {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", fieldErrors...)
	}

	return nil
}

func validateRefreshTokenRequest(req *pb.RefreshTokenRequest) error {
	if req.RefreshToken == "" {
		return errorhandler.NewAPIError(errorhandler.ErrBadRequest, "Invalid validation", errorhandler.FieldError{

			Field:   "refreshToken",
			Message: "Refresh token cannot be empty",
		})
	}

	return nil
}
