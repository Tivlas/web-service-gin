# CRUD Web Api using Golang + PostgreSQL + Gin framework
---
**Model:**
```go
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}
```

You can test it using _curl_ (choose server and set up connection to the database in _album_repository.go_).
Endpoints (https://your-server):
- POST /albums (create)
- GET /albums (get all)
- GET /albums/id  (get by id)
- PUT /albums/id (update)
- DELETE /albums/id (delete)