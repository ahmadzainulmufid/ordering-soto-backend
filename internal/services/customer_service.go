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

const customerRole = "customer"

type CustomerAuthService interface {
	Register(
		ctx context.Context,
		request dto.CustomerRegisterRequest,
	) (*dto.CustomerResponse, error)

	Login(
		ctx context.Context,
		request dto.LoginRequest,
	) (*dto.CustomerLoginResponse, error)

	GetProfile(
		ctx context.Context,
		customerID int64,
	) (*dto.CustomerResponse, error)
}

type customerAuthService struct {
	customerRepository repository.CustomerRepository
	jwtManager         *utils.JWTManager
}

func NewCustomerAuthService(
	customerRepository repository.CustomerRepository,
	jwtManager *utils.JWTManager,
) CustomerAuthService {
	return &customerAuthService{
		customerRepository: customerRepository,
		jwtManager:         jwtManager,
	}
}

func (s *customerAuthService) Register(
	ctx context.Context,
	request dto.CustomerRegisterRequest,
) (*dto.CustomerResponse, error) {
	email := strings.ToLower(strings.TrimSpace(request.Email))
	fullName := strings.TrimSpace(request.FullName)
	phone := strings.TrimSpace(request.Phone)

	existing, err := s.customerRepository.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, ErrEmailAlreadyExist
	}
	if err != nil && !errors.Is(err, repository.ErrCustomerNotFound) {
		return nil, fmt.Errorf("gagal memeriksa email customer: %w", err)
	}

	passwordHash, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("gagal memproses password: %w", err)
	}

	customer := &models.Customer{
		FullName:     fullName,
		Email:        email,
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	if phone != "" {
		customer.Phone = sql.NullString{String: phone, Valid: true}
	}

	if err := s.customerRepository.Create(ctx, customer); err != nil {
		// Jaga-jaga kalau ada dua pendaftaran bersamaan
		if errors.Is(err, repository.ErrEmailAlreadyInUse) {
			return nil, ErrEmailAlreadyExist
		}
		return nil, fmt.Errorf("gagal membuat akun customer: %w", err)
	}

	response := mapCustomerResponse(customer)
	return &response, nil
}

func (s *customerAuthService) Login(
	ctx context.Context,
	request dto.LoginRequest,
) (*dto.CustomerLoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(request.Email))

	customer, err := s.customerRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("gagal mengambil data customer: %w", err)
	}

	if !customer.IsActive {
		return nil, ErrInactiveUser
	}

	if err := utils.CheckPassword(
		request.Password,
		customer.PasswordHash,
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Role "customer" ditetapkan di server, bukan dari input client.
	accessToken, expiresIn, err := s.jwtManager.GenerateAccessToken(
		customer.ID,
		customerRole,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat access token: %w", err)
	}

	if err := s.customerRepository.UpdateLastLogin(ctx, customer.ID); err != nil {
		return nil, fmt.Errorf("gagal memperbarui waktu login terakhir: %w", err)
	}

	return &dto.CustomerLoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		Customer:    mapCustomerResponse(customer),
	}, nil
}

func (s *customerAuthService) GetProfile(
	ctx context.Context,
	customerID int64,
) (*dto.CustomerResponse, error) {
	customer, err := s.customerRepository.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	response := mapCustomerResponse(customer)
	return &response, nil
}

func mapCustomerResponse(customer *models.Customer) dto.CustomerResponse {
	phone := ""
	if customer.Phone.Valid {
		phone = customer.Phone.String
	}

	return dto.CustomerResponse{
		ID:        customer.ID,
		FullName:  customer.FullName,
		Email:     customer.Email,
		Phone:     phone,
		IsActive:  customer.IsActive,
		CreatedAt: customer.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}