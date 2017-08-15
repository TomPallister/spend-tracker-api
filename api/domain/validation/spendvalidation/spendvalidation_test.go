package spendvalidation_test

import (
	"testing"
	"time"

	"github.com/TomPallister/godutch-api/api/domain/validation/spendvalidation"
	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/shopspring/decimal"
)

var newSpend model.Spend
var existingSpend model.Spend
var err error
var result bool
var logger = infrastructure.NilLogger{}
var newTracker model.Tracker
var newUser model.User
var spendValidator = spendvalidation.NewGoDutchSpendValidator()

func TestCanValidateCreateSpendNoName(t *testing.T) {

	spend := model.Spend{
		Value:       decimal.NewFromFloat(12.99),
		TrackerID:   1,
		UserID:      1,
		Currency:    "£",
		DateCreated: time.Now(),
	}

	givenIHaveASpend(spend)
	whenICallTheCreateSpendValidator()
	thenTheCommandIsRejectedWithError(spendvalidation.ErrorInvalidName, t)
}

func TestCanValidateCreateSpendNoValue(t *testing.T) {
	spend := model.Spend{
		TrackerID:   1,
		Name:        "Cheese",
		UserID:      1,
		Currency:    "£",
		DateCreated: time.Now(),
	}

	givenIHaveASpend(spend)
	whenICallTheCreateSpendValidator()
	thenTheCommandIsRejectedWithError(spendvalidation.ErrorSpendCannotBeLessThanZero, t)
}

func TestCanValidateCreateSpendNoTrackerId(t *testing.T) {
	spend := model.Spend{
		Value:       decimal.NewFromFloat(12.99),
		Name:        "Cheese",
		UserID:      1,
		Currency:    "£",
		DateCreated: time.Now(),
	}

	givenIHaveASpend(spend)
	whenICallTheCreateSpendValidator()
	thenTheCommandIsRejectedWithError(spendvalidation.ErrorInvalidTrackerID, t)
}

func TestCanValidateCreateSpendNoUserId(t *testing.T) {
	spend := model.Spend{
		Value:       decimal.NewFromFloat(12.99),
		TrackerID:   1,
		Name:        "Cheese",
		Currency:    "£",
		DateCreated: time.Now(),
	}

	givenIHaveASpend(spend)
	whenICallTheCreateSpendValidator()
	thenTheCommandIsRejectedWithError(spendvalidation.ErrorInvalidUserID, t)
}

func TestCanValidateCreateSpendNoCurrencyId(t *testing.T) {
	spend := model.Spend{
		Value:       decimal.NewFromFloat(12.99),
		TrackerID:   1,
		Name:        "Cheese",
		UserID:      1,
		DateCreated: time.Now(),
	}

	givenIHaveASpend(spend)
	whenICallTheCreateSpendValidator()
	thenTheCommandIsRejectedWithError(spendvalidation.ErrorInvalidCurrency, t)
}

func TestCanValidateCreateSpendNoDateCreated(t *testing.T) {
	spend := model.Spend{
		Value:     decimal.NewFromFloat(12.99),
		TrackerID: 1,
		Name:      "Cheese",
		UserID:    1,
		Currency:  "£",
	}

	givenIHaveASpend(spend)
	whenICallTheCreateSpendValidator()
	thenTheCommandIsRejectedWithError(spendvalidation.ErrorInvalidDateCreated, t)
}

func TestCanValidateCreateSpendUserDoesntMatchUserID(t *testing.T) {
	spend := model.Spend{
		Value:       decimal.NewFromFloat(12.99),
		TrackerID:   1,
		Name:        "Cheese",
		UserID:      1,
		Currency:    "£",
		DateCreated: time.Now(),
	}

	user := model.User{
		ID: 12,
	}

	givenIHaveAUser(user)
	givenIHaveASpend(spend)
	whenICallTheCreateSpendValidator()
	thenTheCommandIsRejectedWithError(spendvalidation.ErrorAdminUserIDIsDifferentToSubjectID, t)
}

