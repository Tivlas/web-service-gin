package handlers

import (
	"net/http"
	"strconv"

	"github.com/Tivlas/web-service-gin/internal/models"
	"github.com/Tivlas/web-service-gin/internal/repositories"
	"github.com/gin-gonic/gin"
)

type AlbumHandler struct {
	repo repositories.AlbumRepository
}

func NewAlbumHandler(repo repositories.AlbumRepository) *AlbumHandler {
	return &AlbumHandler{
		repo: repo,
	}
}

func (nah AlbumHandler) Create(c *gin.Context) {
	var newAlbum models.Album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	id, err := nah.repo.Create(newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.IndentedJSON(http.StatusCreated, id)
	}
}

func (nah AlbumHandler) Edit(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}
	var newAlbum models.Album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	if err = nah.repo.Edit(newAlbum, id); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, newAlbum)
	} else {
		c.IndentedJSON(http.StatusOK, newAlbum)
	}
}

func (nah AlbumHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	if err = nah.repo.Delete(id); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
}

func (nah AlbumHandler) GetById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	alb, err := nah.repo.GetById(id)
	if err == nil {
		c.IndentedJSON(http.StatusOK, alb)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	}
}

func (nah AlbumHandler) GetAll(c *gin.Context) {
	albums, err := nah.repo.GetAll()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}
