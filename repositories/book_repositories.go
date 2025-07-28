package repositories
import (
	"gorm.io/gorm"
	"bookstore-api/models"
)
type BookRepository struct {
	DB *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (r *BookRepository) CreateBook(book *models.Book) error {
	return r.DB.Create(book).Error		}

func (r *BookRepository) GetBookByID(id int) (*models.Book, error) {
	book := &models.Book{}
	if err := r.DB.First(book, id).Error; err != nil {
		return nil, err
	}
	return book, nil
}
// GetAllBooks retrieves all books from the database in a paginated format.
// It returns a slice of Book models and an error if any occurs.
func (r *BookRepository) GetAllBooks(page int) ([]models.Book, error) {
	var books []models.Book
	limit := 10
	offset := (page - 1) * limit
	if err := r.DB.Limit(limit).Offset(offset).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil

}
func (r *BookRepository) DeleteBookByID (ID int ) error {
book := &models.Book{}
	return r.DB.Where("ID = ?", ID).Delete(book).Error
	
}
