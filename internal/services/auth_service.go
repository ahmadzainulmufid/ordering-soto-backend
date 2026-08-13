package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"SotoAyam/internal/dto"
	"SotoAyam/internal/models"
	"SotoAyam/internal/repository"
	"SotoAyam/internal/utils"
)

var (
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrInactiveUser       = errors.New("akun tidak aktif")
	ErrForbiddenRole      = errors.New("role tidak diizinkan")
	ErrEmailAlreadyExist  = errors.New("email sudah terdaftar")
)

type AuthService interface {
	Login(
		ctx context.Context,
		request dto.LoginRequest,
	) (*dto.LoginResponse, error)

	CreateUser(
		ctx context.Context,
		request dto.CreateUserRequest,
		creatorRole string,
	) (*dto.UserResponse, error)

	GetProfile(
		ctx context.Context,
		userID int64,
	) (*dto.UserResponse, error)

	GetAllUsers(ctx context.Context) ([]dto.UserResponse, error)
	UpdateUser(ctx context.Context, id int64, fullName, phone, role string) error
	DeleteUser(ctx context.Context, id int64) error
}

type authService struct {
	userRepository repository.UserRepository
	jwtManager     *utils.JWTManager
}

func NewAuthService(
	userRepository repository.UserRepository,
	jwtManager *utils.JWTManager,
) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtManager:     jwtManager,
	}
}

func (s *authService) Login(
	ctx context.Context,
	request dto.LoginRequest,
) (*dto.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(request.Email))

	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	if !user.IsActive {
		return nil, ErrInactiveUser
	}

	if err := utils.CheckPassword(
		request.Password,
		user.PasswordHash,
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, expiresIn, err := s.jwtManager.GenerateAccessToken(
		user.ID,
		user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat access token: %w", err)
	}

	if err := s.userRepository.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, fmt.Errorf(
			"gagal memperbarui waktu login terakhir: %w",
			err,
		)
	}

	response := &dto.LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User:        mapUserResponse(user),
	}

	return response, nil
}

func (s *authService) CreateUser(
	ctx context.Context,
	request dto.CreateUserRequest,
	creatorRole string,
) (*dto.UserResponse, error) {
	if creatorRole != "owner" {
		return nil, ErrForbiddenRole
	}

	email := strings.ToLower(strings.TrimSpace(request.Email))
	fullName := strings.TrimSpace(request.FullName)
	phone := strings.TrimSpace(request.Phone)
	role := strings.ToLower(strings.TrimSpace(request.Role))

	if role == "owner" {
		return nil, errors.New(
			"akun owner baru tidak dapat dibuat melalui API",
		)
	}

	if !isAllowedEmployeeRole(role) {
		return nil, errors.New("role user tidak valid")
	}

	existingUser, err := s.userRepository.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExist
	}

	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("gagal memeriksa email user: %w", err)
	}

	passwordHash, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("gagal memproses password: %w", err)
	}

	user := &models.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		IsActive:     true,
	}

	if phone != "" {
		user.Phone = sql.NullString{
			String: phone,
			Valid:  true,
		}
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("gagal membuat akun user: %w", err)
	}

	response := mapUserResponse(user)

	return &response, nil
}

func (s *authService) GetProfile(
	ctx context.Context,
	userID int64,
) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := mapUserResponse(user)

	return &response, nil
}

func mapUserResponse(user *models.User) dto.UserResponse {
	phone := ""

	if user.Phone.Valid {
		phone = user.Phone.String
	}

	return dto.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Phone:     phone,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func isAllowedEmployeeRole(role string) bool {
	switch role {
	case "admin", "cashier", "kitchen":
		return true
	default:
		return false
	}
}

func (s *authService) GetAllUsers(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.userRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil seluruh data user: %w", err)
	}

	userResponses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, mapUserResponse(&user))
	}

	return userResponses, nil
}

func (s *authService) UpdateUser(ctx context.Context, id int64, fullName, phone, role string) error {
	return s.userRepository.Update(ctx, id, fullName, phone, role)
}

func (s *authService) DeleteUser(ctx context.Context, id int64) error {
	return s.userRepository.Delete(ctx, id)
}