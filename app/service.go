package app

import (
	"fmt"
	"time"

	"github.com/tyagnii/ecom_test/db"
	"github.com/tyagnii/ecom_test/dto"
)

// Service provides business logic layer
type Service struct {
	repo *db.Repository
}

// NewService creates a new service instance
func NewService(repo *db.Repository) *Service {
	return &Service{repo: repo}
}

// Repo returns the repository (for internal use)
func (s *Service) Repo() interface{} {
	return s.repo
}

// BannerService provides banner business logic
type BannerService struct {
	*Service
}

// NewBannerService creates a new banner service
func NewBannerService(service *Service) *BannerService {
	return &BannerService{Service: service}
}

// CreateBanner creates a new banner with validation
func (s *BannerService) CreateBanner(name string) (*dto.Banner, error) {
	// Validate input
	if name == "" {
		return nil, fmt.Errorf("banner name cannot be empty")
	}
	
	if len(name) > 255 {
		return nil, fmt.Errorf("banner name cannot exceed 255 characters")
	}
	
	// Check if banner with same name already exists
	existingBanner, err := s.repo.GetBannerByName(name)
	if err == nil && existingBanner != nil {
		return nil, fmt.Errorf("banner with name '%s' already exists", name)
	}
	
	// Create new banner
	banner := &dto.Banner{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	if err := s.repo.CreateBanner(banner); err != nil {
		return nil, fmt.Errorf("failed to create banner: %w", err)
	}
	
	return banner, nil
}

// GetBanner retrieves a banner by ID
func (s *BannerService) GetBanner(id int) (*dto.Banner, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid banner ID: %d", id)
	}
	
	banner, err := s.repo.GetBannerByID(id)
	if err != nil {
		return nil, err
	}
	
	return banner, nil
}

// GetAllBanners retrieves all banners
func (s *BannerService) GetAllBanners() ([]*dto.Banner, error) {
	return s.repo.GetAllBanners()
}

// UpdateBanner updates an existing banner
func (s *BannerService) UpdateBanner(id int, name string) (*dto.Banner, error) {
	// Validate input
	if id <= 0 {
		return nil, fmt.Errorf("invalid banner ID: %d", id)
	}
	
	if name == "" {
		return nil, fmt.Errorf("banner name cannot be empty")
	}
	
	if len(name) > 255 {
		return nil, fmt.Errorf("banner name cannot exceed 255 characters")
	}
	
	// Check if banner exists
	existingBanner, err := s.repo.GetBannerByID(id)
	if err != nil {
		return nil, err
	}
	
	// Check if another banner with same name exists
	duplicateBanner, err := s.repo.GetBannerByName(name)
	if err == nil && duplicateBanner != nil && duplicateBanner.ID != id {
		return nil, fmt.Errorf("banner with name '%s' already exists", name)
	}
	
	// Update banner
	existingBanner.Name = name
	existingBanner.UpdatedAt = time.Now()
	
	if err := s.repo.UpdateBanner(existingBanner); err != nil {
		return nil, fmt.Errorf("failed to update banner: %w", err)
	}
	
	return existingBanner, nil
}

// DeleteBanner deletes a banner and all its clicks
func (s *BannerService) DeleteBanner(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid banner ID: %d", id)
	}
	
	// Check if banner exists
	_, err := s.repo.GetBannerByID(id)
	if err != nil {
		return err
	}
	
	// Delete banner (clicks will be deleted due to CASCADE)
	if err := s.repo.DeleteBanner(id); err != nil {
		return fmt.Errorf("failed to delete banner: %w", err)
	}
	
	return nil
}

// ClickService provides click business logic
type ClickService struct {
	*Service
}

// NewClickService creates a new click service
func NewClickService(service *Service) *ClickService {
	return &ClickService{Service: service}
}

