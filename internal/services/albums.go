package services

import (
	"github.com/davidphex/memoryframe-backend/internal/models"
	"github.com/davidphex/memoryframe-backend/internal/repository"
)

type AlbumService interface {
	GetAllAlbums() (*[]models.Album, error)
	GetAlbumByID(id string) (*models.Album, error)
}

type albumService struct {
	repo repository.AlbumRepository
}

func NewAlbumService(repo repository.AlbumRepository) AlbumService {
	return &albumService{repo: repo}
}

func (s *albumService) GetAllAlbums() (*[]models.Album, error) {
	return s.repo.GetAllAlbums()
}

func (s *albumService) GetAlbumByID(id string) (*models.Album, error) {
	return s.repo.GetAlbumByID(id)
}
