package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// APIResponse represents a generic API response
type APIResponse struct {
	BannerID   int       `json:"banner_id,omitempty"`
	ClickCount int       `json:"click_count,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
	Message    string    `json:"message,omitempty"`
	TotalClicks int      `json:"total_clicks,omitempty"`
	FirstClick  time.Time `json:"first_click,omitempty"`
	LastClick   time.Time `json:"last_click,omitempty"`
	PeriodStart time.Time `json:"period_start,omitempty"`
	PeriodEnd   time.Time `json:"period_end,omitempty"`
	ClicksInPeriod int    `json:"clicks_in_period,omitempty"`
}

// StatsRequest represents a stats request
type StatsRequest struct {
	BannerID int       `json:"banner_id"`
	TsFrom   time.Time `json:"ts_from"`
	TsTo     time.Time `json:"ts_to"`
}

const baseURL = "http://localhost:8080"

func main() {
	fmt.Println("=== API Test Suite ===")
	fmt.Println()

	// Test health check
	fmt.Println("1. Testing health check...")
	if err := testHealth(); err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	fmt.Println("✓ Health check passed")
	fmt.Println()

	// Test counter endpoint
	fmt.Println("2. Testing counter endpoint...")
	bannerID := 1
	if err := testCounter(bannerID); err != nil {
		log.Fatalf("Counter test failed: %v", err)
	}
	fmt.Println("✓ Counter test passed")
	fmt.Println()

	// Test multiple clicks
	fmt.Println("3. Testing multiple clicks...")
	for i := 0; i < 5; i++ {
		if err := testCounter(bannerID); err != nil {
			log.Printf("Click %d failed: %v", i+1, err)
		} else {
			fmt.Printf("✓ Click %d recorded\n", i+1)
		}
		time.Sleep(100 * time.Millisecond) // Small delay between clicks
	}
	fmt.Println()

	// Test stats endpoint
	fmt.Println("4. Testing stats endpoint...")
	if err := testStats(bannerID); err != nil {
		log.Fatalf("Stats test failed: %v", err)
	}
	fmt.Println("✓ Stats test passed")
	fmt.Println()

	// Test error cases
	fmt.Println("5. Testing error cases...")
	if err := testErrorCases(); err != nil {
		log.Fatalf("Error case test failed: %v", err)
	}
	fmt.Println("✓ Error case tests passed")
	fmt.Println()

	fmt.Println("=== All tests completed successfully! ===")
}

func testHealth() error {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result["status"] != "healthy" {
		return fmt.Errorf("expected healthy status, got %v", result["status"])
	}

	return nil
}

func testCounter(bannerID int) error {
	url := fmt.Sprintf("%s/api/v1/counter/%d", baseURL, bannerID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("expected status 200, got %d: %s", resp.StatusCode, string(body))
	}

	var result APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.BannerID != bannerID {
		return fmt.Errorf("expected banner ID %d, got %d", bannerID, result.BannerID)
	}

	if result.ClickCount <= 0 {
		return fmt.Errorf("expected click count > 0, got %d", result.ClickCount)
	}

	fmt.Printf("  Recorded click for banner %d (total clicks: %d)\n", result.BannerID, result.ClickCount)
	return nil
}

func testStats(bannerID int) error {
	url := fmt.Sprintf("%s/api/v1/stats/%d", baseURL, bannerID)
	
	// Create stats request for the last hour
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	
	req := StatsRequest{
		BannerID: bannerID,
		TsFrom:   oneHourAgo,
		TsTo:     now,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("expected status 200, got %d: %s", resp.StatusCode, string(body))
	}

	var result APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.BannerID != bannerID {
		return fmt.Errorf("expected banner ID %d, got %d", bannerID, result.BannerID)
	}

	fmt.Printf("  Banner %d stats:\n", result.BannerID)
	fmt.Printf("    Total clicks: %d\n", result.TotalClicks)
	fmt.Printf("    Clicks in period: %d\n", result.ClicksInPeriod)
	if !result.FirstClick.IsZero() {
		fmt.Printf("    First click: %s\n", result.FirstClick.Format(time.RFC3339))
	}
	if !result.LastClick.IsZero() {
		fmt.Printf("    Last click: %s\n", result.LastClick.Format(time.RFC3339))
	}

	return nil
}

func testErrorCases() error {
	// Test invalid banner ID
	fmt.Println("  Testing invalid banner ID...")
	resp, err := http.Get(baseURL + "/api/v1/counter/invalid")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("expected status 400 for invalid banner ID, got %d", resp.StatusCode)
	}
	fmt.Println("    ✓ Invalid banner ID handled correctly")

	// Test non-existent banner
	fmt.Println("  Testing non-existent banner...")
	resp, err = http.Get(baseURL + "/api/v1/counter/99999")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected status 404 for non-existent banner, got %d", resp.StatusCode)
	}
	fmt.Println("    ✓ Non-existent banner handled correctly")

	// Test invalid stats request
	fmt.Println("  Testing invalid stats request...")
	url := baseURL + "/api/v1/stats/1"
	invalidJSON := `{"invalid": "json"}`
	resp, err = http.Post(url, "application/json", bytes.NewBufferString(invalidJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("expected status 400 for invalid JSON, got %d", resp.StatusCode)
	}
	fmt.Println("    ✓ Invalid stats request handled correctly")

	return nil
}
