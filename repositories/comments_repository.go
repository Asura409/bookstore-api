package repositories
import (
	"bookstore-api/models"
	"gorm.io/gorm"
	"time"
)
// CommentsRepository defines the interface for comment-related database operations.
type CommentsRepository struct {
	DB *gorm.DB	
}



// NewCommentsRepository creates a new instance of CommentsRepository.
func NewCommentsRepository(db *gorm.DB) *CommentsRepository {
	return &CommentsRepository{DB: db}
}

// CreateComment inserts a new comment into the database.
func (r *CommentsRepository) AddComment(comment *models.Comment) error {
	return r.DB.Create(comment).Error
}
// ReplyToComment allows a user to reply to an existing comment.
// It returns the newly created reply comment or an error if the operation fails.
func (r *CommentsRepository) ReplyToComment(comment *models.Comment, content string, userID int)  error{
	// replying is basically updating a comment with replies in the comment.reply field
	reply := &models.Comment{
		Content: content,
		CreatedAt: time.Now(),
		UserID: userID,
}
	return r.DB.Model(&comment).Association("Replies").Append(reply)
}

// LikeComment allows a user to like a comment.
// It updates the comment's like count and returns an error if the operation fails.
func (r *CommentsRepository) LikeComment(comment *models.Comment) error {
	comment.Likes++
	return r.DB.Save(comment).Error
}

// DislikeComment allows a user to dislike a comment.
// It updates the comment's dislike count and returns an error if the operation fails.
func (r *CommentsRepository) DislikeComment(comment *models.Comment) error {
	comment.Dislikes++
	return r.DB.Save(comment).Error
}

func (r *CommentsRepository) GetCommentByID(id int) (*models.Comment, error) {
	comment := &models.Comment{}
	if err := r.DB.First(comment, id).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// GetCommentsByBookID retrieves all comments associated with a specific book ID.
// It returns a slice of Comment models and an error if any occurs.
func (r *CommentsRepository) GetCommentsByBookID(bookID int) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.DB.Where("book_id = ?", bookID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}




