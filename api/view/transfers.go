package view

import "github.com/TomPallister/godutch-api/api/model"

// Transfers ...
type Transfers struct {
	Transfers []model.Transfer `json:"transfers"`
}
