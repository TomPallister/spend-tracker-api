package spendservice_test

import (
	"testing"
	"time"

	"github.com/TomPallister/godutch-api/api/domain/spendservice"
	"github.com/TomPallister/godutch-api/api/domain/spendsummaryservice"
	"github.com/TomPallister/godutch-api/api/domain/trackerservice"
	"github.com/TomPallister/godutch-api/api/domain/transferservice"
	"github.com/TomPallister/godutch-api/api/domain/userservice"
	"github.com/TomPallister/godutch-api/api/domain/validation/spendvalidation"
	"github.com/TomPallister/godutch-api/api/domain/validation/trackervalidation"
	"github.com/TomPallister/godutch-api/api/domain/validation/uservalidation"
	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/repository/spendrepository"
	"github.com/TomPallister/godutch-api/api/repository/spendsummaryrepository"
	"github.com/TomPallister/godutch-api/api/repository/trackerrepository"
	"github.com/TomPallister/godutch-api/api/repository/transferrepository"
	"github.com/TomPallister/godutch-api/api/repository/userrepository"
	"github.com/shopspring/decimal"
)

var logger = infrastructure.ConsoleLogger{}
var emailService = &infrastructure.FakeEmailService{}
var trackerRepository = trackerrepository.NewInMemoryTrackerRepository()
var userRepository = userrepository.NewInMemoryUserRepository()
var spendRepository = spendrepository.NewInMemorySpendRepository()
var transferRepository = transferrepository.NewInMemoryTransferRepository()
var spendSummaryRepository = spendsummaryrepository.NewInMemorySpendSummaryRepository()
var spendSummaryService = spendsummaryservice.
	NewGoDutchSpendSummaryService(spendRepository, spendSummaryRepository, trackerRepository, userRepository)
var transferService = transferservice.
	NewGoDutchTransferService(spendRepository, transferRepository, trackerRepository, userRepository)
var userService = userservice.
	NewGoDutchUserService(userRepository, uservalidation.NewGoDutchUserValidator(), logger, emailService, trackerRepository)
var trackerService = trackerservice.
	NewGoDutchTrackerService(trackerRepository, userService, logger, trackervalidation.NewGoDutchTrackerValidator(), transferService, spendSummaryService)
var spendService = spendservice.
	NewGoDutchSpendService(spendRepository, userService, trackerService, spendvalidation.NewGoDutchSpendValidator(), logger, transferService, spendSummaryService)

var savedUser = model.User{}
var savedTracker = model.Tracker{}
var savedSpend = model.Spend{}
var savedSpends = []model.Spend{}
var newSpend = model.Spend{}
var err error
var result bool

func TestCanCreateSpend(t *testing.T) {

	givenIHaveCleanDependencies()

	user := model.User{
		Name:             "Tom",
		AuthenticationID: "sub",
		DateCreated:      time.Now(),
		EmailAddress:     "email@",
	}

	givenIHaveAUser(user, t)

	tracker := model.Tracker{
		AdminUserID:    savedUser.ID,
		DateCreated:    time.Now(),
		Name:           "Tom and Laura",
		TrackerUserIDs: []int{savedUser.ID},
		Currency:       "£",
	}

	givenIHaveATracker(tracker, t)

	spend := model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Cheese",
		TrackerID:   savedTracker.ID,
		UserID:      savedUser.ID,
		Value:       decimal.NewFromFloat(1.99),
	}

	expected := model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Cheese",
		TrackerID:   savedTracker.ID,
		UserID:      savedUser.ID,
		Value:       decimal.NewFromFloat(1.99),
		ID:          1,
	}

	givenIHaveASpend(spend)
	whenICreateTheSpend(t)
	thenTheFollowingIsReturned(expected, t)
}

func TestCanUpdateSpend(t *testing.T) {
	givenIHaveCleanDependencies()

	user := model.User{
		Name:             "Tom",
		AuthenticationID: "sub",
		DateCreated:      time.Now(),
		EmailAddress:     "email@",
	}

	givenIHaveAUser(user, t)

	tracker := model.Tracker{
		AdminUserID:    savedUser.ID,
		DateCreated:    time.Now(),
		Name:           "Tom and Laura",
		TrackerUserIDs: []int{savedUser.ID},
		Currency:       "£",
	}

	givenIHaveATracker(tracker, t)

	spend := model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Cheese",
		TrackerID:   savedTracker.ID,
		UserID:      savedUser.ID,
		Value:       decimal.NewFromFloat(1.99),
	}

	givenIHaveASpend(spend)
	thenICreateTheSpend(t)

	updatedSpend := model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Wine",
		TrackerID:   savedTracker.ID,
		UserID:      savedUser.ID,
		Value:       decimal.NewFromFloat(25.99),
		ID:          savedSpend.ID,
	}

	whenIUpdateTheSpend(updatedSpend, t)

	expected := model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Wine",
		TrackerID:   savedTracker.ID,
		UserID:      savedUser.ID,
		Value:       decimal.NewFromFloat(25.99),
		ID:          1,
	}

	thenTheFollowingIsReturned(expected, t)

}

