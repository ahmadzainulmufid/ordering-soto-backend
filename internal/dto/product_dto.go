package dto

type CreateProductRequest struct {
	CategoryID  int64   `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required,min=2,max=150"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock" binding:"gte=0"`
}

type UpdateProductRequest struct {
	CategoryID  int64   `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required,min=2,max=150"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock" binding:"gte=0"`
	IsAvailable bool    `json:"is_available"`
}

type ProductResponse struct {
	ID           int64   `json:"id"`
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"category_name,omitempty"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	ImageURL     string  `json:"image_url"`
	Stock        int     `json:"stock"`
	IsAvailable  bool    `json:"is_available"`
}