// RecordClick records a new click for a banner
func (s *ClickService) RecordClick(bannerID int, timestamp time.Time) (*dto.Click, error) {
	// Validate input
	if bannerID <= 0 {
		return nil, fmt.Errorf("invalid banner ID: %d", bannerID)
	}
	
	// Check if banner exists
	_, err := s.repo.GetBannerByID(bannerID)
	if err != nil {
		return nil, fmt.Errorf("banner with ID %d not found: %w", bannerID, err)
	}
	
	// Use current time if timestamp is zero
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	
	// Create new click
	click := &dto.Click{
		Timestamp: timestamp,
		BannerID:  bannerID,
		CreatedAt: time.Now(),
	}
	
	if err := s.repo.CreateClick(click); err != nil {
		return nil, fmt.Errorf("failed to record click: %w", err)
	}
	
	return click, nil
}

// GetClick retrieves a click by ID
func (s *ClickService) GetClick(id int) (*dto.Click, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid click ID: %d", id)
	}
	
	return s.repo.GetClickByID(id)
}

// GetClicksForBanner retrieves all clicks for a specific banner
func (s *ClickService) GetClicksForBanner(bannerID int) ([]*dto.Click, error) {
	if bannerID <= 0 {
		return nil, fmt.Errorf("invalid banner ID: %d", bannerID)
	}
	
	// Check if banner exists
	_, err := s.repo.GetBannerByID(bannerID)
	if err != nil {
		return nil, fmt.Errorf("banner with ID %d not found: %w", bannerID, err)
	}
	
	return s.repo.GetClicksByBannerID(bannerID)
}

// GetClicksInDateRange retrieves clicks within a date range
func (s *ClickService) GetClicksInDateRange(start, end time.Time) ([]*dto.Click, error) {
	if start.After(end) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}
	
	return s.repo.GetClicksByDateRange(start, end)
}

// GetClicksForBannerInDateRange retrieves clicks for a specific banner within a date range
func (s *ClickService) GetClicksForBannerInDateRange(bannerID int, start, end time.Time) ([]*dto.Click, error) {
	if bannerID <= 0 {
		return nil, fmt.Errorf("invalid banner ID: %d", bannerID)
	}
	
	if start.After(end) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}
	
	// Check if banner exists
	_, err := s.repo.GetBannerByID(bannerID)
	if err != nil {
		return nil, fmt.Errorf("banner with ID %d not found: %w", bannerID, err)
	}
	
	return s.repo.GetClicksByBannerIDAndDateRange(bannerID, start, end)
}

// GetClickStats retrieves click statistics for a banner
func (s *ClickService) GetClickStats(bannerID int) (*db.ClickStats, error) {
	if bannerID <= 0 {
		return nil, fmt.Errorf("invalid banner ID: %d", bannerID)
	}
	
	// Check if banner exists
	_, err := s.repo.GetBannerByID(bannerID)
	if err != nil {
		return nil, fmt.Errorf("banner with ID %d not found: %w", bannerID, err)
	}
	
	return s.repo.GetClickStats(bannerID)
}

// DeleteClick deletes a click by ID
func (s *ClickService) DeleteClick(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid click ID: %d", id)
	}
	
	return s.repo.DeleteClick(id)
}

// AnalyticsService provides analytics functionality
type AnalyticsService struct {
	*Service
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(service *Service) *AnalyticsService {
	return &AnalyticsService{Service: service}
}

// GetBannerPerformance retrieves performance metrics for all banners
func (s *AnalyticsService) GetBannerPerformance() ([]*BannerPerformance, error) {
	// Get all banners
	banners, err := s.repo.GetAllBanners()
	if err != nil {
		return nil, fmt.Errorf("failed to get banners: %w", err)
	}
	
	var performances []*BannerPerformance
	for _, banner := range banners {
		stats, err := s.repo.GetClickStats(banner.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get stats for banner %d: %w", banner.ID, err)
		}
		
		performance := &BannerPerformance{
			Banner:      banner,
			TotalClicks: stats.TotalClicks,
			FirstClick:  stats.FirstClick,
			LastClick:   stats.LastClick,
		}
		
		performances = append(performances, performance)
	}
	
	return performances, nil
}

// BannerPerformance represents banner performance metrics
type BannerPerformance struct {
	Banner      *dto.Banner `json:"banner"`
	TotalClicks int         `json:"total_clicks"`
	FirstClick  time.Time   `json:"first_click"`
	LastClick   time.Time   `json:"last_click"`
}