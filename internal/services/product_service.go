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
	ErrProductNotFound = errors.New("produk tidak ditemukan")
	ErrProductExists   = errors.New("produk sudah ada")
	ErrCategoryInvalid = errors.New("kategori tidak ditemukan")
)

type ProductService interface {
	GetAllProducts(
		ctx context.Context,
	) ([]dto.ProductResponse, error)

	GetProductByID(
		ctx context.Context,
		productID int64,
	) (*dto.ProductResponse, error)

	SearchProducts(
		ctx context.Context,
		name string,
	) ([]dto.ProductResponse, error)

	CreateProduct(
		ctx context.Context,
		request dto.CreateProductRequest,
	) (*dto.ProductResponse, error)

	UpdateProduct(
		ctx context.Context,
		productID int64,
		request dto.UpdateProductRequest,
	) (*dto.ProductResponse, error)

	DeleteProduct(
		ctx context.Context,
		productID int64,
	) error
}

type productService struct {
	productRepository  repository.ProductRepository
	categoryRepository repository.CategoryRepository
}

func NewProductService(
	productRepository repository.ProductRepository,
	categoryRepository repository.CategoryRepository,
) ProductService {
	return &productService{
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
	}
}

func (s *productService) GetAllProducts(
	ctx context.Context,
) ([]dto.ProductResponse, error) {
	products, err := s.productRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get products: %w",
			err,
		)
	}

	responses := make(
		[]dto.ProductResponse,
		0,
		len(products),
	)

	for _, product := range products {
		responses = append(
			responses,
			mapProductResponse(&product),
		)
	}

	return responses, nil
}

func (s *productService) GetProductByID(
	ctx context.Context,
	productID int64,
) (*dto.ProductResponse, error) {
	product, err := s.productRepository.FindByID(
		ctx,
		productID,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrProductNotFound,
		) {
			return nil, ErrProductNotFound
		}

		return nil, fmt.Errorf(
			"failed to get product by ID: %w",
			err,
		)
	}

	response := mapProductResponse(product)

	return &response, nil
}

func (s *productService) SearchProducts(
	ctx context.Context,
	name string,
) ([]dto.ProductResponse, error) {
	name = strings.TrimSpace(name)

	products, err := s.productRepository.FindByName(
		ctx,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to search products: %w",
			err,
		)
	}

	responses := make(
		[]dto.ProductResponse,
		0,
		len(products),
	)

	for _, product := range products {
		responses = append(
			responses,
			mapProductResponse(&product),
		)
	}

	return responses, nil
}

func (s *productService) CreateProduct(
	ctx context.Context,
	request dto.CreateProductRequest,
) (*dto.ProductResponse, error) {
	name := strings.TrimSpace(request.Name)

	category, err := s.categoryRepository.FindByID(
		ctx,
		request.CategoryID,
	)
	if err != nil {
		if errors.Is(
			err,
			repository.ErrCategoryNotFound,
		) {
			return nil, ErrCategoryInvalid
		}

		return nil, fmt.Errorf(
			"failed to check product category: %w",
			err,
		)
	}

	if !category.IsActive {
		return nil, errors.New(
			"kategori tidak aktif",
		)
	}

	existingProduct, err :=
		s.productRepository.FindExactByName(
			ctx,
			name,
		)

	if err == nil && existingProduct != nil {
		return nil, ErrProductExists
	}

	if err != nil &&
		!errors.Is(
			err,
			repository.ErrProductNotFound,
		) {
		return nil, fmt.Errorf(
			"failed to check existing product: %w",
			err,
		)
	}

	product := &models.Product{
		CategoryID:  request.CategoryID,
		Name:        name,
		Slug:        utils.GenerateSlug(name),
		Price:       request.Price,
		Stock:       request.Stock,
		IsAvailable: true,
	}

	if strings.TrimSpace(request.Description) != "" {
		product.Description = sql.NullString{
			String: strings.TrimSpace(
				request.Description,
			),
			Valid: true,
		}
	}

	if strings.TrimSpace(request.ImageURL) != "" {
		product.ImageURL = sql.NullString{
			String: strings.TrimSpace(
				request.ImageURL,
			),
			Valid: true,
		}
	}

	if err := s.productRepository.Create(
		ctx,
		product,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to create product: %w",
			err,
		)
	}

	response := mapProductResponse(product)

	return &response, nil
}

