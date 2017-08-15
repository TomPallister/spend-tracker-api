package spendsummaryservice_test

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

var spendUserOne = model.Spend{}
var spendUserTwo = model.Spend{}
var savedUserOne = model.User{}
var savedUserTwo = model.User{}
var savedTracker = model.Tracker{}
var savedSpend = model.Spend{}
var savedSpends = []model.Spend{}
var newSpend = model.Spend{}
var err error
var result bool
var savedSpendSummaries = []model.SpendSummary{}
var emailService = &infrastructure.FakeEmailService{}
var logger = infrastructure.ConsoleLogger{}
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

func TestCanCreateBasicSpendSummaries(t *testing.T) {
	givenIHaveCleanDependencies()
	givenUsersAndTrackerHaveBeenCreated(t)

	spendUserOne := model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Cheese",
		TrackerID:   savedTracker.ID,
		UserID:      savedUserOne.ID,
		Value:       decimal.NewFromFloat(56),
	}

	givenIHaveASpend(savedUserOne.AuthenticationID, spendUserOne, t)

	spendUserTwo := model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Cheese",
		TrackerID:   savedTracker.ID,
		UserID:      savedUserTwo.ID,
		Value:       decimal.NewFromFloat(1),
	}

	givenIHaveASpend(savedUserTwo.AuthenticationID, spendUserTwo, t)
	whenICreateTheSpendSummaries(t)

	expectedSpendSummaries := []model.SpendSummary{
		model.SpendSummary{
			UserID:       savedUserOne.ID,
			EmailAddress: savedUserOne.EmailAddress,
			TrackerID:    savedTracker.ID,
			Value:        spendUserOne.Value,
			Currency:     savedTracker.Currency,
		},
		model.SpendSummary{
			UserID:       savedUserTwo.ID,
			EmailAddress: savedUserTwo.EmailAddress,
			TrackerID:    savedTracker.ID,
			Value:        spendUserTwo.Value,
			Currency:     savedTracker.Currency,
		},
	}

	thenTheFollowingSpendSummariesreReturned(expectedSpendSummaries, t)
}

func TestCanUpdateSpendSummariesForTrackerID(t *testing.T) {
	givenIHaveCleanDependencies()
	givenThereAreSpendSummariesForATracker(t)

	spendTwoUserTwo := model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Meat",
		TrackerID:   savedTracker.ID,
		UserID:      savedUserTwo.ID,
		Value:       decimal.NewFromFloat(10),
	}

	givenIHaveASpend(savedUserTwo.AuthenticationID, spendTwoUserTwo, t)

	whenIUpdateTheSpendSummaries(t)

	expectedSpendSummaries := []model.SpendSummary{
		model.SpendSummary{
			UserID:       savedUserOne.ID,
			EmailAddress: savedUserOne.EmailAddress,
			TrackerID:    savedTracker.ID,
			Value:        spendUserOne.Value,
			Currency:     savedTracker.Currency,
		},
		model.SpendSummary{
			UserID:       savedUserTwo.ID,
			EmailAddress: savedUserTwo.EmailAddress,
			TrackerID:    savedTracker.ID,
			Value:        decimal.NewFromFloat(11),
			Currency:     savedTracker.Currency,
		},
	}

	thenTheFollowingSpendSummariesreReturned(expectedSpendSummaries, t)

}

func TestCanFindSpendSummariesForTrackerID(t *testing.T) {
	givenIHaveCleanDependencies()

	givenThereAreSpendSummariesForATracker(t)

	whenIGetTheSpendSummaries(savedUserOne.AuthenticationID, t)

	expectedSpendSummaries := []model.SpendSummary{
		model.SpendSummary{
			UserID:       savedUserOne.ID,
			EmailAddress: savedUserOne.EmailAddress,
			TrackerID:    savedTracker.ID,
			Value:        spendUserOne.Value,
			Currency:     savedTracker.Currency,
		},
		model.SpendSummary{
			UserID:       savedUserTwo.ID,
			EmailAddress: savedUserTwo.EmailAddress,
			TrackerID:    savedTracker.ID,
			Value:        spendUserTwo.Value,
			Currency:     savedTracker.Currency,
		},
	}

	thenTheFollowingSpendSummariesreReturned(expectedSpendSummaries, t)
}

func TestCanDeleteSpendSummariesForTrackerID(t *testing.T) {
	givenIHaveCleanDependencies()
	givenThereAreSpendSummariesForATracker(t)
	whenIDeleteTheSpendSummaries(t)
	thenNoSpendSummariesAreReturned(t)
}

func whenIDeleteTheSpendSummaries(t *testing.T) {
	result, err = spendSummaryService.DeleteSpendSummaries(savedTracker.ID)
	if err != nil {
		t.Fatalf("There was an error: %v", err)
	}
}

func thenNoSpendSummariesAreReturned(t *testing.T) {
	if result == false {
		t.Fatalf("Transfers were not deleted")
	}
}

