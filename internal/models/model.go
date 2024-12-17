package models

import (
	"time"
)

// Decision represents a user's decision on another user's profile.
type Decision struct {
	UserID    string    // The ID of the user making the decision
	TargetID  string    // The ID of the target user profile
	Decision  string    // "LIKE" or "PASS"
	Timestamp time.Time // The time when the decision was made
}
