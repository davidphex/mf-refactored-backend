package services

import (
	"github.com/davidphex/memoryframe-backend/internal/models"
	"github.com/davidphex/memoryframe-backend/internal/repository"
)

type AlbumService interface {
	GetAllAlbums() (*[]models.Album, error)
	GetAlbumByID(id string) (*models.Album, error)
	InsertAlbum(album *models.Album) error
	UpdateAlbum(album *models.Album) error
	DeleteAlbum(id string) error
}

type albumService struct {
	repo repository.AlbumRepository
}

func NewAlbumService(repo repository.AlbumRepository) AlbumService {
	return &albumService{repo: repo}
}

func (s *albumService) GetAllAlbums() (*[]models.Album, error) {
	return s.repo.GetAll()
}

func (s *albumService) GetAlbumByID(id string) (*models.Album, error) {
	return s.repo.GetById(id)
}

func (s *albumService) InsertAlbum(album *models.Album) error {
	return s.repo.Insert(album)
}

func (s *albumService) UpdateAlbum(album *models.Album) error {
	return s.repo.Update(album)
}

func (s *albumService) DeleteAlbum(id string) error {
	return s.repo.Delete(id)
}
