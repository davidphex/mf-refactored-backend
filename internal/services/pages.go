package services

import (
	"github.com/davidphex/memoryframe-backend/internal/models"
	"github.com/davidphex/memoryframe-backend/internal/repository"
)

type PagesService interface {
	// Basic CRUD operations
	GetPage(pageId string) (*models.AlbumPage, error)
	InsertPage(page *models.AlbumPage) error
	UpdatePage(page *models.AlbumPage) error
	DeletePage(pageId string) error

	// Custom query
	GetAlbumPages(albumId string) (*[]models.AlbumPage, error)
}

type pagesService struct {
	pagesRepo  repository.PagesRepository
	albumsRepo repository.AlbumRepository
}

func NewPagesService(pagesRepo repository.PagesRepository, albumsRepo repository.AlbumRepository) PagesService {
	return &pagesService{pagesRepo: pagesRepo, albumsRepo: albumsRepo}
}

func (s *pagesService) GetPage(pageId string) (*models.AlbumPage, error) {
	return s.pagesRepo.Get(pageId)
}

func (s *pagesService) InsertPage(page *models.AlbumPage) error {

	err := s.pagesRepo.Insert(page)
	if err != nil {
		return err
	}

	// Also update the corresponding album to include this page ID
	album, err := s.albumsRepo.GetById(page.AlbumID.Hex())
	if err != nil {
		return err
	}
	if album == nil {
		return nil
	}

	album.PagesID = append(album.PagesID, page.ID)
	err = s.albumsRepo.Update(album)

	return err
}

func (s *pagesService) UpdatePage(page *models.AlbumPage) error {
	return s.pagesRepo.Update(page)
}

func (s *pagesService) DeletePage(pageId string) error {
	return s.pagesRepo.Delete(pageId)
}

func (s *pagesService) GetAlbumPages(albumId string) (*[]models.AlbumPage, error) {
	return s.pagesRepo.GetByAlbumId(albumId)
}
