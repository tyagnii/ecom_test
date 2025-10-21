package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tyagnii/ecom_test/dto"
)

// Repository provides database operations
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// ClickStats represents click statistics
type ClickStats struct {
	BannerID    int       `json:"banner_id"`
	TotalClicks int       `json:"total_clicks"`
	FirstClick  time.Time `json:"first_click"`
	LastClick   time.Time `json:"last_click"`
}

// BannerWithStats represents a banner with click statistics
type BannerWithStats struct {
	Banner      *dto.Banner `json:"banner"`
	ClickCount int          `json:"click_count"`
	LastClick   *time.Time  `json:"last_click"`
}

// BannerClickCount represents banner click count for top banners
type BannerClickCount struct {
	BannerID   int    `json:"banner_id"`
	BannerName string `json:"banner_name"`
	ClickCount int    `json:"click_count"`
}

// HourlyClicks represents clicks per hour
type HourlyClicks struct {
	Hour       int `json:"hour"`
	ClickCount int `json:"click_count"`
}

// DailyClicks represents clicks per day
type DailyClicks struct {
	Date       time.Time `json:"date"`
	ClickCount int       `json:"click_count"`
}

// Banner CRUD Operations

// CreateBanner creates a new banner
func (r *Repository) CreateBanner(banner *dto.Banner) error {
	query := `
		INSERT INTO banners (name, created_at, updated_at) 
		VALUES ($1, $2, $3) 
		RETURNING id`
	
	err := r.db.QueryRow(
		query,
		banner.Name,
		banner.CreatedAt,
		banner.UpdatedAt,
	).Scan(&banner.ID)
	
	if err != nil {
		return fmt.Errorf("failed to create banner: %w", err)
	}
	
	return nil
}