func TestCanValidateCreateSpendTrackerDoesntMatchTrackerID(t *testing.T) {
	spend := model.Spend{
		Value:       decimal.NewFromFloat(12.99),
		TrackerID:   1,
		Name:        "Cheese",
		UserID:      1,
		Currency:    "£",
		DateCreated: time.Now(),
	}

	tracker := model.Tracker{
		ID:       12,
		Currency: "£",
	}

	user := model.User{
		ID: 1,
	}

	givenIHaveATracker(tracker)
	givenIHaveAUser(user)
	givenIHaveASpend(spend)
	whenICallTheCreateSpendValidator()
	thenTheCommandIsRejectedWithError(spendvalidation.ErrorTrackerIDIsDifferntToTrackTrackerID, t)
}

func TestCanValidateCreateSpendCurrencyDoesntMatchTrackerCurrency(t *testing.T) {
	spend := model.Spend{
		Value:       decimal.NewFromFloat(12.99),
		TrackerID:   1,
		Name:        "Cheese",
		UserID:      1,
		Currency:    "£",
		DateCreated: time.Now(),
	}

	tracker := model.Tracker{
		ID:       1,
		Currency: "$",
	}

	user := model.User{
		ID: 1,
	}

	givenIHaveATracker(tracker)
	givenIHaveAUser(user)
	givenIHaveASpend(spend)
	whenICallTheCreateSpendValidator()
	thenTheCommandIsRejectedWithError(spendvalidation.ErrorSpendCurrencyDoesNotMatchTrackerCurrency, t)
}

func TestCanValidateCreateSpend(t *testing.T) {
	spend := model.Spend{
		Value:       decimal.NewFromFloat(12.99),
		TrackerID:   1,
		Name:        "Cheese",
		UserID:      1,
		Currency:    "£",
		DateCreated: time.Now(),
	}

	tracker := model.Tracker{
		ID:       1,
		Currency: "£",
	}

	user := model.User{
		ID: 1,
	}

	givenIHaveATracker(tracker)
	givenIHaveAUser(user)
	givenIHaveASpend(spend)
	whenICallTheCreateSpendValidator()
	thenTheCommandIsAccepted(t)
}

func TestCanValidateDeleteSpend(t *testing.T) {
	spendAlreadyExists := model.Spend{
		Value:       decimal.NewFromFloat(12.99),
		TrackerID:   1,
		Name:        "Cheese",
		UserID:      1,
		Currency:    "£",
		DateCreated: time.Now(),
	}

	spend := model.Spend{
		Value:       decimal.NewFromFloat(118.99),
		TrackerID:   1,
		Name:        "Wine",
		UserID:      1,
		Currency:    "£",
		DateCreated: time.Now(),
	}

	tracker := model.Tracker{
		ID:       1,
		Currency: "£",
	}

	user := model.User{
		ID: 1,
	}

	givenASpendAlreadyExists(spendAlreadyExists)
	givenIHaveATracker(tracker)
	givenIHaveAUser(user)
	givenIHaveASpend(spend)
	whenICallTheUpdateSpendValidator()
	thenTheCommandIsAccepted(t)
}

func givenIHaveATracker(tracker model.Tracker) {
	newTracker = tracker
}

func givenIHaveAUser(user model.User) {
	newUser = user
}

func givenASpendAlreadyExists(spend model.Spend) {
	existingSpend = spend
}

func givenIHaveASpend(spend model.Spend) {
	newSpend = spend
}

func whenICallTheCreateSpendValidator() {
	result, err = spendValidator.IsValidCreateSpend(newSpend, logger, newUser, newTracker)
}

func whenICallTheDeleteSpendValidator(id int) {
	result, err = spendValidator.IsValidDeleteSpend(id, logger, newUser, newSpend)
}

func whenICallTheUpdateSpendValidator() {
	result, err = spendValidator.IsValidUpdateSpend(newSpend, logger, newUser, existingSpend, newTracker)
}

func thenTheCommandIsRejectedWithError(e error, t *testing.T) {
	if err != e {
		t.Fatalf("Error should be %v but was %v", e, err)
	}
}

func thenTheCommandIsAccepted(t *testing.T) {
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
	if result == false {
		t.Fatalf("There result was false and it should be true")
	}
}
