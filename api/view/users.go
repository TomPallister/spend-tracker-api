package view

import "github.com/TomPallister/godutch-api/api/model"

// Users ...
type Users struct {
	Users []model.User `json:"users"`
}
