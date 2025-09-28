package database

import (
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

func InitCloudinary() *cloudinary.Cloudinary {
	cld, _ := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if cld == nil {
		panic("Failed to initialize Cloudinary")
	}

	return cld
}