func TestCanDeleteSpend(t *testing.T) {
	givenIHaveCleanDependencies()

	user := model.User{
		Name:             "Tom",
		AuthenticationID: "sub",
		DateCreated:      time.Now(),
		EmailAddress:     "email@",
	}

	givenIHaveAUser(user, t)

	tracker := model.Tracker{
		AdminUserID:    savedUser.ID,
		DateCreated:    time.Now(),
		Name:           "Tom and Laura",
		TrackerUserIDs: []int{savedUser.ID},
		Currency:       "£",
	}

	givenIHaveATracker(tracker, t)

	spend := model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Cheese",
		TrackerID:   savedTracker.ID,
		UserID:      savedUser.ID,
		Value:       decimal.NewFromFloat(1.99),
	}

	givenIHaveASpend(spend)
	thenICreateTheSpend(t)
	whenIDeleteTheSpend(t)
	thenTheSpendIsDeleted(t)
}

func TestCanGetSpendsForTracker(t *testing.T) {
	givenIHaveCleanDependencies()

	user := model.User{
		Name:             "Tom",
		AuthenticationID: "sub",
		DateCreated:      time.Now(),
		EmailAddress:     "email@",
	}

	givenIHaveAUser(user, t)

	tracker := model.Tracker{
		AdminUserID:    savedUser.ID,
		DateCreated:    time.Now(),
		Name:           "Tom and Laura",
		TrackerUserIDs: []int{savedUser.ID},
		Currency:       "£",
	}

	givenIHaveATracker(tracker, t)

	spend := model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Cheese",
		TrackerID:   savedTracker.ID,
		UserID:      savedUser.ID,
		Value:       decimal.NewFromFloat(1.99),
	}

	givenIHaveASpend(spend)
	thenICreateTheSpend(t)

	spend = model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Cheese",
		TrackerID:   savedTracker.ID,
		UserID:      savedUser.ID,
		Value:       decimal.NewFromFloat(1.99),
	}

	givenIHaveASpend(spend)
	thenICreateTheSpend(t)
	whenIFindTheSpendsByTrackerID(t)
	thenTheSpendsAreReturnedForTheTracker(t)
}

func whenIFindTheSpendsByTrackerID(t *testing.T) {
	savedSpends, err = spendService.FindByTrackerID(savedUser.AuthenticationID, savedTracker.ID)
	if err != nil {
		t.Fatalf("Error was %v", err)
	}
}

func thenTheSpendsAreReturnedForTheTracker(t *testing.T) {
	if len(savedSpends) != 2 {
		t.Fatalf("Expected %v, got %v", 2, len(savedSpends))
	}
}

func whenIDeleteTheSpend(t *testing.T) {
	result, err = spendService.DeleteSpend(savedUser.AuthenticationID, savedSpend.ID)
	if err != nil {
		t.Fatalf("Error was %v", err)
	}
}

func thenTheSpendIsDeleted(t *testing.T) {
	if result == false {
		t.Fatalf("The tracker was not deleted")
	}
}

func givenIHaveCleanDependencies() {
	logger = infrastructure.ConsoleLogger{}
	trackerRepository = trackerrepository.NewInMemoryTrackerRepository()
	userRepository = userrepository.NewInMemoryUserRepository()
	spendRepository = spendrepository.NewInMemorySpendRepository()
	transferRepository = transferrepository.NewInMemoryTransferRepository()
	spendSummaryRepository = spendsummaryrepository.NewInMemorySpendSummaryRepository()
	spendSummaryService = spendsummaryservice.
		NewGoDutchSpendSummaryService(spendRepository, spendSummaryRepository, trackerRepository, userRepository)
	transferService = transferservice.
		NewGoDutchTransferService(spendRepository, transferRepository, trackerRepository, userRepository)
	userService = userservice.
		NewGoDutchUserService(userRepository, uservalidation.NewGoDutchUserValidator(), logger, emailService, trackerRepository)
	trackerService = trackerservice.
		NewGoDutchTrackerService(trackerRepository, userService, logger, trackervalidation.NewGoDutchTrackerValidator(), transferService, spendSummaryService)
	spendService = spendservice.
		NewGoDutchSpendService(spendRepository, userService, trackerService, spendvalidation.NewGoDutchSpendValidator(), logger, transferService, spendSummaryService)
}

func whenIUpdateTheSpend(spend model.Spend, t *testing.T) {
	savedSpend, err = spendService.UpdateSpend(savedUser.AuthenticationID, spend)
	if err != nil {
		t.Fatalf("There was an error %v", err)

	}
}

func givenIHaveAUser(user model.User, t *testing.T) {
	savedUser, err = userService.CreateUser(user.AuthenticationID, user)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func givenIHaveATracker(tracker model.Tracker, t *testing.T) {
	savedTracker, err = trackerService.CreateTracker(savedUser.AuthenticationID, tracker)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func givenIHaveASpend(spend model.Spend) {
	newSpend = spend
}

func whenICreateTheSpend(t *testing.T) {
	savedSpend, err = spendService.CreateSpend(savedUser.AuthenticationID, newSpend)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func thenICreateTheSpend(t *testing.T) {
	savedSpend, err = spendService.CreateSpend(savedUser.AuthenticationID, newSpend)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func thenTheFollowingIsReturned(expected model.Spend, t *testing.T) {
	if savedSpend.ID != expected.ID {
		t.Fatalf("Expected %v, got %v", expected.ID, savedSpend.ID)
	}
	if savedSpend.Name != expected.Name {
		t.Fatalf("Expected %v, got %v", expected.Name, savedSpend.Name)
	}
	if savedSpend.Currency != expected.Currency {
		t.Fatalf("Expected %v, got %v", expected.Currency, savedSpend.Currency)
	}
	if savedSpend.Value.Cmp(expected.Value) != 0 {
		t.Fatalf("Expected %v, got %v", expected.Value, savedSpend.Value)
	}
}
