package main

import (
	"log"
	"bookstore-api/middleware"
	"bookstore-api/config"
	"bookstore-api/database"
	"bookstore-api/handlers"
	"bookstore-api/services"
	"bookstore-api/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

 func main(){
	// Initialize database
	cfg := config.New() // Load configuration
	 DB := database.ConnectAndMigrate(cfg.DSN)
	 userRepo := repositories.NewUserRepository(DB)
	 bookRepo := repositories.NewBookRepository(DB)
	Crepo := repositories.NewCommentsRepository(DB)

	authService := services.NewAuthService(userRepo, cfg.JWT_SECRET)
	// Initialize email service
	emailService, err := services.NewEmailService(cfg.EMAIL_API_KEY, cfg.EMAIL_SENDER," ")
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}
	 //initialize handlers
	 userHandler := handlers.NewUserHandler(userRepo)
	 bookHandler :=  handlers.NewBookHandler(bookRepo)
	 CHandler := handlers.NewCommentHandler(Crepo)
	 AuthHandler := handlers.NewAuthHandler(authService, emailService) // Assuming AuthService is implemented in userRepo
	 // Initialize Fiber app
	 app := fiber.New()
	 // Middleware CORS and logger
	app.Use(cors.New(cors.Config{
	AllowOrigins: "https://booksbyjohn.herokuapp.com, http://localhost:3000",
	AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	// AllowCredentials: true, // Enable if using cookies/auth headers
	// MaxAge: 86400, // Cache CORS preflight for 1 day
}))
	app.Use(logger.New())


	 
	
	
	// Initialize middleware
	authMiddleware := middleware.JWTProtected(userRepo)
	
	// Public routes (no JWT required)
	user := app.Group("/user")
	{
		user.Post("/register", userHandler.CreateUserHandler)
		user.Post("/login", userHandler.LoginUserHandler)
		 // Reset password route
		user.Post("/request-password-reset", AuthHandler.RequestPasswordReset) // Request password reset route
		user.Post("/reset-password", AuthHandler.ResetPassword) // Reset password route
	}

	// Protected routes (require JWT)
	protectedUser := app.Group("/user", authMiddleware)
	{
		protectedUser.Get("/:id", userHandler.GetUserByIDHandler)
		protectedUser.Get("/:username", userHandler.GetUserByUsernameHandler)
	}

	book := app.Group("/book", authMiddleware) // All book routes protected
	{
		book.Post("/", bookHandler.CreateBookHandler)
		book.Get("/", bookHandler.GetAllBooksHandler)
		book.Get("/:id", bookHandler.GetBookByIDHandler) // Fixed typo from GetAllBooksHandler
	}

	comments := book.Group("/:bookID/comments", authMiddleware)
	{
		comments.Post("/", CHandler.CreateCommentHandler)
		comments.Get("/", CHandler.GetCommentsByBookIDHandler)
		comments.Post("/:commentID/like", CHandler.LikeCommentHandler)
		comments.Post("/:commentID/dislike", CHandler.DislikeCommentHandler)
	}

	



	 //Server
	 log.Println("Server starting on :3001")
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
 }