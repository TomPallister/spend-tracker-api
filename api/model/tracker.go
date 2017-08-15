package model

import "time"

// Tracker ...
type Tracker struct {
	ID             int64      `json:"id"`
	AdminUserID    int64      `json:"adminUserId"`
	TrackerUserIDs []int64    `json:"trackerUserIds"`
	Name           string    `json:"name"`
	DateCreated    time.Time `json:"dateCreated"`
	Currency       string    `json:"currency"`
}