func (s *productService) UpdateProduct(
	ctx context.Context,
	productID int64,
	request dto.UpdateProductRequest,
) (*dto.ProductResponse, error) {
	product, err := s.productRepository.FindByID(
		ctx,
		productID,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrProductNotFound,
		) {
			return nil, ErrProductNotFound
		}

		return nil, fmt.Errorf(
			"failed to get product by ID: %w",
			err,
		)
	}

	category, err := s.categoryRepository.FindByID(
		ctx,
		request.CategoryID,
	)
	if err != nil {
		if errors.Is(
			err,
			repository.ErrCategoryNotFound,
		) {
			return nil, ErrCategoryInvalid
		}

		return nil, fmt.Errorf(
			"failed to check product category: %w",
			err,
		)
	}

	if !category.IsActive {
		return nil, errors.New(
			"kategori tidak aktif",
		)
	}

	name := strings.TrimSpace(request.Name)

	if !strings.EqualFold(
		name,
		product.Name,
	) {
		existingProduct, err :=
			s.productRepository.FindExactByName(
				ctx,
				name,
			)

		if err == nil &&
			existingProduct != nil &&
			existingProduct.ID != productID {
			return nil, ErrProductExists
		}

		if err != nil &&
			!errors.Is(
				err,
				repository.ErrProductNotFound,
			) {
			return nil, fmt.Errorf(
				"failed to check product name: %w",
				err,
			)
		}
	}

	product.CategoryID = request.CategoryID
	product.Name = name
	product.Slug = utils.GenerateSlug(name)
	product.Price = request.Price
	product.Stock = request.Stock
	product.IsAvailable = request.IsAvailable

	product.Description = sql.NullString{}

	if strings.TrimSpace(
		request.Description,
	) != "" {
		product.Description = sql.NullString{
			String: strings.TrimSpace(
				request.Description,
			),
			Valid: true,
		}
	}

	product.ImageURL = sql.NullString{}

	if strings.TrimSpace(
		request.ImageURL,
	) != "" {
		product.ImageURL = sql.NullString{
			String: strings.TrimSpace(
				request.ImageURL,
			),
			Valid: true,
		}
	}

	if err := s.productRepository.Update(
		ctx,
		product,
	); err != nil {
		if errors.Is(
			err,
			repository.ErrProductNotFound,
		) {
			return nil, ErrProductNotFound
		}

		return nil, fmt.Errorf(
			"failed to update product: %w",
			err,
		)
	}

	response := mapProductResponse(product)

	return &response, nil
}

func (s *productService) DeleteProduct(
	ctx context.Context,
	productID int64,
) error {
	_, err := s.productRepository.FindByID(
		ctx,
		productID,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrProductNotFound,
		) {
			return ErrProductNotFound
		}

		return fmt.Errorf(
			"failed to get product by ID: %w",
			err,
		)
	}

	if err := s.productRepository.Delete(
		ctx,
		productID,
	); err != nil {
		if errors.Is(
			err,
			repository.ErrProductNotFound,
		) {
			return ErrProductNotFound
		}

		return fmt.Errorf(
			"failed to delete product: %w",
			err,
		)
	}

	return nil
}

func mapProductResponse(
	product *models.Product,
) dto.ProductResponse {
	description := ""
	imageURL := ""

	if product.Description.Valid {
		description = product.Description.String
	}

	if product.ImageURL.Valid {
		imageURL = product.ImageURL.String
	}

	return dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: description,
		Price:       product.Price,
		ImageURL:    imageURL,
		Stock:       product.Stock,
		IsAvailable: product.IsAvailable,
	}
}