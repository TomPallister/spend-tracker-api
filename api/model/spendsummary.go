package model

import "github.com/shopspring/decimal"

// SpendSummary ...
type SpendSummary struct {
	ID           int64            `json:"id"`
	TrackerID    int64            `json:"trackerId"`
	UserID       int64            `json:"userId"`
	Value        decimal.Decimal  `json:"value"`
	Currency     string           `json:"currency"`
}
