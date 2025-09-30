package services

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/davidphex/memoryframe-backend/internal/models"
	"github.com/davidphex/memoryframe-backend/internal/repository"
	"github.com/davidphex/memoryframe-backend/internal/utils"
)

type AlbumService interface {
	GetAllAlbums() (*[]models.Album, error)
	GetAlbumByID(id string) (*models.Album, error)
	InsertAlbum(album *models.Album) error
	UpdateAlbum(album *models.Album) error
	DeleteAlbum(id string) error

	GeneratePDF(albumID string) ([]byte, error)
}

type albumService struct {
	albumRepo repository.AlbumRepository
	photoRepo repository.PhotoRepository
	pageRepo  repository.PagesRepository
}

func NewAlbumService(
	albumRepo repository.AlbumRepository,
	photoRepo repository.PhotoRepository,
	pageRepo repository.PagesRepository,
) AlbumService {
	return &albumService{
		albumRepo: albumRepo,
		photoRepo: photoRepo,
		pageRepo:  pageRepo,
	}
}

func (s *albumService) GetAllAlbums() (*[]models.Album, error) {
	return s.albumRepo.GetAll()
}

func (s *albumService) GetAlbumByID(id string) (*models.Album, error) {
	return s.albumRepo.GetById(id)
}

func (s *albumService) InsertAlbum(album *models.Album) error {
	return s.albumRepo.Insert(album)
}

func (s *albumService) UpdateAlbum(album *models.Album) error {
	return s.albumRepo.Update(album)
}

func (s *albumService) DeleteAlbum(id string) error {
	return s.albumRepo.Delete(id)
}

func (s *albumService) GeneratePDF(albumId string) ([]byte, error) {
	album, err := s.albumRepo.GetById(albumId)
	if err != nil {
		return nil, err
	}
	if album == nil {
		return nil, nil // Album not found
	}

	// Get pages for the album
	pages, err := s.pageRepo.GetByAlbumId(albumId)
	if err != nil {
		return nil, err
	}

	fmt.Println("Pages fetched for album:", len(*pages))

	// Prepare the HTML template
	tmpl, err := template.ParseFiles("internal/templates/index.html")
	if err != nil {
		return nil, err
	}

	// Render the template with album data
	var htmlBuffer bytes.Buffer
	err = tmpl.Execute(&htmlBuffer, struct {
		Title       string
		Description string
		Pages       *[]models.AlbumPage
	}{
		Title:       album.Title,
		Description: album.Description,
		Pages:       pages,
	})
	if err != nil {
		return nil, err
	}

	// Write the HTML buffer to a temp file
	htmlFilePath := "/tmp/index.html"
	err = os.WriteFile(htmlFilePath, htmlBuffer.Bytes(), 0644)
	if err != nil {
		return nil, err
	}

	// Write CSS to file
	cssFilePath := "/tmp/style.css"
	cssContent, err := os.ReadFile("internal/templates/style.css")
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(cssFilePath, cssContent, 0644)
	if err != nil {
		return nil, err
	}

	// Generate PDF using Gotenberg
	pdfBytes, err := utils.GenerateGotenbergPDF(htmlFilePath, cssFilePath)
	if err != nil {
		return nil, err
	}

	return pdfBytes, nil

}
