package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Spend ...
type Spend struct {
	ID          int64           `json:"id"`
	Value       decimal.Decimal `json:"value"`
	TrackerID   int64           `json:"trackerId"`
	Name        string          `json:"name"`
	UserID      int64           `json:"userId"`
	Currency    string          `json:"currency"`
	DateCreated time.Time       `json:"dateCreated"`
}
