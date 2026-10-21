package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadService interface {
	UploadImage(ctx context.Context, file multipart.File) (string, error)
}

type uploadService struct {
	cld *cloudinary.Cloudinary
}

func NewUploadService(cloudinaryURL string) (UploadService, error) {
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to init cloudinary: %w", err)
	}

	return &uploadService{cld: cld}, nil
}

func (s *uploadService) UploadImage(
	ctx context.Context,
	file multipart.File,
) (string, error) {
	result, err := s.cld.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{
			Folder: "sotoayam/products",
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	if result.Error.Message != "" {
		return "", fmt.Errorf("cloudinary: %s", result.Error.Message)
	}

	if result.SecureURL == "" {
		return "", errors.New("cloudinary tidak mengembalikan URL")
	}

	return result.SecureURL, nil
}