package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

	router.Run("localhost:8080")
}
