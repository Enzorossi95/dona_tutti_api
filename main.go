package main

import (
	"database/sql"
	"dona_tutti_api/campaign"
	"dona_tutti_api/campaigncategory"
	"dona_tutti_api/database"
	"dona_tutti_api/docs"
	"dona_tutti_api/donation"
	"dona_tutti_api/donor"
	"dona_tutti_api/migrations"
	"dona_tutti_api/organizer"
	"dona_tutti_api/rbac"
	"dona_tutti_api/user"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

// @title Dona Tutti API
// @version 1.0
// @description API for managing donations and campaigns in the Dona Tutti platform
// @host localhost:9999
// @BasePath /api
// @schemes http https
// @contact.name API Support
// @contact.email support@donatutti.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "Dona Tutti API"
	docs.SwaggerInfo.Description = "API for managing donations and campaigns in the Dona Tutti platform"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:9999"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Database setup
	db, sqlDB, err := setupDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := migrations.Up(sqlDB); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Get underlying sql.DB for proper cleanup
	sqlDB, err = db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "[${time_rfc3339}] ${status} ${method} ${uri} - ${latency_human}\n",
		CustomTimeFormat: time.RFC3339,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// API Group
	api := e.Group("/api")

	// Initialize services
	userRepo := user.NewUserRepository(db)
	userService := user.NewService(userRepo)

	campaignRepo := campaign.NewCampaignRepository(db)
	campaignService := campaign.NewService(campaignRepo)

	categoryRepo := campaigncategory.NewCategoryRepository(db)
	categoryService := campaigncategory.NewService(categoryRepo)

	organizerRepo := organizer.NewOrganizerRepository(db)
	organizerService := organizer.NewService(organizerRepo)

	donorRepo := donor.NewDonorRepository(db)
	donorService := donor.NewService(donorRepo)

	donationRepo := donation.NewDonationRepository(db)
	donationService := donation.NewService(donationRepo)

	// Initialize RBAC service
	rbacRepo := rbac.NewRepository(db)
	rbacService := rbac.NewService(rbacRepo)

	// Register routes
	user.RegisterRoutes(api, userService)
	campaign.RegisterRoutes(api, campaignService, rbacService)
	campaigncategory.RegisterRoutes(api, categoryService)
	organizer.RegisterRoutes(api, organizerService)
	donor.RegisterRoutes(api, donorService)
	donation.RegisterRoutes(api, donationService)
	rbac.RegisterRoutes(api, rbacService)

	// Start server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "9999"
	}

	fmt.Printf("\n🚀 Server running on port :%s\n", port)
	fmt.Printf("📚 Swagger documentation available at http://localhost:%s/swagger/index.html\n", port)
	fmt.Printf("🌐 API Base URL: http://localhost:%s/api\n\n", port)
	log.Fatal(e.Start(":" + port))
}

func setupDatabase() (*gorm.DB, *sql.DB, error) {
	gormDB, err := database.Connect()
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, err
	}

	return gormDB, sqlDB, nil
}
