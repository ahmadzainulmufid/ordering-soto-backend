package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"SotoAyam/internal/dto"
	"SotoAyam/internal/models"
	"SotoAyam/internal/repository"
	"SotoAyam/internal/utils"
)

var (
	ErrCategoryNotFound = errors.New("kategori tidak ditemukan")
	ErrCategoryExists   = errors.New("kategori sudah ada")
)

type CategoryService interface {
	GetAllCategories(
		ctx context.Context,
	) ([]dto.CategoryResponse, error)

	GetCategoryByID(
		ctx context.Context,
		categoryID int64,
	) (*dto.CategoryResponse, error)

	GetCategoryByName(
		ctx context.Context,
		categoryName string,
	) (*dto.CategoryResponse, error)

	CreateCategory(
		ctx context.Context,
		request dto.CreateCategoryRequest,
	) (*dto.CategoryResponse, error)

	UpdateCategory(
		ctx context.Context,
		categoryID int64,
		request dto.UpdateCategoryRequest,
	) (*dto.CategoryResponse, error)

	DeleteCategory(
		ctx context.Context,
		categoryID int64,
	) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(
	categoryRepository repository.CategoryRepository,
) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}

func (s *categoryService) GetAllCategories(
	ctx context.Context,
) ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get categories: %w",
			err,
		)
	}

	responses := make(
		[]dto.CategoryResponse,
		0,
		len(categories),
	)

	for _, category := range categories {
		response := mapCategoryToCategoryResponse(
			&category,
		)

		responses = append(
			responses,
			response,
		)
	}

	return responses, nil
}

func (s *categoryService) GetCategoryByID(
	ctx context.Context,
	categoryID int64,
) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepository.FindByID(
		ctx,
		categoryID,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrCategoryNotFound,
		) {
			return nil, ErrCategoryNotFound
		}

		return nil, fmt.Errorf(
			"failed to get category by ID: %w",
			err,
		)
	}

	response := mapCategoryToCategoryResponse(
		category,
	)

	return &response, nil
}

func (s *categoryService) GetCategoryByName(
	ctx context.Context,
	categoryName string,
) (*dto.CategoryResponse, error) {
	categoryName = strings.TrimSpace(categoryName)

	category, err := s.categoryRepository.FindByName(
		ctx,
		categoryName,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrCategoryNotFound,
		) {
			return nil, ErrCategoryNotFound
		}

		return nil, fmt.Errorf(
			"failed to get category by name: %w",
			err,
		)
	}

	response := mapCategoryToCategoryResponse(
		category,
	)

	return &response, nil
}

func (s *categoryService) CreateCategory(
	ctx context.Context,
	request dto.CreateCategoryRequest,
) (*dto.CategoryResponse, error) {
	name := strings.TrimSpace(request.Name)

	existingCategory, err :=
		s.categoryRepository.FindByName(
			ctx,
			name,
		)

	if err == nil && existingCategory != nil {
		return nil, ErrCategoryExists
	}

	if err != nil &&
		!errors.Is(
			err,
			repository.ErrCategoryNotFound,
		) {
		return nil, fmt.Errorf(
			"failed to check existing category: %w",
			err,
		)
	}

	category := &models.Category{
		Name:     name,
		Slug:     utils.GenerateSlug(name),
		IsActive: true,
	}

	if err := s.categoryRepository.Create(
		ctx,
		category,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to create category: %w",
			err,
		)
	}

	response := mapCategoryToCategoryResponse(
		category,
	)

	return &response, nil
}

func (s *categoryService) UpdateCategory(
	ctx context.Context,
	categoryID int64,
	request dto.UpdateCategoryRequest,
) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepository.FindByID(
		ctx,
		categoryID,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrCategoryNotFound,
		) {
			return nil, ErrCategoryNotFound
		}

		return nil, fmt.Errorf(
			"failed to get category by ID: %w",
			err,
		)
	}

	name := strings.TrimSpace(request.Name)

	if !strings.EqualFold(
		name,
		category.Name,
	) {
		existingCategory, err :=
			s.categoryRepository.FindByName(
				ctx,
				name,
			)

		if err == nil &&
			existingCategory != nil &&
			existingCategory.ID != categoryID {
			return nil, ErrCategoryExists
		}

		if err != nil &&
			!errors.Is(
				err,
				repository.ErrCategoryNotFound,
			) {
			return nil, fmt.Errorf(
				"failed to check category name: %w",
				err,
			)
		}
	}

	category.Name = name
	category.Slug = utils.GenerateSlug(name)
	category.IsActive = request.IsActive

	if err := s.categoryRepository.Update(
		ctx,
		category,
	); err != nil {
		if errors.Is(
			err,
			repository.ErrCategoryNotFound,
		) {
			return nil, ErrCategoryNotFound
		}

		return nil, fmt.Errorf(
			"failed to update category: %w",
			err,
		)
	}

	response := mapCategoryToCategoryResponse(
		category,
	)

	return &response, nil
}

func (s *categoryService) DeleteCategory(
	ctx context.Context,
	categoryID int64,
) error {
	_, err := s.categoryRepository.FindByID(
		ctx,
		categoryID,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrCategoryNotFound,
		) {
			return ErrCategoryNotFound
		}

		return fmt.Errorf(
			"failed to get category by ID: %w",
			err,
		)
	}

	if err := s.categoryRepository.Delete(
		ctx,
		categoryID,
	); err != nil {
		if errors.Is(
			err,
			repository.ErrCategoryNotFound,
		) {
			return ErrCategoryNotFound
		}

		return fmt.Errorf(
			"failed to delete category: %w",
			err,
		)
	}

	return nil
}

func mapCategoryToCategoryResponse(
	category *models.Category,
) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:       category.ID,
		Name:     category.Name,
		Slug:     category.Slug,
		IsActive: category.IsActive,
	}
}


