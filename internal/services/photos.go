package services

import (
	"mime/multipart"

	"github.com/davidphex/memoryframe-backend/internal/repository"
)

type PhotoService interface {
	UploadPhoto(fileHeader *multipart.FileHeader, albumId string) error
}

type photoService struct {
	photoRepo repository.PhotoRepository
	albumRepo repository.AlbumRepository
}

func NewPhotoService(photoRepo repository.PhotoRepository, albumRepo repository.AlbumRepository) PhotoService {
	return &photoService{photoRepo: photoRepo, albumRepo: albumRepo}
}

func (s *photoService) UploadPhoto(fileHeader *multipart.FileHeader, albumId string) error {

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	fileName := fileHeader.Filename

	// Call the photo repository to handle the upload
	photoId, err := s.photoRepo.Insert(file, fileName, albumId)
	if err != nil {
		return err
	}

	// Update the album to include the new photo
	err = s.albumRepo.AddPhotoToAlbum(albumId, photoId)
	if err != nil {
		return err
	}

	return nil
}
