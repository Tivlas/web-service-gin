package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	dbname = "data-access"
)

func injectDB(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func getAlbums(db *sql.DB) ([]album, error) {
	var albums []album
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("getAlbums: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("getAlbums: %v", err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAlbums: %v", err)
	}
	return albums, nil
}

func getAlbumsHandler(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	albums, err := getAlbums(db)
	status := http.StatusOK
	if err != nil {
		log.Fatal(err)
		status = http.StatusInternalServerError
	}
	c.IndentedJSON(status, albums)
}

func addAlbum(alb album, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return fmt.Errorf("addAlbum: %v", err)
	}
	return nil
}

func addAlbumHandler(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	status := http.StatusCreated
	if err := addAlbum(newAlbum, db); err != nil {
		log.Fatal(err)
		status = http.StatusInternalServerError
	}
	c.IndentedJSON(status, newAlbum)
}

func getAlbumByID(id int64, db *sql.DB) (album, error) {
	var alb album

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("getAlbumByID %d: no such album", id)
		}
		return alb, fmt.Errorf("getAlbumByID %d: %v", id, err)
	}
	return alb, nil
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByIDHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Fatal("Error in getAlbumByIDHandler:", err)
		return
	}
	db := c.MustGet("db").(*sql.DB)
	alb, err := getAlbumByID(id, db)
	if err == nil {
		c.IndentedJSON(http.StatusOK, alb)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	}
}

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, os.Getenv("DBUSER"), os.Getenv("DBPASS"), dbname)

	db, err := sql.Open("postgres", psqlconn)
	defer db.Close()
	if err != nil {
		fmt.Println("Connection error")
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Connection error")
	}

	router := gin.Default()
	router.Use(injectDB(db))
	router.GET("/albums", getAlbumsHandler)
	router.POST("/albums", addAlbumHandler)
	router.GET("/albums/:id", getAlbumByIDHandler)
	router.Run("localhost:8080")
}
