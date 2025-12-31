package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type User struct {
	ID           int64       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	GetUserByUsername(username string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	GetUserByToken(scope, tokenPlaintext string) (*User, error)
}

func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query := `SELECT id, username, email, password_hash, created_at, updated_at
			  FROM "user"
			  WHERE username = $1`
	
	err := s.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `INSERT INTO "user" (username, email, password_hash)
			  VALUES ($1, $2, $3)
			  RETURNING id, created_at, updated_at`
			
	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserStore) UpdateUser(user *User) error {
	query := `UPDATE "user" SET email = $1, password_hash = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	res, err := s.db.Exec(query, user.Email, user.PasswordHash.hash, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *PostgresUserStore) GetUserByToken(scope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.created_at, u.updated_at
		FROM "user" u
		INNER JOIN tokens t ON u.id = t.user_id
		WHERE t.hash = $1 AND t.scope = $2 AND t.expiry > $3
	`

	user := &User{
		PasswordHash: password{},
	}

	err := s.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}