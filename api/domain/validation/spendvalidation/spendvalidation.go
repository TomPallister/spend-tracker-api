package spendvalidation

import (
	"errors"

	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/shopspring/decimal"
)

// ErrorInvalidName ...
var ErrorInvalidName = errors.New("Invalid name")

// ErrorInvalidCurrency ...
var ErrorInvalidCurrency = errors.New("Invalid currency")

// ErrorSpendCannotBeLessThanZero ...
var ErrorSpendCannotBeLessThanZero = errors.New("Invalid spend amount")

// ErrorInvalidTrackerID ...
var ErrorInvalidTrackerID = errors.New("Invalid tracker id")

// ErrorInvalidUserID ...
var ErrorInvalidUserID = errors.New("Invalid user id")
 
// ErrorAdminUserIDIsDifferentToSubjectID ...
var ErrorAdminUserIDIsDifferentToSubjectID = errors.New("user id is different to subject id")

// ErrorInvalidDateCreated ...
var ErrorInvalidDateCreated = errors.New("Invalid date created")

// ErrorTrackerIDIsDifferntToTrackTrackerID ...
var ErrorTrackerIDIsDifferntToTrackTrackerID = errors.New("tracker id is different to track tracker id")

// ErrorTheSpendDoesNotExist ...
var ErrorTheSpendDoesNotExist = errors.New("The spend does not exist")

// ErrorSpendCurrencyDoesNotMatchTrackerCurrency ...
var ErrorSpendCurrencyDoesNotMatchTrackerCurrency = errors.New("The spend currency does not match the tracker currency")

// SpendValidator ...
type SpendValidator interface {
	IsValidCreateSpend(spend model.Spend, logger infrastructure.Logger,
		user model.User, tracker model.Tracker) (bool, error)

	IsValidUpdateSpend(spend model.Spend,
		logger infrastructure.Logger, user model.User,
		existingSpend model.Spend, tracker model.Tracker) (bool, error)

	IsValidDeleteSpend(id int64,
		logger infrastructure.Logger, user model.User,
		existingSpend model.Spend) (bool, error)
}

// GoDutchSpendValidator ...
type GoDutchSpendValidator struct {
}

// NewGoDutchSpendValidator ...
func NewGoDutchSpendValidator() *GoDutchSpendValidator {

	service := GoDutchSpendValidator{}

	return &service
}

// IsValidCreateSpend ...
func (validator *GoDutchSpendValidator) IsValidCreateSpend(spend model.Spend,
	logger infrastructure.Logger, user model.User, tracker model.Tracker) (bool, error) {

	if len(spend.Name) <= 0 {
		logger.Error("Error: ", ErrorInvalidName)
		return false, ErrorInvalidName
	}

	if spend.Value.Cmp(decimal.NewFromFloat(0)) == -1 {
		logger.Error("Error: ", ErrorSpendCannotBeLessThanZero)
		return false, ErrorSpendCannotBeLessThanZero
	}

	if spend.Value.Cmp(decimal.NewFromFloat(0)) == 0 {
		logger.Error("Error: ", ErrorSpendCannotBeLessThanZero)
		return false, ErrorSpendCannotBeLessThanZero
	}

	if spend.TrackerID <= 0 {
		logger.Error("Error: ", ErrorInvalidTrackerID)
		return false, ErrorInvalidTrackerID
	}

	if spend.UserID <= 0 {
		logger.Error("Error: ", ErrorInvalidUserID)
		return false, ErrorInvalidUserID
	}

	if len(spend.Currency) <= 0 {
		logger.Error("Error: ", ErrorInvalidCurrency)
		return false, ErrorInvalidCurrency
	}

	if spend.DateCreated.IsZero() {
		logger.Error("Error: ", ErrorInvalidDateCreated)
		return false, ErrorInvalidDateCreated
	}

	if spend.UserID != user.ID {
		logger.Error("Error: ", ErrorAdminUserIDIsDifferentToSubjectID)
		return false, ErrorAdminUserIDIsDifferentToSubjectID
	}

	if spend.TrackerID != tracker.ID {
		logger.Error("Error: ", ErrorTrackerIDIsDifferntToTrackTrackerID)
		return false, ErrorTrackerIDIsDifferntToTrackTrackerID
	}

	if spend.Currency != tracker.Currency {
		logger.Error("Error: ", ErrorSpendCurrencyDoesNotMatchTrackerCurrency)
		return false, ErrorSpendCurrencyDoesNotMatchTrackerCurrency
	}

	return true, nil
}

// IsValidUpdateSpend ...
func (validator *GoDutchSpendValidator) IsValidUpdateSpend(spend model.Spend,
	logger infrastructure.Logger,
	user model.User,
	existingSpend model.Spend, tracker model.Tracker) (bool, error) {

	if len(spend.Name) <= 0 {
		logger.Error("Error: ", ErrorInvalidName)
		return false, ErrorInvalidName
	}

	if spend.Value.Cmp(decimal.NewFromFloat(0)) == -1 {
		logger.Error("Error: ", ErrorSpendCannotBeLessThanZero)
		return false, ErrorSpendCannotBeLessThanZero
	}

	if spend.Value.Cmp(decimal.NewFromFloat(0)) == 0 {
		logger.Error("Error: ", ErrorSpendCannotBeLessThanZero)
		return false, ErrorSpendCannotBeLessThanZero
	}

	if spend.TrackerID <= 0 {
		logger.Error("Error: ", ErrorInvalidTrackerID)
		return false, ErrorInvalidTrackerID
	}

	if spend.UserID <= 0 {
		logger.Error("Error: ", ErrorInvalidUserID)
		return false, ErrorInvalidUserID
	}

	if len(spend.Currency) <= 0 {
		logger.Error("Error: ", ErrorInvalidCurrency)
		return false, ErrorInvalidCurrency
	}

	if spend.DateCreated.IsZero() {
		logger.Error("Error: ", ErrorInvalidDateCreated)
		return false, ErrorInvalidDateCreated
	}

	if spend.UserID != user.ID {
		logger.Error("Error: ", ErrorAdminUserIDIsDifferentToSubjectID)
		return false, ErrorAdminUserIDIsDifferentToSubjectID
	}

	if spend.ID != existingSpend.ID {
		logger.Error("Error: ", ErrorTheSpendDoesNotExist)
		return false, ErrorTheSpendDoesNotExist
	}

	if spend.Currency != tracker.Currency {
		logger.Error("Error: ", ErrorSpendCurrencyDoesNotMatchTrackerCurrency)
		return false, ErrorSpendCurrencyDoesNotMatchTrackerCurrency
	} 

	return true, nil
}

// IsValidDeleteSpend ...
func (validator *GoDutchSpendValidator) IsValidDeleteSpend(id int64,
	logger infrastructure.Logger,
	user model.User,
	existingSpend model.Spend) (bool, error) {

	if existingSpend.UserID != user.ID {
		logger.Error("Error: ", ErrorAdminUserIDIsDifferentToSubjectID)
		return false, ErrorAdminUserIDIsDifferentToSubjectID
	}

	if id != existingSpend.ID {
		logger.Error("Error: ", ErrorTheSpendDoesNotExist)
		return false, ErrorTheSpendDoesNotExist
	}

	return true, nil
}
