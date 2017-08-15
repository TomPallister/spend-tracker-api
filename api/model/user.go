package model

import "time"

// User ....
type User struct {
	ID               int64      `json:"id"`
	Name             string    `json:"name"`
	AuthenticationID string    `json:"authenticationID"`
	EmailAddress     string    `json:"emailAddress"`
	DateCreated      time.Time `json:"dateCreated"`
}
