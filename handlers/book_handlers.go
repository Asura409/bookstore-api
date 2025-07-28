package handlers

import (
	
	"bookstore-api/models"	
	"bookstore-api/repositories"
	"github.com/gofiber/fiber/v2"
)

// CreateBookHandler handles the creation of a new book.
type bookHandler struct {
	repo *repositories.BookRepository
}

func NewBookHandler(repo *repositories.BookRepository) *bookHandler {
	return &bookHandler{
		repo: repo,
	}
}

func (h *bookHandler) CreateBookHandler(c *fiber.Ctx) error {
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	
	if err := h.repo.CreateBook(book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create book"})
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

func (h *bookHandler) GetBookByIDHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	book, err := h.repo.GetBookByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	return c.JSON(book)
}



func (h *bookHandler) GetAllBooksHandler(c *fiber.Ctx) error {
	books, err := h.repo.GetAllBooks(1) // Assuming page 1 for simplicity
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve books"})
	}

	return c.JSON(books)
}

