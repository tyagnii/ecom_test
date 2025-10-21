package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/tyagnii/ecom_test/app"
	"github.com/tyagnii/ecom_test/db"
	_ "github.com/lib/pq"
)

func main() {
	// Database configuration
	config := &db.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "your_password",
		DBName:   "ecom_test",
		SSLMode:  "disable",
	}

	// Connect to database
	database, err := db.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create repository and services
	repo := db.NewRepository(database)
	service := app.NewService(repo)
	bannerService := app.NewBannerService(service)
	clickService := app.NewClickService(service)
	analyticsService := app.NewAnalyticsService(service)

	fmt.Println("=== E-commerce CRUD Example ===\n")

	// Create banners
	fmt.Println("1. Creating banners...")
	banner1, err := bannerService.CreateBanner("Summer Sale Banner")
	if err != nil {
		log.Fatalf("Failed to create banner 1: %v", err)
	}
	fmt.Printf("Created banner: %s (ID: %d)\n", banner1.Name, banner1.ID)

	banner2, err := bannerService.CreateBanner("Black Friday Banner")
	if err != nil {
		log.Fatalf("Failed to create banner 2: %v", err)
	}
	fmt.Printf("Created banner: %s (ID: %d)\n", banner2.Name, banner2.ID)

	// Record some clicks
	fmt.Println("\n2. Recording clicks...")
	clicks := []struct {
		bannerID int
		delay    time.Duration
	}{
		{banner1.ID, 0},
		{banner1.ID, 1 * time.Second},
		{banner2.ID, 2 * time.Second},
		{banner1.ID, 3 * time.Second},
		{banner2.ID, 4 * time.Second},
	}

	for i, click := range clicks {
		clickTime := time.Now().Add(click.delay)
		recordedClick, err := clickService.RecordClick(click.bannerID, clickTime)
		if err != nil {
			log.Fatalf("Failed to record click %d: %v", i+1, err)
		}
		fmt.Printf("Recorded click %d for banner %d at %s\n", 
			recordedClick.ID, recordedClick.BannerID, recordedClick.Timestamp.Format(time.RFC3339))
	}

	// List all banners
	fmt.Println("\n3. Listing all banners...")
	banners, err := bannerService.GetAllBanners()
	if err != nil {
		log.Fatalf("Failed to get banners: %v", err)
	}
	for _, banner := range banners {
		fmt.Printf("Banner %d: %s (created: %s)\n", 
			banner.ID, banner.Name, banner.CreatedAt.Format(time.RFC3339))
	}

	// Get clicks for banner 1
	fmt.Println("\n4. Getting clicks for banner 1...")
	clicksForBanner1, err := clickService.GetClicksForBanner(banner1.ID)
	if err != nil {
		log.Fatalf("Failed to get clicks for banner 1: %v", err)
	}
	fmt.Printf("Found %d clicks for banner %d:\n", len(clicksForBanner1), banner1.ID)
	for _, click := range clicksForBanner1 {
		fmt.Printf("  Click %d at %s\n", click.ID, click.Timestamp.Format(time.RFC3339))
	}

	// Get click statistics
	fmt.Println("\n5. Getting click statistics...")
	stats1, err := clickService.GetClickStats(banner1.ID)
	if err != nil {
		log.Fatalf("Failed to get stats for banner 1: %v", err)
	}
	fmt.Printf("Banner %d stats: %d total clicks\n", banner1.ID, stats1.TotalClicks)
	if stats1.TotalClicks > 0 {
		fmt.Printf("  First click: %s\n", stats1.FirstClick.Format(time.RFC3339))
		fmt.Printf("  Last click: %s\n", stats1.LastClick.Format(time.RFC3339))
	}

	stats2, err := clickService.GetClickStats(banner2.ID)
	if err != nil {
		log.Fatalf("Failed to get stats for banner 2: %v", err)
	}
	fmt.Printf("Banner %d stats: %d total clicks\n", banner2.ID, stats2.TotalClicks)
	if stats2.TotalClicks > 0 {
		fmt.Printf("  First click: %s\n", stats2.FirstClick.Format(time.RFC3339))
		fmt.Printf("  Last click: %s\n", stats2.LastClick.Format(time.RFC3339))
	}

	// Update a banner
	fmt.Println("\n6. Updating banner...")
	updatedBanner, err := bannerService.UpdateBanner(banner1.ID, "Updated Summer Sale Banner")
	if err != nil {
		log.Fatalf("Failed to update banner: %v", err)
	}
	fmt.Printf("Updated banner: %s (updated: %s)\n", 
		updatedBanner.Name, updatedBanner.UpdatedAt.Format(time.RFC3339))

	// Get banner performance analytics
	fmt.Println("\n7. Banner performance analytics...")
	performances, err := analyticsService.GetBannerPerformance()
	if err != nil {
		log.Fatalf("Failed to get banner performance: %v", err)
	}
	for _, perf := range performances {
		fmt.Printf("Banner '%s': %d clicks\n", perf.Banner.Name, perf.TotalClicks)
	}

	// Get clicks in date range
	fmt.Println("\n8. Getting clicks in date range...")
	start := time.Now().Add(-1 * time.Hour)
	end := time.Now().Add(1 * time.Hour)
	clicksInRange, err := clickService.GetClicksInDateRange(start, end)
	if err != nil {
		log.Fatalf("Failed to get clicks in date range: %v", err)
	}
	fmt.Printf("Found %d clicks between %s and %s\n", 
		len(clicksInRange), start.Format(time.RFC3339), end.Format(time.RFC3339))

	fmt.Println("\n=== Example completed successfully! ===")
}
