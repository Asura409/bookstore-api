package models

import "time"

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
	Comments []Comment `json:"comments"` // comments on the book
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"` // store as hash in production
	Email     string `json:"email"`
	Role      string `json:"role"` // "admin" or "user"
	CreatedAt string `json:"created_at"`
	Books    []Book `json:"books"` // books owned by the user
}


type Comment struct {
	ID        int       `json:"id"`
	BookID    int       `json:"book_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Replies   []*Comment `json:"replies"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
}
// reply tp a comment is an array of comments




