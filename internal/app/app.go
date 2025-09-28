package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/davidphex/memoryframe-backend/internal/database"
	"github.com/davidphex/memoryframe-backend/internal/handlers"
	"github.com/davidphex/memoryframe-backend/internal/repository"
	"github.com/davidphex/memoryframe-backend/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Config struct {
	Port int
	Env  string
}

type Application struct {
	Config       Config
	DBClient     *mongo.Client
	AlbumHandler *handlers.AlbumHandler
	PageHandler  *handlers.PageHandler
	PhotoHandler *handlers.PhotoHandler
}

func New(cfg Config, client *mongo.Client) *Application {

	DB_NAME := os.Getenv("DB_NAME")
	if DB_NAME == "" {
		panic("DB_NAME environment variable is not set")
	}

	db := client.Database(DB_NAME)
	cld := database.InitCloudinary()

	// Initialize repositories
	albumRepo := repository.NewAlbumRepository(db)
	pageRepo := repository.NewPagesRepository(db)
	photoRepo := repository.NewPhotoRepository(db, cld)

	albumService := services.NewAlbumService(albumRepo, photoRepo)
	albumHandler := handlers.NewAlbumHandler(albumService)

	pageService := services.NewPagesService(pageRepo, albumRepo)
	pageHandler := handlers.NewPageHandler(pageService)

	photoService := services.NewPhotoService(photoRepo, albumRepo)
	photoHandler := handlers.NewPhotoHandler(photoService)

	return &Application{
		Config:       cfg,
		DBClient:     client,
		AlbumHandler: albumHandler,
		PageHandler:  pageHandler,
		PhotoHandler: photoHandler,
	}
}

func (app *Application) setupRoutes() http.Handler {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", handlers.HealthCheck)
		v1.GET("/albums", app.AlbumHandler.GetAllAlbums)
		v1.GET("/albums/:id", app.AlbumHandler.GetAlbumByID)
		v1.POST("/albums", app.AlbumHandler.InsertAlbum)
		v1.PUT("/albums/:id", app.AlbumHandler.UpdateAlbum)
		v1.DELETE("/albums/:id", app.AlbumHandler.DeleteAlbum)
		v1.GET("/albums/:id/pages", app.PageHandler.GetAlbumPages)
		v1.GET("/albums/:id/export", app.AlbumHandler.ExportAlbum)

		v1.POST("/albums/:id/photos", app.PhotoHandler.UploadPhoto)
		v1.GET("/albums/:id/photos", app.PhotoHandler.GetAlbumPhotos)

		v1.GET("/photos/:id", app.PhotoHandler.GetPhoto)

		v1.GET("/pages/:id", app.PageHandler.GetPage)
		v1.POST("/pages", app.PageHandler.InsertPage)
		v1.PUT("/pages/:id", app.PageHandler.UpdatePage)
		v1.DELETE("/pages/:id", app.PageHandler.DeletePage)

		v1.PUT("/pages/:id/elements", app.PageHandler.UpdatePageElements)
	}

	return router
}

func (app *Application) GetHandler() http.Handler {
	return app.setupRoutes()
}

func (app *Application) Serve() error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Port),
		Handler: app.GetHandler(),
	}

	fmt.Printf("Server starting on port %d\n", app.Config.Port)
	return server.ListenAndServe()
}