func whenIGetTheSpendSummaries(sub string, t *testing.T) {
	savedSpendSummaries, err = spendSummaryService.FindSpendSummariesForTrackerID(sub, savedTracker.ID)
	if err != nil {
		t.Fatalf("There was an error: %v", err)
	}
}

func givenThereAreSpendSummariesForATracker(t *testing.T) {

	givenUsersAndTrackerHaveBeenCreated(t)

	spendUserOne = model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Cheese",
		TrackerID:   savedTracker.ID,
		UserID:      savedUserOne.ID,
		Value:       decimal.NewFromFloat(56),
	}

	givenIHaveASpend(savedUserOne.AuthenticationID, spendUserOne, t)

	spendUserTwo = model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Cheese",
		TrackerID:   savedTracker.ID,
		UserID:      savedUserTwo.ID,
		Value:       decimal.NewFromFloat(1),
	}

	givenIHaveASpend(savedUserTwo.AuthenticationID, spendUserTwo, t)
	whenICreateTheSpendSummaries(t)

	expectedSpendSummaries := []model.SpendSummary{
		model.SpendSummary{
			UserID:       savedUserOne.ID,
			EmailAddress: savedUserOne.EmailAddress,
			TrackerID:    savedTracker.ID,
			Value:        spendUserOne.Value,
			Currency:     savedTracker.Currency,
		},
		model.SpendSummary{
			UserID:       savedUserTwo.ID,
			EmailAddress: savedUserTwo.EmailAddress,
			TrackerID:    savedTracker.ID,
			Value:        spendUserTwo.Value,
			Currency:     savedTracker.Currency,
		},
	}

	thenTheFollowingSpendSummariesreReturned(expectedSpendSummaries, t)
}

func thenTheFollowingSpendSummariesreReturned(expectedSpendSummary []model.SpendSummary, t *testing.T) {
	if len(savedSpendSummaries) <= 0 {
		t.Fatalf("There were not saved spend summaries")
	}
	for i := 0; i < len(savedSpendSummaries); i++ {
		if savedSpendSummaries[i].Currency != expectedSpendSummary[i].Currency {
			t.Fatalf("Expected %v, got %v", expectedSpendSummary[i].Currency, savedSpendSummaries[i].Currency)
		}
		if savedSpendSummaries[i].TrackerID != expectedSpendSummary[i].TrackerID {
			t.Fatalf("Expected %v, got %v", expectedSpendSummary[i].TrackerID, savedSpendSummaries[i].TrackerID)
		}
		if savedSpendSummaries[i].EmailAddress != expectedSpendSummary[i].EmailAddress {
			t.Fatalf("Expected %v, got %v", expectedSpendSummary[i].EmailAddress, savedSpendSummaries[i].EmailAddress)
		}
		if savedSpendSummaries[i].UserID != expectedSpendSummary[i].UserID {
			t.Fatalf("Expected %v, got %v", expectedSpendSummary[i].UserID, savedSpendSummaries[i].UserID)
		}
		if savedSpendSummaries[i].Value.Cmp(expectedSpendSummary[i].Value) != 0 {
			t.Fatalf("Expected %v, got %v", expectedSpendSummary[i].Value, savedSpendSummaries[i].Value)
		}
	}
}

func whenICreateTheSpendSummaries(t *testing.T) {
	savedSpendSummaries, err = spendSummaryService.UpsertSpendSummaries(savedTracker.ID)
	if err != nil {
		t.Fatalf("There was an error %v", err)
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

func givenUsersAndTrackerHaveBeenCreated(t *testing.T) {

	userOne := model.User{
		Name:             "Tom",
		AuthenticationID: "sub",
		DateCreated:      time.Now(),
		EmailAddress:     "email@tom.com",
	}

	givenIHaveUserOne(userOne, t)

	userTwo := model.User{
		Name:             "Laura",
		AuthenticationID: "sub123",
		DateCreated:      time.Now(),
		EmailAddress:     "email@laz.com",
	}

	givenIHaveUserTwo(userTwo, t)

	tracker := model.Tracker{
		AdminUserID:    savedUserOne.ID,
		DateCreated:    time.Now(),
		Name:           "Tom and Laura",
		TrackerUserIDs: []int{savedUserOne.ID, savedUserTwo.ID},
		Currency:       "£",
	}

	givenIHaveATracker(tracker, t)

}

func givenIHaveUserOne(user model.User, t *testing.T) {
	savedUserOne, err = userService.CreateUser(user.AuthenticationID, user)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func givenIHaveUserTwo(user model.User, t *testing.T) {
	savedUserTwo, err = userService.CreateUser(user.AuthenticationID, user)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func givenIHaveATracker(tracker model.Tracker, t *testing.T) {
	savedTracker, err = trackerService.CreateTracker(savedUserOne.AuthenticationID, tracker)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func givenIHaveASpend(sub string, spend model.Spend, t *testing.T) {
	savedSpend, err = spendService.CreateSpend(sub, spend)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func whenIUpdateTheSpendSummaries(t *testing.T) {
	savedSpendSummaries, err = spendSummaryService.UpsertSpendSummaries(savedTracker.ID)
}
