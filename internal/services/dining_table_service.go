package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"SotoAyam/internal/dto"
	"SotoAyam/internal/models"
	"SotoAyam/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrDiningTableNotFound = errors.New("meja makan tidak ditemukan")
	ErrDiningTableExists    = errors.New("meja makan sudah ada")
)

type DiningTableService interface {
	GetAllDiningTables(
		ctx context.Context,
	) ([]dto.DiningTableResponse, error)

	GetDiningTableByID(
		ctx context.Context,
		diningTableID int64,
	) (*dto.DiningTableResponse, error)

	GetDiningTableByQRToken(
		ctx context.Context,
		qrToken string,
	) (*dto.DiningTableResponse, error)

	CreateDiningTable(
		ctx context.Context,
		request dto.CreateDiningTableRequest,
	) (*dto.DiningTableResponse, error)

	UpdateDiningTable(
		ctx context.Context,
		diningTableID int64,
		request dto.UpdateDiningTableRequest,
	) (*dto.DiningTableResponse, error)

	DeleteDiningTable(
		ctx context.Context,
		diningTableID int64,
	) error
}

type diningTableService struct {
	diningTableRepository repository.DiningTableRepository
}

func NewDiningTableService(
	diningTableRepository repository.DiningTableRepository,
) DiningTableService {
	return &diningTableService{
		diningTableRepository: diningTableRepository,
	}
}

func (s *diningTableService) GetAllDiningTables(
	ctx context.Context,
) ([]dto.DiningTableResponse, error) {
	diningTables, err := s.diningTableRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get all dining tables: %w",
			err,
		)
	}

	responses := make([]dto.DiningTableResponse, 0, len(diningTables))

	for _, diningTable := range diningTables {
		responses = append(
			responses,
			mapDiningTableResponse(&diningTable),
		)
	}

	return responses, nil
}

func (s *diningTableService) GetDiningTableByID(
	ctx context.Context,
	diningTableID int64,
) (*dto.DiningTableResponse, error) {
	diningTable, err := s.diningTableRepository.FindByID(
		ctx,
		diningTableID,
	)

	if err != nil {
		if errors.Is(err, repository.ErrDiningTableNotFound) {
			return nil, ErrDiningTableNotFound
		}

		return nil, fmt.Errorf(
			"failed to get dining table by ID: %w",
			err,
		)
	}

	response := mapDiningTableResponse(diningTable)

	return &response, nil
}

func (s *diningTableService) GetDiningTableByQRToken(
	ctx context.Context,
	qrToken string,
) (*dto.DiningTableResponse, error) {
	qrToken = strings.TrimSpace(qrToken)

	diningTable, err := s.diningTableRepository.FindByQRToken(
		ctx,
		qrToken,
	)

	if err != nil {
		if errors.Is(err, repository.ErrDiningTableNotFound) {
			return nil, ErrDiningTableNotFound
		}

		return nil, fmt.Errorf(
			"failed to get dining table by QR token: %w",
			err,
		)
	}

	if !diningTable.IsActive {
		return nil, ErrDiningTableNotFound
	}

	response := mapDiningTableResponse(diningTable)

	return &response, nil
}

func (s *diningTableService) CreateDiningTable(
	ctx context.Context,
	request dto.CreateDiningTableRequest,
) (*dto.DiningTableResponse, error) {
	tableNumber := strings.TrimSpace(request.TableNumber)

	existingDiningTable, err :=
		s.diningTableRepository.FindByTableNumber(
			ctx,
			tableNumber,
		)

	if err == nil && existingDiningTable != nil {
		return nil, ErrDiningTableExists
	}

	if err != nil &&
		!errors.Is(err, repository.ErrDiningTableNotFound) {
		return nil, fmt.Errorf(
			"failed to check existing dining table: %w",
			err,
		)
	}

	diningTable := &models.DiningTable{
		TableNumber: tableNumber,
		QRToken:     uuid.NewString(),
		IsActive:    true,
	}

	if err := s.diningTableRepository.Create(
		ctx,
		diningTable,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to create dining table: %w",
			err,
		)
	}

	response := mapDiningTableResponse(diningTable)

	return &response, nil
}

func (s *diningTableService) UpdateDiningTable(
	ctx context.Context,
	diningTableID int64,
	request dto.UpdateDiningTableRequest,
) (*dto.DiningTableResponse, error) {
	diningTable, err := s.diningTableRepository.FindByID(
		ctx,
		diningTableID,
	)

	if err != nil {
		if errors.Is(err, repository.ErrDiningTableNotFound) {
			return nil, ErrDiningTableNotFound
		}

		return nil, fmt.Errorf(
			"failed to get dining table by ID: %w",
			err,
		)
	}

	tableNumber := strings.TrimSpace(request.TableNumber)

	if tableNumber != diningTable.TableNumber {
		existingDiningTable, err :=
			s.diningTableRepository.FindByTableNumber(
				ctx,
				tableNumber,
			)

		if err == nil &&
			existingDiningTable != nil &&
			existingDiningTable.ID != diningTableID {
			return nil, ErrDiningTableExists
		}

		if err != nil &&
			!errors.Is(
				err,
				repository.ErrDiningTableNotFound,
			) {
			return nil, fmt.Errorf(
				"failed to check table number: %w",
				err,
			)
		}
	}

	diningTable.TableNumber = tableNumber
	diningTable.IsActive = request.IsActive

	if err := s.diningTableRepository.Update(
		ctx,
		diningTable,
	); err != nil {
		if errors.Is(err, repository.ErrDiningTableNotFound) {
			return nil, ErrDiningTableNotFound
		}

		return nil, fmt.Errorf(
			"failed to update dining table: %w",
			err,
		)
	}

	response := mapDiningTableResponse(diningTable)

	return &response, nil
}

func (s *diningTableService) DeleteDiningTable(
	ctx context.Context,
	diningTableID int64,
) error {
	_, err := s.diningTableRepository.FindByID(
		ctx,
		diningTableID,
	)

	if err != nil {
		if errors.Is(err, repository.ErrDiningTableNotFound) {
			return ErrDiningTableNotFound
		}

		return fmt.Errorf(
			"failed to get dining table by ID: %w",
			err,
		)
	}

	if err := s.diningTableRepository.Delete(
		ctx,
		diningTableID,
	); err != nil {
		if errors.Is(err, repository.ErrDiningTableNotFound) {
			return ErrDiningTableNotFound
		}

		return fmt.Errorf(
			"failed to delete dining table: %w",
			err,
		)
	}

	return nil
}

func mapDiningTableResponse(
	diningTable *models.DiningTable,
) dto.DiningTableResponse {
	return dto.DiningTableResponse{
		ID:          diningTable.ID,
		TableNumber: diningTable.TableNumber,
		QRToken:     diningTable.QRToken,
		IsActive:    diningTable.IsActive,
	}
}