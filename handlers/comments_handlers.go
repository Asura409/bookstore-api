package handlers

import (
	"bookstore-api/models"
	"bookstore-api/repositories"
	"time"
	"log"
	"github.com/gofiber/fiber/v2"
)

// CreateCommentHandler handles the creation of a new comment.
type CommentsHandler struct {
	repo repositories.CommentsRepository

}

func NewCommentHandler(repo *repositories.CommentsRepository) *CommentsHandler {
	return &CommentsHandler{
		repo: *repo,
	}
}

func (h *CommentsHandler) CreateCommentHandler(c *fiber.Ctx) error {
	comment := new(models.Comment)
	
	// Validate message
	message := c.FormValue("message")
	if message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Message is required"})
	}
	comment.Content = message
	
	// Get and validate user ID from context
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user"})
	}
	comment.UserID = user.ID
	
	// Parse and validate book ID
	book := c.Locals("book").(*models.Book)
	if book == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book"})
	}
	
	comment.BookID = book.ID
	
	comment.CreatedAt = time.Now()
   book.Comments = append(book.Comments, *comment)
	// Save comment
	
	if err := h.repo.AddComment(comment); err != nil {
		log.Printf("Failed to create comment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create comment"})
	}

	return c.Status(fiber.StatusCreated).JSON(comment)
}

func (h *CommentsHandler) LikeCommentHandler(c *fiber.Ctx) error {
	commentID, err := c.ParamsInt("commentID")
	if err != nil || commentID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment ID"})
	}
	userID, ok := c.Locals("userID").(int)
	if !ok || userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user"})
	}
	comment, errc := h.repo.GetCommentByID(commentID)
	if errc != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
	}

	if err := h.repo.LikeComment(comment); err != nil {
		log.Printf("Failed to like comment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to like comment"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *CommentsHandler) DislikeCommentHandler(c *fiber.Ctx) error {
	commentID, err := c.ParamsInt("commentID")
	if err != nil || commentID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment ID"})
	}
	userID, ok := c.Locals("userID").(int)
	if !ok || userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user"})
	}
	comment, errc := h.repo.GetCommentByID(commentID)
	if errc != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
	}

	if err := h.repo.DislikeComment(comment); err != nil {
		log.Printf("Failed to dislike comment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to dislike comment"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *CommentsHandler) GetCommentsByBookIDHandler(c *fiber.Ctx) error {
	bookID, err := c.ParamsInt("bookID")
	if err != nil || bookID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	comments, err := h.repo.GetCommentsByBookID(bookID)
	if err != nil {
		log.Printf("Failed to retrieve comments: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve comments"})
	}

	return c.JSON(comments)
}