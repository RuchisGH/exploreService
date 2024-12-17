package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	models "example.com/mod/internal/models"
)

// Database is a wrapper around the SQL database connection
type Database struct {
	Conn *sql.DB
}

// NewDatabase creates a new Database instance
func NewDatabase(conn *sql.DB) *Database {
	return &Database{Conn: conn}
}

// DatabaseInterface defines the required methods for the database.
type DatabaseInterface interface {
	GetReceivedLikes(userID string) ([]string, error)
	GetGivenDecisions(userID string) ([]models.Decision, error)
	UpsertDecision(userID, targetID string, decision bool) error
}

// Ensure Database implements DatabaseInterface
var _ DatabaseInterface = (*Database)(nil)

// UpsertDecision inserts or updates a decision in the database.
func (d *Database) UpsertDecision(userID, targetID string, decision bool) error {
	query := `
        INSERT INTO decisions (user_id, target_id, decision, timestamp)
        VALUES (?, ?, ?, NOW())
        ON DUPLICATE KEY UPDATE decision = VALUES(decision), timestamp = NOW()`
	_, err := d.Conn.Exec(query, userID, targetID, decision)
	return err
}

// GetReceivedLikes fetches the IDs of users who liked the given user.
func (d *Database) GetReceivedLikes(userID string) ([]string, error) {
	// Query to get all users who liked the given user
	query := `
        SELECT user_id
        FROM decisions
        WHERE target_id = ? AND decision = 'LIKE'
		LIMIT 50 OFFSET ?`
	rows, err := d.Conn.Query(query, userID)
	if err != nil {
		log.Printf("Error fetching received likes: %v", err)
		return nil, err
	}
	defer rows.Close()

	var likers []string
	for rows.Next() {
		var liker string
		if err := rows.Scan(&liker); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		likers = append(likers, liker)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return likers, nil
}

// GetGivenDecisions fetches all decisions (LIKE or PASS) made by the given user.
func (d *Database) GetGivenDecisions(userID string) ([]models.Decision, error) {
	query := `
        SELECT target_id, decision, timestamp
        FROM decisions
        WHERE user_id = ?`
	rows, err := d.Conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decisions []models.Decision
	for rows.Next() {
		var d models.Decision
		var timestamp []byte // Temporarily store timestamp as byte slice

		// Scan the values, including the timestamp as a byte slice
		if err := rows.Scan(&d.TargetID, &d.Decision, &timestamp); err != nil {
			return nil, err
		}

		// Convert timestamp to time.Time
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(timestamp))
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp: %v", err)
		}
		d.Timestamp = parsedTime

		d.UserID = userID
		decisions = append(decisions, d)
	}

	return decisions, nil
}

// DeleteDecision removes a specific decision made by a user for a target profile.
func DeleteDecision(db *sql.DB, userID, targetID string) error {
	query := `
        DELETE FROM decisions
        WHERE user_id = ? AND target_id = ?`
	_, err := db.Exec(query, userID, targetID)
	return err
}