// GetBannerByID retrieves a banner by ID
func (r *Repository) GetBannerByID(id int) (*dto.Banner, error) {
	query := `
		SELECT id, name, created_at, updated_at 
		FROM banners 
		WHERE id = $1`
	
	banner := &dto.Banner{}
	err := r.db.QueryRow(query, id).Scan(
		&banner.ID,
		&banner.Name,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("banner with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}
	
	return banner, nil
}

// GetAllBanners retrieves all banners
func (r *Repository) GetAllBanners() ([]*dto.Banner, error) {
	query := `
		SELECT id, name, created_at, updated_at 
		FROM banners 
		ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get banners: %w", err)
	}
	defer rows.Close()
	
	var banners []*dto.Banner
	for rows.Next() {
		banner := &dto.Banner{}
		err := rows.Scan(
			&banner.ID,
			&banner.Name,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan banner: %w", err)
		}
		banners = append(banners, banner)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating banners: %w", err)
	}
	
	return banners, nil
}

// UpdateBanner updates an existing banner
func (r *Repository) UpdateBanner(banner *dto.Banner) error {
	query := `
		UPDATE banners 
		SET name = $1, updated_at = $2 
		WHERE id = $3`
	
	result, err := r.db.Exec(query, banner.Name, banner.UpdatedAt, banner.ID)
	if err != nil {
		return fmt.Errorf("failed to update banner: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("banner with ID %d not found", banner.ID)
	}
	
	return nil
}

// DeleteBanner deletes a banner by ID
func (r *Repository) DeleteBanner(id int) error {
	query := `DELETE FROM banners WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete banner: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("banner with ID %d not found", id)
	}
	
	return nil
}

// GetBannerByName retrieves a banner by name
func (r *Repository) GetBannerByName(name string) (*dto.Banner, error) {
	query := `
		SELECT id, name, created_at, updated_at 
		FROM banners 
		WHERE name = $1`
	
	banner := &dto.Banner{}
	err := r.db.QueryRow(query, name).Scan(
		&banner.ID,
		&banner.Name,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("banner with name '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get banner by name: %w", err)
	}
	
	return banner, nil
}

// SearchBannersByName searches banners by name pattern
func (r *Repository) SearchBannersByName(namePattern string) ([]*dto.Banner, error) {
	query := `
		SELECT id, name, created_at, updated_at 
		FROM banners 
		WHERE name ILIKE $1 
		ORDER BY name`
	
	rows, err := r.db.Query(query, "%"+namePattern+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search banners: %w", err)
	}
	defer rows.Close()
	
	var banners []*dto.Banner
	for rows.Next() {
		banner := &dto.Banner{}
		err := rows.Scan(
			&banner.ID,
			&banner.Name,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan banner: %w", err)
		}
		banners = append(banners, banner)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating banners: %w", err)
	}
	
	return banners, nil
}

// GetBannersWithClickCount retrieves banners with their click counts
func (r *Repository) GetBannersWithClickCount() ([]*BannerWithStats, error) {
	query := `
		SELECT 
			b.id, b.name, b.created_at, b.updated_at,
			COUNT(c.id) as click_count,
			MAX(c.timestamp) as last_click
		FROM banners b
		LEFT JOIN clicks c ON b.id = c.bannerid
		GROUP BY b.id, b.name, b.created_at, b.updated_at
		ORDER BY click_count DESC, b.created_at DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get banners with click count: %w", err)
	}
	defer rows.Close()
	
	var results []*BannerWithStats
	for rows.Next() {
		result := &BannerWithStats{
			Banner: &dto.Banner{},
		}
		var lastClick sql.NullTime
		
		err := rows.Scan(
			&result.Banner.ID,
			&result.Banner.Name,
			&result.Banner.CreatedAt,
			&result.Banner.UpdatedAt,
			&result.ClickCount,
			&lastClick,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan banner with stats: %w", err)
		}
		
		if lastClick.Valid {
			result.LastClick = &lastClick.Time
		}
		
		results = append(results, result)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating banners with stats: %w", err)
	}
	
	return results, nil
}

// Click CRUD Operations

// CreateClick creates a new click
func (r *Repository) CreateClick(click *dto.Click) error {
	query := `
		INSERT INTO clicks (timestamp, bannerid, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id`
	
	err := r.db.QueryRow(
		query,
		click.Timestamp,
		click.BannerID,
		click.CreatedAt,
	).Scan(&click.ID)
	
	if err != nil {
		return fmt.Errorf("failed to create click: %w", err)
	}
	
	return nil
}

// GetClickByID retrieves a click by ID
func (r *Repository) GetClickByID(id int) (*dto.Click, error) {
	query := `
		SELECT id, timestamp, bannerid, created_at 
		FROM clicks 
		WHERE id = $1`
	
	click := &dto.Click{}
	err := r.db.QueryRow(query, id).Scan(
		&click.ID,
		&click.Timestamp,
		&click.BannerID,
		&click.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("click with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get click: %w", err)
	}
	
	return click, nil
}

// GetAllClicks retrieves all clicks
func (r *Repository) GetAllClicks() ([]*dto.Click, error) {
	query := `
		SELECT id, timestamp, bannerid, created_at 
		FROM clicks 
		ORDER BY timestamp DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks: %w", err)
	}
	defer rows.Close()
	
	var clicks []*dto.Click
	for rows.Next() {
		click := &dto.Click{}
		err := rows.Scan(
			&click.ID,
			&click.Timestamp,
			&click.BannerID,
			&click.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan click: %w", err)
		}
		clicks = append(clicks, click)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating clicks: %w", err)
	}
	
	return clicks, nil
}

// GetClicksByBannerID retrieves clicks for a specific banner
func (r *Repository) GetClicksByBannerID(bannerID int) ([]*dto.Click, error) {
	query := `
		SELECT id, timestamp, bannerid, created_at 
		FROM clicks 
		WHERE bannerid = $1 
		ORDER BY timestamp DESC`
	
	rows, err := r.db.Query(query, bannerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks by banner ID: %w", err)
	}
	defer rows.Close()
	
	var clicks []*dto.Click
	for rows.Next() {
		click := &dto.Click{}
		err := rows.Scan(
			&click.ID,
			&click.Timestamp,
			&click.BannerID,
			&click.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan click: %w", err)
		}
		clicks = append(clicks, click)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating clicks: %w", err)
	}
	
	return clicks, nil
}

// GetClicksByDateRange retrieves clicks within a date range
func (r *Repository) GetClicksByDateRange(start, end time.Time) ([]*dto.Click, error) {
	query := `
		SELECT id, timestamp, bannerid, created_at 
		FROM clicks 
		WHERE timestamp BETWEEN $1 AND $2 
		ORDER BY timestamp DESC`
	
	rows, err := r.db.Query(query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks by date range: %w", err)
	}
	defer rows.Close()
	
	var clicks []*dto.Click
	for rows.Next() {
		click := &dto.Click{}
		err := rows.Scan(
			&click.ID,
			&click.Timestamp,
			&click.BannerID,
			&click.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan click: %w", err)
		}
		clicks = append(clicks, click)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating clicks: %w", err)
	}
	
	return clicks, nil
}

// GetClicksByBannerIDAndDateRange retrieves clicks for a specific banner within a date range
func (r *Repository) GetClicksByBannerIDAndDateRange(bannerID int, start, end time.Time) ([]*dto.Click, error) {
	query := `
		SELECT id, timestamp, bannerid, created_at 
		FROM clicks 
		WHERE bannerid = $1 AND timestamp BETWEEN $2 AND $3 
		ORDER BY timestamp DESC`
	
	rows, err := r.db.Query(query, bannerID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks by banner ID and date range: %w", err)
	}
	defer rows.Close()
	
	var clicks []*dto.Click
	for rows.Next() {
		click := &dto.Click{}
		err := rows.Scan(
			&click.ID,
			&click.Timestamp,
			&click.BannerID,
			&click.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan click: %w", err)
		}
		clicks = append(clicks, click)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating clicks: %w", err)
	}
	
	return clicks, nil
}

// DeleteClick deletes a click by ID
func (r *Repository) DeleteClick(id int) error {
	query := `DELETE FROM clicks WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete click: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("click with ID %d not found", id)
	}
	
	return nil
}

// GetClickStats retrieves click statistics for a banner
func (r *Repository) GetClickStats(bannerID int) (*ClickStats, error) {
	query := `
		SELECT 
			bannerid,
			COUNT(*) as total_clicks,
			MIN(timestamp) as first_click,
			MAX(timestamp) as last_click
		FROM clicks 
		WHERE bannerid = $1
		GROUP BY bannerid`
	
	stats := &ClickStats{}
	err := r.db.QueryRow(query, bannerID).Scan(
		&stats.BannerID,
		&stats.TotalClicks,
		&stats.FirstClick,
		&stats.LastClick,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return &ClickStats{
				BannerID:    bannerID,
				TotalClicks: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to get click stats: %w", err)
	}
	
	return stats, nil
}

// GetTopBanners retrieves top banners by click count
func (r *Repository) GetTopBanners(limit int) ([]*BannerClickCount, error) {
	query := `
		SELECT 
			b.id as banner_id,
			b.name as banner_name,
			COUNT(c.id) as click_count
		FROM banners b
		LEFT JOIN clicks c ON b.id = c.bannerid
		GROUP BY b.id, b.name
		ORDER BY click_count DESC, b.name
		LIMIT $1`
	
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top banners: %w", err)
	}
	defer rows.Close()
	
	var results []*BannerClickCount
	for rows.Next() {
		result := &BannerClickCount{}
		err := rows.Scan(
			&result.BannerID,
			&result.BannerName,
			&result.ClickCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan top banner: %w", err)
		}
		results = append(results, result)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating top banners: %w", err)
	}
	
	return results, nil
}

// GetClicksByHour retrieves hourly click distribution for a banner on a specific date
func (r *Repository) GetClicksByHour(bannerID int, date time.Time) ([]*HourlyClicks, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	query := `
		SELECT 
			EXTRACT(HOUR FROM timestamp) as hour,
			COUNT(*) as click_count
		FROM clicks 
		WHERE bannerid = $1 
		AND timestamp >= $2 
		AND timestamp < $3
		GROUP BY EXTRACT(HOUR FROM timestamp)
		ORDER BY hour`
	
	rows, err := r.db.Query(query, bannerID, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks by hour: %w", err)
	}
	defer rows.Close()
	
	var results []*HourlyClicks
	for rows.Next() {
		result := &HourlyClicks{}
		err := rows.Scan(
			&result.Hour,
			&result.ClickCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan hourly clicks: %w", err)
		}
		results = append(results, result)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating hourly clicks: %w", err)
	}
	
	return results, nil
}

// GetClicksByDay retrieves daily click distribution for a banner within a date range
func (r *Repository) GetClicksByDay(bannerID int, startDate, endDate time.Time) ([]*DailyClicks, error) {
	query := `
		SELECT 
			DATE(timestamp) as date,
			COUNT(*) as click_count
		FROM clicks 
		WHERE bannerid = $1 
		AND timestamp >= $2 
		AND timestamp <= $3
		GROUP BY DATE(timestamp)
		ORDER BY date`
	
	rows, err := r.db.Query(query, bannerID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks by day: %w", err)
	}
	defer rows.Close()
	
	var results []*DailyClicks
	for rows.Next() {
		result := &DailyClicks{}
		err := rows.Scan(
			&result.Date,
			&result.ClickCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan daily clicks: %w", err)
		}
		results = append(results, result)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating daily clicks: %w", err)
	}
	
	return results, nil
}
