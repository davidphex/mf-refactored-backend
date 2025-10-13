package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/davidphex/memoryframe-backend/internal/database"
	"github.com/davidphex/memoryframe-backend/internal/handlers"
	"github.com/davidphex/memoryframe-backend/internal/repository"
	"github.com/davidphex/memoryframe-backend/internal/services"
	"github.com/gin-contrib/cors"
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
	UserHandler  *handlers.UserHandler
	UserService  services.UserService
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
	userRepo := repository.NewUserRepository(db)

	albumService := services.NewAlbumService(albumRepo, photoRepo, pageRepo)
	albumHandler := handlers.NewAlbumHandler(albumService)

	pageService := services.NewPagesService(pageRepo, albumRepo)
	pageHandler := handlers.NewPageHandler(pageService)

	photoService := services.NewPhotoService(photoRepo, albumRepo)
	photoHandler := handlers.NewPhotoHandler(photoService)

	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	return &Application{
		Config:       cfg,
		DBClient:     client,
		AlbumHandler: albumHandler,
		PageHandler:  pageHandler,
		PhotoHandler: photoHandler,
		UserHandler:  userHandler,
		UserService:  userService,
	}
}

func (app *Application) setupRoutes() http.Handler {
	router := gin.Default()

	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", handlers.HealthCheck)

		v1.POST("/users/register", app.UserHandler.RegisterUser)
		v1.POST("/users/login", app.UserHandler.LoginUser)
	}

	secureGroup := v1.Group("/")
	secureGroup.Use(app.AuthMiddleware())
	{
		secureGroup.GET("/albums", app.AlbumHandler.GetAllAlbums)
		secureGroup.GET("/albums/:id", app.AlbumHandler.GetAlbumByID)
		secureGroup.POST("/albums", app.AlbumHandler.InsertAlbum)
		secureGroup.PUT("/albums/:id", app.AlbumHandler.UpdateAlbum)
		secureGroup.DELETE("/albums/:id", app.AlbumHandler.DeleteAlbum)
		secureGroup.GET("/albums/:id/pages", app.PageHandler.GetAlbumPages)
		secureGroup.GET("/albums/:id/export", app.AlbumHandler.ExportAlbum)

		secureGroup.POST("/albums/:id/photos", app.PhotoHandler.UploadPhoto)
		secureGroup.GET("/albums/:id/photos", app.PhotoHandler.GetAlbumPhotos)

		secureGroup.GET("/photos/:id", app.PhotoHandler.GetPhoto)

		secureGroup.GET("/pages/:id", app.PageHandler.GetPage)
		secureGroup.POST("/pages", app.PageHandler.InsertPage)
		secureGroup.PUT("/pages/:id", app.PageHandler.UpdatePage)
		secureGroup.DELETE("/pages/:id", app.PageHandler.DeletePage)

		secureGroup.PUT("/pages/:id/elements", app.PageHandler.UpdatePageElements)

		secureGroup.GET("/users/:id", app.UserHandler.GetUser)
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
