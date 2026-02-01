package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepositories
}

func NewCategoryService(repo *repositories.CategoryRepositories) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAllCategory() ([]models.Category, error) {
	return s.repo.GetAllCategory()
}

func (s *CategoryService) CreateCategory(data *models.Category) error {
	return s.repo.CreateCategory(data)
}

func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *CategoryService) UpdateCategory(data *models.Category) error {
	return s.repo.UpdateCategory(data)
}

func (s *CategoryService) DeleteCategory(id int) error {
	return s.repo.DeleteCategory(id)
}
 