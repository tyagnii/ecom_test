/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/tyagnii/ecom_test/cache"
	"github.com/tyagnii/ecom_test/db"
)

// cacheCmd represents the cache command
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Cache management operations",
	Long:  `Manage in-memory cache for the application.`,
}

// cacheStatsCmd represents the cache stats command
var cacheStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show cache statistics",
	Long:  `Show cache performance statistics and metrics.`,
	Run: func(cmd *cobra.Command, args []string) {
		showCacheStats()
	},
}

// cacheClearCmd represents the cache clear command
var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all cache data",
	Long:  `Clear all cached data from memory.`,
	Run: func(cmd *cobra.Command, args []string) {
		clearCache()
	},
}

// cacheWarmCmd represents the cache warm command
var cacheWarmCmd = &cobra.Command{
	Use:   "warm",
	Short: "Warm up the cache",
	Long:  `Preload frequently accessed data into cache.`,
	Run: func(cmd *cobra.Command, args []string) {
		warmCache()
	},
}

// cacheTestCmd represents the cache test command
var cacheTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test cache performance",
	Long:  `Run cache performance tests and benchmarks.`,
	Run: func(cmd *cobra.Command, args []string) {
		testCachePerformance()
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)
	cacheCmd.AddCommand(cacheStatsCmd)
	cacheCmd.AddCommand(cacheClearCmd)
	cacheCmd.AddCommand(cacheWarmCmd)
	cacheCmd.AddCommand(cacheTestCmd)
}

func getCachedRepository() (*cache.CachedRepository, error) {
	// Connect to database
	database, err := connectToDatabase()
	if err != nil {
		return nil, err
	}
	
	// Create repository and cache
	repo := db.NewRepository(database)
	cacheInstance := cache.NewInMemoryCache(cache.DefaultCleanupInterval)
	cachedRepo := cache.NewCachedRepository(repo, cacheInstance)
	
	return cachedRepo, nil
}

func showCacheStats() {
	cachedRepo, err := getCachedRepository()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer cachedRepo.GetCache().Stop()
	
	stats := cachedRepo.GetCacheStats()
	
	fmt.Printf("Cache Statistics\n")
	fmt.Printf("================\n\n")
	fmt.Printf("Size: %d items\n", stats.Size)
	fmt.Printf("Hits: %d\n", stats.Hits)
	fmt.Printf("Misses: %d\n", stats.Misses)
	fmt.Printf("Sets: %d\n", stats.Sets)
	fmt.Printf("Deletes: %d\n", stats.Deletes)
	fmt.Printf("Expirations: %d\n", stats.Expirations)
	
	if stats.Hits+stats.Misses > 0 {
		hitRate := float64(stats.Hits) / float64(stats.Hits+stats.Misses) * 100
		fmt.Printf("Hit Rate: %.2f%%\n", hitRate)
	} else {
		fmt.Printf("Hit Rate: N/A (no requests yet)\n")
	}
}

func clearCache() {
	cachedRepo, err := getCachedRepository()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer cachedRepo.GetCache().Stop()
	
	cachedRepo.ClearCache()
	fmt.Println("Cache cleared successfully!")
}

func warmCache() {
	cachedRepo, err := getCachedRepository()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer cachedRepo.GetCache().Stop()
	
	fmt.Println("Warming up cache...")
	err = cachedRepo.WarmCache()
	if err != nil {
		log.Fatalf("Failed to warm cache: %v", err)
	}
	
	stats := cachedRepo.GetCacheStats()
	fmt.Printf("Cache warmed successfully! Cached %d items.\n", stats.Size)
}

func testCachePerformance() {
	cachedRepo, err := getCachedRepository()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer cachedRepo.GetCache().Stop()
	
	fmt.Println("Running cache performance test...")
	
	// Test parameters
	testBannerID := 1
	iterations := 1000
	
	// Warm up cache first
	cachedRepo.WarmCache()
	
	// Test cache hits
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, _ = cachedRepo.GetClickStats(testBannerID)
	}
	hitDuration := time.Since(start)
	
	// Clear cache for miss test
	cachedRepo.ClearCache()
	
	// Test cache misses
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_, _ = cachedRepo.GetClickStats(testBannerID)
	}
	missDuration := time.Since(start)
	
	// Get final stats
	stats := cachedRepo.GetCacheStats()
	
	fmt.Printf("\nPerformance Test Results\n")
	fmt.Printf("=======================\n\n")
	fmt.Printf("Iterations: %d\n", iterations)
	fmt.Printf("Cache Hits Duration: %v\n", hitDuration)
	fmt.Printf("Cache Misses Duration: %v\n", missDuration)
	fmt.Printf("Hit Rate: %.2f%%\n", float64(stats.Hits)/float64(stats.Hits+stats.Misses)*100)
	
	if hitDuration > 0 && missDuration > 0 {
		speedup := float64(missDuration) / float64(hitDuration)
		fmt.Printf("Speedup: %.2fx faster with cache\n", speedup)
	}
}
