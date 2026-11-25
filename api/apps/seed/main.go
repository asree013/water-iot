package main

import (
	"api/config"
	"api/modules/entities"
	"api/pkg/database"
	"log"
)

func main() {
	log.Println("[Seeder CLI] Starting seed runner...")
	cfg := config.LoadConfig()

	db, err := database.NewSQLConnection(cfg)
	if err != nil {
		log.Fatalf("[Seeder CLI] Database connection failed: %v", err)
	}

	// AutoMigrate all tables
	if err := db.AutoMigrate(
		&entities.WaterData{},
		&entities.TDMAZone{},
		&entities.TPipeAssets{},
		&entities.TLeakageIncident{},
		&entities.TDemandForecast{},
		&entities.TPumpAsset{},
		&entities.TPumpHealthScore{},
	); err != nil {
		log.Fatalf("[Seeder CLI] AutoMigrate failed: %v", err)
	}

	// Run Seed
	if err := database.SeedAll(db); err != nil {
		log.Fatalf("[Seeder CLI] Seeding failed: %v", err)
	}

	log.Println("[Seeder CLI] All mock data seeded successfully for สงขลา, เพชรบุรี, กาญจนบุรี, กทม!")
}
