package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Debug: Print the template files that are available
	files, err := filepath.Glob("./templates/**/*.html")
	if err != nil {
		log.Printf("Error finding templates: %v", err)
	} else {
		log.Printf("Found templates: %v", files)
	}

	//setup the templates engine
	engine := html.New("./templates", ".html")
	engine.AddFunc("now", func() string {
		return time.Now().Format("Jan 2")
	})

	engine.Reload(true)

	engine.Load()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "Portfolio Website",
		Views:                 engine,
		DisableStartupMessage: true,
		IdleTimeout:           5 * time.Second,
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          5 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${ip} ${status} ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Static files
	app.Static("/", "./static")

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("layouts/base", fiber.Map{
			"title": "Home",
		})
	})

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	// Start server
	log.Printf("Server starting on port %s\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
