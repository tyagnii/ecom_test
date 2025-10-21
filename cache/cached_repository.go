package cache

import (
	"fmt"
	"time"

	"github.com/tyagnii/ecom_test/db"
	"github.com/tyagnii/ecom_test/dto"
)

// CachedRepository wraps a repository with caching functionality
type CachedRepository struct {
	repo  *db.Repository
	cache Cache
}

// NewCachedRepository creates a new cached repository
func NewCachedRepository(repo *db.Repository, cache Cache) *CachedRepository {
	return &CachedRepository{
		repo:  repo,
		cache: cache,
	}
}

// Banner operations with caching

// CreateBanner creates a new banner and invalidates cache
func (r *CachedRepository) CreateBanner(banner *dto.Banner) error {
	err := r.repo.CreateBanner(banner)
	if err != nil {
		return err
	}

	// Cache the new banner
	r.cache.SetBanner(banner, DefaultBannerTTL)
	
	// Invalidate related caches
	r.cache.InvalidateTopBanners()

	return nil
}

// GetBannerByID retrieves a banner with caching
func (r *CachedRepository) GetBannerByID(id int) (*dto.Banner, error) {
	// Try cache first
	if banner, found := r.cache.GetBanner(id); found {
		return banner, nil
	}

	// Get from database
	banner, err := r.repo.GetBannerByID(id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	r.cache.SetBanner(banner, DefaultBannerTTL)

	return banner, nil
}

// GetAllBanners retrieves all banners (not cached due to frequent changes)
func (r *CachedRepository) GetAllBanners() ([]*dto.Banner, error) {
	return r.repo.GetAllBanners()
}

// UpdateBanner updates a banner and invalidates cache
func (r *CachedRepository) UpdateBanner(banner *dto.Banner) error {
	err := r.repo.UpdateBanner(banner)
	if err != nil {
		return err
	}

	// Update cache
	r.cache.SetBanner(banner, DefaultBannerTTL)
	
	// Invalidate related caches
	r.cache.InvalidateBanner(banner.ID)

	return nil
}

// DeleteBanner deletes a banner and invalidates cache
func (r *CachedRepository) DeleteBanner(id int) error {
	err := r.repo.DeleteBanner(id)
	if err != nil {
		return err
	}

	// Invalidate all related cache entries
	r.cache.InvalidateBanner(id)

	return nil
}

// GetBannerByName retrieves a banner by name (not cached due to low frequency)
func (r *CachedRepository) GetBannerByName(name string) (*dto.Banner, error) {
	return r.repo.GetBannerByName(name)
}

// SearchBannersByName searches banners by name (not cached due to low frequency)
func (r *CachedRepository) SearchBannersByName(namePattern string) ([]*dto.Banner, error) {
	return r.repo.SearchBannersByName(namePattern)
}

// GetBannersWithClickCount retrieves banners with click counts with caching
func (r *CachedRepository) GetBannersWithClickCount() ([]*db.BannerWithStats, error) {
	// This is expensive, so we don't cache the full result
	// Instead, we cache individual banner stats
	return r.repo.GetBannersWithClickCount()
}

// Click operations with caching

// CreateClick creates a new click and invalidates related caches
func (r *CachedRepository) CreateClick(click *dto.Click) error {
	err := r.repo.CreateClick(click)
	if err != nil {
		return err
	}

	// Invalidate click-related caches for this banner
	r.cache.InvalidateClickStats(click.BannerID)
	r.cache.InvalidateBannerWithStats(click.BannerID)
	r.cache.InvalidateTopBanners()

	return nil
}

// GetClickByID retrieves a click by ID (not cached due to low frequency)
func (r *CachedRepository) GetClickByID(id int) (*dto.Click, error) {
	return r.repo.GetClickByID(id)
}

// GetAllClicks retrieves all clicks (not cached due to high volume)
func (r *CachedRepository) GetAllClicks() ([]*dto.Click, error) {
	return r.repo.GetAllClicks()
}

// GetClicksByBannerID retrieves clicks for a banner (not cached due to high volume)
func (r *CachedRepository) GetClicksByBannerID(bannerID int) ([]*dto.Click, error) {
	return r.repo.GetClicksByBannerID(bannerID)
}

// GetClicksByDateRange retrieves clicks in date range (not cached due to high volume)
func (r *CachedRepository) GetClicksByDateRange(start, end time.Time) ([]*dto.Click, error) {
	return r.repo.GetClicksByDateRange(start, end)
}

// GetClicksByBannerIDAndDateRange retrieves clicks for banner in date range (not cached due to high volume)
func (r *CachedRepository) GetClicksByBannerIDAndDateRange(bannerID int, start, end time.Time) ([]*dto.Click, error) {
	return r.repo.GetClicksByBannerIDAndDateRange(bannerID, start, end)
}

// DeleteClick deletes a click and invalidates related caches
func (r *CachedRepository) DeleteClick(id int) error {
	// Get the click first to know which banner to invalidate
	click, err := r.repo.GetClickByID(id)
	if err != nil {
		return err
	}

	err = r.repo.DeleteClick(id)
	if err != nil {
		return err
	}

	// Invalidate click-related caches for this banner
	r.cache.InvalidateClickStats(click.BannerID)
	r.cache.InvalidateBannerWithStats(click.BannerID)
	r.cache.InvalidateTopBanners()

	return nil
}

// GetClickStats retrieves click statistics with caching
func (r *CachedRepository) GetClickStats(bannerID int) (*db.ClickStats, error) {
	// Try cache first
	if stats, found := r.cache.GetClickStats(bannerID); found {
		return stats, nil
	}

	// Get from database
	stats, err := r.repo.GetClickStats(bannerID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	r.cache.SetClickStats(bannerID, stats, DefaultClickStatsTTL)

	return stats, nil
}

// GetTopBanners retrieves top banners with caching
func (r *CachedRepository) GetTopBanners(limit int) ([]*db.BannerClickCount, error) {
	// Try cache first
	if banners, found := r.cache.GetTopBanners(limit); found {
		return banners, nil
	}

	// Get from database
	banners, err := r.repo.GetTopBanners(limit)
	if err != nil {
		return nil, err
	}

	// Cache the result
	r.cache.SetTopBanners(limit, banners, DefaultTopBannersTTL)

	return banners, nil
}

// GetClicksByHour retrieves hourly clicks (not cached due to low frequency)
func (r *CachedRepository) GetClicksByHour(bannerID int, date time.Time) ([]*db.HourlyClicks, error) {
	return r.repo.GetClicksByHour(bannerID, date)
}

// GetClicksByDay retrieves daily clicks (not cached due to low frequency)
func (r *CachedRepository) GetClicksByDay(bannerID int, startDate, endDate time.Time) ([]*db.DailyClicks, error) {
	return r.repo.GetClicksByDay(bannerID, startDate, endDate)
}

// Cache management methods

// GetCacheStats returns cache performance statistics
func (r *CachedRepository) GetCacheStats() CacheStats {
	return r.cache.Stats()
}

// ClearCache clears all cached data
func (r *CachedRepository) ClearCache() {
	r.cache.Clear()
}

// InvalidateBannerCache invalidates all cache entries for a banner
func (r *CachedRepository) InvalidateBannerCache(bannerID int) {
	r.cache.InvalidateBanner(bannerID)
}

// WarmCache preloads frequently accessed data
func (r *CachedRepository) WarmCache() error {
	// Get all banners and cache them
	banners, err := r.repo.GetAllBanners()
	if err != nil {
		return fmt.Errorf("failed to warm banner cache: %w", err)
	}

	for _, banner := range banners {
		r.cache.SetBanner(banner, DefaultBannerTTL)
	}

	// Cache click stats for all banners
	for _, banner := range banners {
		stats, err := r.repo.GetClickStats(banner.ID)
		if err != nil {
			continue // Skip if stats can't be retrieved
		}
		r.cache.SetClickStats(banner.ID, stats, DefaultClickStatsTTL)
	}

	// Cache top banners
	topBanners, err := r.repo.GetTopBanners(10)
	if err == nil {
		r.cache.SetTopBanners(10, topBanners, DefaultTopBannersTTL)
	}

	return nil
}

// GetCache returns the underlying cache instance
func (r *CachedRepository) GetCache() Cache {
	return r.cache
}
