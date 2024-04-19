package services

import (
	"context"

	models "github.com/sm888sm/halten-backend/models"

	"github.com/sm888sm/halten-backend/common/errorhandler"
	pb_user "github.com/sm888sm/halten-backend/user-service/api/pb"
	"github.com/sm888sm/halten-backend/user-service/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo   repositories.UserRepository
	bcryptCost int
	pb_user.UnimplementedUserServiceServer
}

func NewUserService(userRepo repositories.UserRepository, bcryptCost int) *UserService {
	return &UserService{userRepo: userRepo, bcryptCost: bcryptCost}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb_user.CreateUserRequest) (*pb_user.CreateUserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.bcryptCost)
	if err != nil {
		return nil, errorhandler.NewGrpcInternalError()
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Fullname: req.Fullname,
		// Add other fields as necessary
	}

	repoRes, err := s.userRepo.CreateUser(&repositories.CreateUserRequest{User: user})
	if err != nil {
		return nil, err
	}

	return &pb_user.CreateUserResponse{
		UserID:   uint64(repoRes.User.ID),
		Username: repoRes.User.Email}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, req *pb_user.GetUserByIDRequest) (*pb_user.GetUserByIDResponse, error) {
	repoRes, err := s.userRepo.GetUserByID(&repositories.GetUserByIDRequest{
		UserID: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	user := &pb_user.User{
		UserID:   uint64(repoRes.User.ID),
		Username: repoRes.User.Username,
		Email:    repoRes.User.Email,
		// Add other fields as necessary
	}

	return &pb_user.GetUserByIDResponse{User: user}, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, req *pb_user.GetUserByUsernameRequest) (*pb_user.GetUserByUsernameResponse, error) {
	repoRes, err := s.userRepo.GetUserByUsername(&repositories.GetUserByUsernameRequest{
		Username: req.Username,
	})

	if err != nil {
		return nil, err
	}
	user := &pb_user.User{
		UserID:   uint64(repoRes.User.ID),
		Username: repoRes.User.Username,
		Email:    repoRes.User.Email,
		// Add other fields as necessary
	}
	return &pb_user.GetUserByUsernameResponse{User: user}, nil
}

func (s *UserService) UpdateEmail(ctx context.Context, req *pb_user.UpdateEmailRequest) (*pb_user.UpdateEmailResponse, error) {
	err := s.userRepo.UpdateEmail(&repositories.UpdateEmailRequest{
		UserID: req.UserID,
	}) // Updated
	if err != nil {
		return nil, err
	}
	return &pb_user.UpdateEmailResponse{Message: "Email changed successfully"}, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, req *pb_user.UpdatePasswordRequest) (*pb_user.UpdatePasswordResponse, error) {
	err := s.userRepo.UpdatePassword(&repositories.UpdatePasswordRequest{
		UserID:      req.UserID,
		NewPassword: req.NewPassword,
	}) // Updated
	if err != nil {
		return nil, err
	}
	return &pb_user.UpdatePasswordResponse{Message: "Password changed successfully"}, nil
}

func (s *UserService) UpdateUsername(ctx context.Context, req *pb_user.UpdateUsernameRequest) (*pb_user.UpdateUsernameResponse, error) {
	err := s.userRepo.UpdateUsername(&repositories.UpdateUsernameRequest{
		UserID:   req.UserID,
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}
	return &pb_user.UpdateUsernameResponse{Message: "Username changed successfully"}, nil
}

func (s *UserService) ConfirmNewEmail(ctx context.Context, req *pb_user.ConfirmNewEmailRequest) (*pb_user.ConfirmNewEmailResponse, error) {
	err := s.userRepo.ConfirmNewEmail(&repositories.ConfirmNewEmailRequest{
		UserID: req.UserID,
		Token:  req.Token,
	})
	if err != nil {
		return nil, err
	}
	return &pb_user.ConfirmNewEmailResponse{Message: "Email confirmed successfully"}, nil
}
