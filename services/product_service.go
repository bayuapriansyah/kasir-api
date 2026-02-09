package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepositories
}

func NewProductService(repo *repositories.ProductRepositories) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProduct(name string) ([]models.Produk, error) {
	return s.repo.GetAllProduct(name)
}

func (s *ProductService) CreateProduct(data *models.Produk) error {
	return s.repo.CreateProduct(data)
}
func (s *ProductService) GetProdukByID(id int) (*models.Produk, error) {
	return s.repo.GetProdukByID(id)
}
func (s *ProductService) UpdateProduk(data *models.Produk) error {
	return s.repo.UpdateProduk(data)
}
func (s *ProductService) DeleteProduk(id int) error {
	return s.repo.DeleteProduk(id)
}
