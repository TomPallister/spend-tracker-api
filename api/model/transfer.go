package model

import "github.com/shopspring/decimal"

// Transfer ...
type Transfer struct {
	ID               int64            `json:"id"`
	TrackerID        int64            `json:"trackerId"`
	FromUserID       int64            `json:"fromUserId"`
	ToUserID         int64            `json:"toUserId"`
	Value            decimal.Decimal `json:"value"`
	Currency         string          `json:"currency"`
}
