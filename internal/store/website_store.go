package store

import (
	"database/sql"
	"time"

	"github.com/Adi-ty/betteruptime-go/internal/stream"
)

type WebsiteStatus string

const (
    StatusUp      WebsiteStatus = "UP"
    StatusDown    WebsiteStatus = "DOWN"
    StatusUnknown WebsiteStatus = "UNKNOWN"
)

type Website struct {
	ID        int64 `json:"id"`
	Url       string `json:"url"`
	UserID    int64  `json:"user_id"`
	TimeAdded time.Time `json:"time_added"`
	// Regions []Region
	WebsiteTicks []WebsiteTick	 `json:"website_ticks"`
}

// type Region struct {
// 	ID   string
// 	Name string
// }

type WebsiteTick struct {
	ID             string			`json:"id"`
	ResponseTimeMs int64			`json:"response_time_ms"`
	StatusCode     WebsiteStatus	`json:"status_code"`
	WebsiteID      int64		    `json:"website_id"`
	RegionID       string			`json:"region_id"`
}

type PostgresWebsiteStore struct {
	db *sql.DB
}

func NewPostgresWebsiteStore(db *sql.DB) *PostgresWebsiteStore {
	return &PostgresWebsiteStore{db: db}
}

type WebsiteStore interface {
	CreateWebsite(Website *Website) error
	GetWebsiteStatusByID(userId int64, id int64) (*Website, error)
	GetAllWebsites() ([]*stream.WebsiteEvent, error)
	MarkWebsiteTickProcessed(websiteTick *WebsiteTick) error
}

func (s *PostgresWebsiteStore) CreateWebsite(website *Website) error {
	query := "INSERT INTO \"website\" (url, user_id, time_added) VALUES ($1, $2, $3) RETURNING id"
	err := s.db.QueryRow(query, website.Url, website.UserID, website.TimeAdded).Scan(&website.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresWebsiteStore) GetWebsiteStatusByID(userId int64, id int64) (*Website, error) {
	website := &Website{}
	var tickID sql.NullString
    var responseTimeMs sql.NullInt64
    var statusCode sql.NullString
    var websiteID sql.NullInt64
    var regionID sql.NullString

	query := `
			SELECT
				w.id, w.url, w.user_id, w.time_added,
				t.id, t.response_time_ms, t.status_code, t.website_id, t.region_id
			FROM "website" AS w
			LEFT JOIN "website_tick" AS t
			ON t.website_id = w.id
			WHERE
			w.user_id = $1
			AND w.id = $2
			ORDER BY t."created_at" DESC
			LIMIT 1;
		`

	err := s.db.QueryRow(query, userId, id).Scan(
        &website.ID, &website.Url, &website.UserID, &website.TimeAdded,
        &tickID, &responseTimeMs, &statusCode, &websiteID, &regionID,
    )
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if tickID.Valid {
        tick := WebsiteTick{
            ID:             tickID.String,
            ResponseTimeMs: responseTimeMs.Int64,
            StatusCode:     WebsiteStatus(statusCode.String),
            WebsiteID:      websiteID.Int64,
            RegionID:       regionID.String,
        }
        website.WebsiteTicks = []WebsiteTick{tick}
    } else {
        website.WebsiteTicks = []WebsiteTick{}
    }

	return website, nil
}

func (s *PostgresWebsiteStore) GetAllWebsites() ([]*stream.WebsiteEvent, error) {
	query := `SELECT id, url FROM "website"`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var websites []*stream.WebsiteEvent
	for rows.Next() {
		website := &stream.WebsiteEvent{}
		err := rows.Scan(&website.ID, &website.Url)
		if err != nil {
			return nil, err
		}
		websites = append(websites, website)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return websites, nil
}

func (s *PostgresWebsiteStore) MarkWebsiteTickProcessed(websiteTick *WebsiteTick) error {
	query := `INSERT INTO "website_tick" (id, website_id, region_id, response_time_ms, status_code, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.Exec(query, websiteTick.ID, websiteTick.WebsiteID, websiteTick.RegionID, websiteTick.ResponseTimeMs, websiteTick.StatusCode, time.Now())
	if err != nil {
		return err
	}
	
	return nil
}