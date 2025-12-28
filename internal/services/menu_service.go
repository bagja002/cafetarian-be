package services

import (
	"project-kelas-santai/internal/models"
	"project-kelas-santai/internal/repository"
)

type MenuService interface {
	GetAllMenus() ([]models.Menu, error)
	GetMenuByID(id string) (*models.Menu, error)
}

type menuService struct {
	repo repository.MenuRepository
}

func NewMenuService(repo repository.MenuRepository) MenuService {
	return &menuService{repo: repo}
}

func (s *menuService) GetAllMenus() ([]models.Menu, error) {
	return s.repo.FindAll()
}

func (s *menuService) GetMenuByID(id string) (*models.Menu, error) {
	return s.repo.FindByID(id)
}
