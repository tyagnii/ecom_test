package cache

import (
	"sync"
	"time"

	"github.com/tyagnii/ecom_test/db"
	"github.com/tyagnii/ecom_test/dto"
)

// Cache interface defines cache operations
type Cache interface {
	// Banner operations
	GetBanner(id int) (*dto.Banner, bool)
	SetBanner(banner *dto.Banner, ttl time.Duration)
	DeleteBanner(id int)
	InvalidateBanner(id int)

	// Click statistics
	GetClickStats(bannerID int) (*db.ClickStats, bool)
	SetClickStats(bannerID int, stats *db.ClickStats, ttl time.Duration)
	InvalidateClickStats(bannerID int)

	// Banner with stats
	GetBannerWithStats(id int) (*db.BannerWithStats, bool)
	SetBannerWithStats(id int, stats *db.BannerWithStats, ttl time.Duration)
	InvalidateBannerWithStats(id int)

	// Top banners
	GetTopBanners(limit int) ([]*db.BannerClickCount, bool)
	SetTopBanners(limit int, banners []*db.BannerClickCount, ttl time.Duration)
	InvalidateTopBanners()

	// Cache management
	Clear()
	Size() int
	Stats() CacheStats
	Stop()
}

// CacheItem represents a cached item with expiration
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// IsExpired checks if the cache item has expired
func (item *CacheItem) IsExpired() bool {
	return time.Now().After(item.ExpiresAt)
}

// InMemoryCache implements an in-memory cache with TTL support
type InMemoryCache struct {
	mu       sync.RWMutex
	items    map[string]*CacheItem
	stats    CacheStats
	cleanup  *time.Ticker
	stopChan chan struct{}
}

// CacheStats provides cache performance metrics
type CacheStats struct {
	Hits       int64 `json:"hits"`
	Misses     int64 `json:"misses"`
	Sets       int64 `json:"sets"`
	Deletes    int64 `json:"deletes"`
	Expirations int64 `json:"expirations"`
	Size       int   `json:"size"`
}

// NewInMemoryCache creates a new in-memory cache
func NewInMemoryCache(cleanupInterval time.Duration) *InMemoryCache {
	cache := &InMemoryCache{
		items:    make(map[string]*CacheItem),
		cleanup:  time.NewTicker(cleanupInterval),
		stopChan: make(chan struct{}),
	}

	// Start cleanup goroutine
	go cache.cleanupExpired()

	return cache
}

// Stop stops the cache cleanup goroutine
func (c *InMemoryCache) Stop() {
	c.cleanup.Stop()
	close(c.stopChan)
}

// cleanupExpired removes expired items from the cache
func (c *InMemoryCache) cleanupExpired() {
	for {
		select {
		case <-c.cleanup.C:
			c.mu.Lock()
			expiredCount := 0
			for key, item := range c.items {
				if item.IsExpired() {
					delete(c.items, key)
					expiredCount++
				}
			}
			c.stats.Expirations += int64(expiredCount)
			c.mu.Unlock()
		case <-c.stopChan:
			return
		}
	}
}

// get retrieves an item from cache
func (c *InMemoryCache) get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		c.mu.Lock()
		c.stats.Misses++
		c.mu.Unlock()
		return nil, false
	}

	if item.IsExpired() {
		c.mu.Lock()
		delete(c.items, key)
		c.stats.Expirations++
		c.stats.Misses++
		c.mu.Unlock()
		return nil, false
	}

	c.mu.Lock()
	c.stats.Hits++
	c.mu.Unlock()
	return item.Value, true
}

// set stores an item in cache with TTL
func (c *InMemoryCache) set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
	c.stats.Sets++
}

// delete removes an item from cache
func (c *InMemoryCache) delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.items[key]; exists {
		delete(c.items, key)
		c.stats.Deletes++
	}
}

// Banner operations

// GetBanner retrieves a banner from cache
func (c *InMemoryCache) GetBanner(id int) (*dto.Banner, bool) {
	key := "banner:" + string(rune(id))
	value, found := c.get(key)
	if !found {
		return nil, false
	}
	banner, ok := value.(*dto.Banner)
	return banner, ok
}

