package app

import (
	"fmt"
	"net/http"
	"os"

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
}

func New(cfg Config, client *mongo.Client) *Application {

	DB_NAME := os.Getenv("DB_NAME")
	if DB_NAME == "" {
		panic("DB_NAME environment variable is not set")
	}

	db := client.Database(DB_NAME)

	albumRepo := repository.NewAlbumRepository(db)
	albumService := services.NewAlbumService(albumRepo)
	albumHandler := handlers.NewAlbumHandler(albumService)

	return &Application{
		Config:       cfg,
		DBClient:     client,
		AlbumHandler: albumHandler,
	}
}

func (app *Application) setupRoutes() http.Handler {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", handlers.HealthCheck)
		v1.GET("/albums", app.AlbumHandler.GetAllAlbums)
		v1.GET("/albums/:id", app.AlbumHandler.GetAlbumByID)
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
