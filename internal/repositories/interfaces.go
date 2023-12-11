package repositories

import "github.com/Tivlas/web-service-gin/internal/models"

type AlbumRepository interface {
	Delete(id int64) error
	Edit(newAlb models.Album, id int64) error
	GetById(id int64) (models.Album, error)
	GetAll() ([]models.Album, error)
	Create(newAlbum models.Album) (int64, error)
	CloseDb()
}
