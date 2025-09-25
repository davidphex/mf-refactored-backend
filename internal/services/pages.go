package services

import (
	"github.com/davidphex/memoryframe-backend/internal/models"
	"github.com/davidphex/memoryframe-backend/internal/repository"
)

type PagesService interface {
	GetPages(albumId string) (*[]models.AlbumPage, error)
}

type pagesService struct {
	repo repository.PagesRepository
}

func NewPagesService(repo repository.PagesRepository) PagesService {
	return &pagesService{repo: repo}
}

func (s *pagesService) GetPages(albumId string) (*[]models.AlbumPage, error) {
	return s.repo.GetPages(albumId)
}