// SetBanner stores a banner in cache
func (c *InMemoryCache) SetBanner(banner *dto.Banner, ttl time.Duration) {
	key := "banner:" + string(rune(banner.ID))
	c.set(key, banner, ttl)
}

// DeleteBanner removes a banner from cache
func (c *InMemoryCache) DeleteBanner(id int) {
	key := "banner:" + string(rune(id))
	c.delete(key)
}

// InvalidateBanner invalidates banner and related data
func (c *InMemoryCache) InvalidateBanner(id int) {
	c.DeleteBanner(id)
	c.InvalidateClickStats(id)
	c.InvalidateBannerWithStats(id)
	c.InvalidateTopBanners()
}

// Click statistics operations

// GetClickStats retrieves click statistics from cache
func (c *InMemoryCache) GetClickStats(bannerID int) (*db.ClickStats, bool) {
	key := "click_stats:" + string(rune(bannerID))
	value, found := c.get(key)
	if !found {
		return nil, false
	}
	stats, ok := value.(*db.ClickStats)
	return stats, ok
}

// SetClickStats stores click statistics in cache
func (c *InMemoryCache) SetClickStats(bannerID int, stats *db.ClickStats, ttl time.Duration) {
	key := "click_stats:" + string(rune(bannerID))
	c.set(key, stats, ttl)
}

// InvalidateClickStats removes click statistics from cache
func (c *InMemoryCache) InvalidateClickStats(bannerID int) {
	key := "click_stats:" + string(rune(bannerID))
	c.delete(key)
}

// Banner with stats operations

// GetBannerWithStats retrieves banner with stats from cache
func (c *InMemoryCache) GetBannerWithStats(id int) (*db.BannerWithStats, bool) {
	key := "banner_stats:" + string(rune(id))
	value, found := c.get(key)
	if !found {
		return nil, false
	}
	stats, ok := value.(*db.BannerWithStats)
	return stats, ok
}

// SetBannerWithStats stores banner with stats in cache
func (c *InMemoryCache) SetBannerWithStats(id int, stats *db.BannerWithStats, ttl time.Duration) {
	key := "banner_stats:" + string(rune(id))
	c.set(key, stats, ttl)
}

// InvalidateBannerWithStats removes banner with stats from cache
func (c *InMemoryCache) InvalidateBannerWithStats(id int) {
	key := "banner_stats:" + string(rune(id))
	c.delete(key)
}

// Top banners operations

// GetTopBanners retrieves top banners from cache
func (c *InMemoryCache) GetTopBanners(limit int) ([]*db.BannerClickCount, bool) {
	key := "top_banners:" + string(rune(limit))
	value, found := c.get(key)
	if !found {
		return nil, false
	}
	banners, ok := value.([]*db.BannerClickCount)
	return banners, ok
}

// SetTopBanners stores top banners in cache
func (c *InMemoryCache) SetTopBanners(limit int, banners []*db.BannerClickCount, ttl time.Duration) {
	key := "top_banners:" + string(rune(limit))
	c.set(key, banners, ttl)
}

// InvalidateTopBanners removes top banners from cache
func (c *InMemoryCache) InvalidateTopBanners() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove all top_banners keys
	for key := range c.items {
		if len(key) > 11 && key[:11] == "top_banners:" {
			delete(c.items, key)
			c.stats.Deletes++
		}
	}
}

// Cache management

// Clear removes all items from cache
func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*CacheItem)
}

// Size returns the number of items in cache
func (c *InMemoryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Stats returns cache performance statistics
func (c *InMemoryCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := c.stats
	stats.Size = len(c.items)
	return stats
}

// Default TTL values
const (
	DefaultBannerTTL      = 5 * time.Minute
	DefaultClickStatsTTL  = 2 * time.Minute
	DefaultBannerStatsTTL = 3 * time.Minute
	DefaultTopBannersTTL  = 1 * time.Minute
	DefaultCleanupInterval = 30 * time.Second
)
