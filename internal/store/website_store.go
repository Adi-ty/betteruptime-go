package store

import (
	"database/sql"
	"fmt"
)

// type WebsiteStatus string

// const (
//     StatusUp      WebsiteStatus = "UP"
//     StatusDown    WebsiteStatus = "DOWN"
//     StatusUnknown WebsiteStatus = "UNKNOWN"
// )

type Website struct {
	ID        int64 `json:"id"`
	Url       string `json:"url"`
	// Regions []Region
	// WebsiteTicks []WebsiteTick
}

// type Region struct {
// 	ID   string
// 	Name string
// }

// type WebsiteTick struct {
// 	ID             string
// 	ResponseTimeMs int
// 	StatusCode     WebsiteStatus
// 	WebsiteID      string
// 	RegionID       string
// }

type PostgresWebsiteStore struct {
	db *sql.DB
}

func NewPostgresWebsiteStore(db *sql.DB) *PostgresWebsiteStore {
	return &PostgresWebsiteStore{db: db}
}

type WebsiteStore interface {
	CreateWebsite(Website *Website) (*Website, error)
	GetWebsiteByID(id int64) (*Website, error)
}

func (s *PostgresWebsiteStore) CreateWebsite(website *Website) (*Website, error) {
	query := "INSERT INTO website (url) VALUES ($1) RETURNING id"
	err := s.db.QueryRow(query, website.Url).Scan(&website.ID)
	if err != nil {
		fmt.Println("here: %w", err)
		return nil, err
	}

	return website, nil
}

func (s *PostgresWebsiteStore) GetWebsiteByID(id int64) (*Website, error) {
	website := &Website{}
	query := "SELECT id, url FROM website WHERE id = $1"

	err := s.db.QueryRow(query, id).Scan(&website.ID, &website.Url)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return website, nil
}
