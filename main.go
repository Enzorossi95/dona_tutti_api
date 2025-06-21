package main

import (
	"database/sql"
	"fmt"
	"log"
	"microservice_go/campaign"
	"microservice_go/campaigncategory"
	"microservice_go/database"
	"microservice_go/donation"
	"microservice_go/donor"
	"microservice_go/migrations"
	"microservice_go/organizer"
	"microservice_go/publishing"
	"microservice_go/repository"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func main() {
	// Configuración de la base de datos
	db, sqlDB, err := setupDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Ejecutar migraciones
	if err := migrations.Up(sqlDB); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Get underlying sql.DB for proper cleanup
	sqlDB, err = db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Article services
	articleRepo := repository.NewArticlesRepository(db)
	articleService := publishing.NewService(articleRepo)

	// Campaign services
	campaignRepo := campaign.NewCampaignRepository(db)
	campaignService := campaign.NewService(campaignRepo)

	// Category services
	categoryRepo := campaigncategory.NewCategoryRepository(db)
	categoryService := campaigncategory.NewService(categoryRepo)

	// Organizer services
	organizerRepo := organizer.NewOrganizerRepository(db)
	organizerService := organizer.NewService(organizerRepo)

	// Donor services
	donorRepo := donor.NewDonorRepository(db)
	donorService := donor.NewService(donorRepo)

	// Donation services
	donationRepo := donation.NewDonationRepository(db)
	donationService := donation.NewService(donationRepo)

	router := httprouter.New()

	// Register routes
	publishing.RegisterRoutes(router, articleService)
	campaign.RegisterRoutes(router, campaignService)
	campaigncategory.RegisterRoutes(router, categoryService)
	organizer.RegisterRoutes(router, organizerService)
	donor.RegisterRoutes(router, donorService)
	donation.RegisterRoutes(router, donationService)

	// Configurar CORS globalmente
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			header.Set("Access-Control-Allow-Origin", "*")
		}
		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "9999"
	}

	fmt.Printf("🚀 Hot Reload funcionando! Servidor escuchando en puerto :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func setupDatabase() (*gorm.DB, *sql.DB, error) {
	// Tu configuración actual de la base de datos con GORM
	gormDB, err := database.Connect()
	if err != nil {
		return nil, nil, err
	}

	// Obtener la conexión SQL subyacente
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, err
	}

	return gormDB, sqlDB, nil
}
