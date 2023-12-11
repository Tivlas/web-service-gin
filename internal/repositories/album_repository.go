package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Tivlas/web-service-gin/internal/database"
	"github.com/Tivlas/web-service-gin/internal/models"
)

type PsqlAlbumRepository struct {
	db *sql.DB
}

func NewPsqlAlbumRepository() (*PsqlAlbumRepository, error) {
	db, err := database.Connect("localhost", "5432", os.Getenv("DBUSER"), os.Getenv("DBPASS"), "data-access")
	if err != nil {
		return nil, err
	}

	return &PsqlAlbumRepository{
		db: db,
	}, nil
}

func (par PsqlAlbumRepository) CloseDb() {
	par.db.Close()
}

func (par PsqlAlbumRepository) Delete(id int64) error {
	_, err := par.db.Exec("DELETE FROM album WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("Delete: %v", err)
	}
	return nil
}

func (par PsqlAlbumRepository) Edit(newAlb models.Album, id int64) error {
	_, err := par.db.Exec("UPDATE album SET title = $1, artist = $2, price = $3 WHERE id = $4", newAlb.Title, newAlb.Artist, newAlb.Price, id)
	if err != nil {
		return fmt.Errorf("Edit: %v", err)
	}
	return nil
}

func (par PsqlAlbumRepository) GetById(id int64) (models.Album, error) {
	var alb models.Album

	row := par.db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("GetById %d: no such album", id)
		}
		return alb, fmt.Errorf("GetById %d: %v", id, err)
	}
	return alb, nil
}

func (par PsqlAlbumRepository) GetAll() ([]models.Album, error) {
	var albums []models.Album
	rows, err := par.db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("GetAll: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb models.Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("GetAll: %v", err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAll: %v", err)
	}
	return albums, nil
}

func (par PsqlAlbumRepository) Create(newAlbum models.Album) (id int64, err error) {
	row := par.db.QueryRow("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", newAlbum.Title,
		newAlbum.Artist, newAlbum.Price)
	err = row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
