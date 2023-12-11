package main

import (
	"log"

	"github.com/Tivlas/web-service-gin/internal/handlers"
	"github.com/Tivlas/web-service-gin/internal/repositories"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {
	repo, err := repositories.NewPsqlAlbumRepository()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer repo.CloseDb()
	handler := handlers.NewAlbumHandler(repo)
	router := gin.Default()
	router.GET("/albums", handler.GetAll)
	router.POST("/albums", handler.Create)
	router.GET("/albums/:id", handler.GetById)
	router.PUT("/albums/:id", handler.Edit)
	router.DELETE("/albums/:id", handler.Delete)
	router.RunTLS("localhost:8080", "cert.cer", "key.pkey")
}
