package transferservice_test

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

var savedUserOne = model.User{}
var savedUserTwo = model.User{}
var savedTracker = model.Tracker{}
var savedSpend = model.Spend{}
var savedSpends = []model.Spend{}
var newSpend = model.Spend{}
var err error
var result bool
var savedTransfers = []model.Transfer{}
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

func TestCanCreateBasicTransfers(t *testing.T) {
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
	whenICreateTheTransfers(t)

	expectedTransfers := []model.Transfer{
		model.Transfer{
			FromUserID:       savedUserTwo.ID,
			FromEmailAddress: savedUserTwo.EmailAddress,
			ToUserID:         savedUserOne.ID,
			ToEmailAddress:   savedUserOne.EmailAddress,
			Value:            decimal.NewFromFloat(27.5),
			Currency:         savedTracker.Currency,
			TrackerID:        savedTracker.ID,
		},
	}

	thenTheFollowingTransfersAreReturned(expectedTransfers, t)
}

func TestCanUpdateTransfersForTrackerID(t *testing.T) {
	givenIHaveCleanDependencies()
	givenThereAreTransferForATracker(t)

	spendTwoUserTwo := model.Spend{
		Currency:    "£",
		DateCreated: time.Now(),
		Name:        "Meat",
		TrackerID:   savedTracker.ID,
		UserID:      savedUserTwo.ID,
		Value:       decimal.NewFromFloat(10),
	}

	givenIHaveASpend(savedUserTwo.AuthenticationID, spendTwoUserTwo, t)

	whenIUpdateTheTransfers(t)

	expectedTransfers := []model.Transfer{
		model.Transfer{
			FromUserID:       savedUserTwo.ID,
			FromEmailAddress: savedUserTwo.EmailAddress,
			ToUserID:         savedUserOne.ID,
			ToEmailAddress:   savedUserOne.EmailAddress,
			Value:            decimal.NewFromFloat(22.5),
			Currency:         savedTracker.Currency,
			TrackerID:        savedTracker.ID,
		},
	}

	thenTheFollowingTransfersAreReturned(expectedTransfers, t)

}

func TestCanFindTransfersForTrackerID(t *testing.T) {
	givenIHaveCleanDependencies()

	givenThereAreTransferForATracker(t)

	whenIGetTheTransfers(savedUserOne.AuthenticationID, t)

	expectedTransfers := []model.Transfer{
		model.Transfer{
			FromUserID:       savedUserTwo.ID,
			FromEmailAddress: savedUserTwo.EmailAddress,
			ToUserID:         savedUserOne.ID,
			ToEmailAddress:   savedUserOne.EmailAddress,
			Value:            decimal.NewFromFloat(27.5),
			Currency:         savedTracker.Currency,
			TrackerID:        savedTracker.ID,
		},
	}

	thenTheFollowingTransfersAreReturned(expectedTransfers, t)
}

func TestCanDeleteTransfersForTrackerID(t *testing.T) {
	givenIHaveCleanDependencies()
	givenThereAreTransferForATracker(t)
	whenIDeleteTheTransfers(t)
	thenNoTransfersAreReturned(t)
}

func whenIDeleteTheTransfers(t *testing.T) {
	result, err = transferService.DeleteTransfers(savedTracker.ID)
	if err != nil {
		t.Fatalf("There was an error: %v", err)
	}
}

func thenNoTransfersAreReturned(t *testing.T) {
	if result == false {
		t.Fatalf("Transfers were not deleted")
	}
}

func whenIGetTheTransfers(sub string, t *testing.T) {
	savedTransfers, err = transferService.FindTransfersForTrackerID(sub, savedTracker.ID)
	if err != nil {
		t.Fatalf("There was an error: %v", err)
	}
}

func givenThereAreTransferForATracker(t *testing.T) {

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
	whenICreateTheTransfers(t)

	expectedTransfers := []model.Transfer{
		model.Transfer{
			FromUserID:       savedUserTwo.ID,
			FromEmailAddress: savedUserTwo.EmailAddress,
			ToUserID:         savedUserOne.ID,
			ToEmailAddress:   savedUserOne.EmailAddress,
			Value:            decimal.NewFromFloat(27.5),
			Currency:         savedTracker.Currency,
			TrackerID:        savedTracker.ID,
		},
	}

	thenTheFollowingTransfersAreReturned(expectedTransfers, t)
}

func thenTheFollowingTransfersAreReturned(expectedTransfers []model.Transfer, t *testing.T) {
	if len(savedTransfers) <= 0 {
		t.Fatalf("There were not saved transfers")
	}
	for i := 0; i < len(savedTransfers); i++ {
		if savedTransfers[i].Currency != expectedTransfers[i].Currency {
			t.Fatalf("Expected %v, got %v", expectedTransfers[i].Currency, savedTransfers[i].Currency)
		}
		if savedTransfers[i].TrackerID != expectedTransfers[i].TrackerID {
			t.Fatalf("Expected %v, got %v", expectedTransfers[i].TrackerID, savedTransfers[i].TrackerID)
		}
		if savedTransfers[i].FromEmailAddress != expectedTransfers[i].FromEmailAddress {
			t.Fatalf("Expected %v, got %v", expectedTransfers[i].FromEmailAddress, savedTransfers[i].FromEmailAddress)
		}
		if savedTransfers[i].FromUserID != expectedTransfers[i].FromUserID {
			t.Fatalf("Expected %v, got %v", expectedTransfers[i].FromUserID, savedTransfers[i].FromUserID)
		}
		if savedTransfers[i].ToEmailAddress != expectedTransfers[i].ToEmailAddress {
			t.Fatalf("Expected %v, got %v", expectedTransfers[i].ToEmailAddress, savedTransfers[i].ToEmailAddress)
		}
		if savedTransfers[i].ToUserID != expectedTransfers[i].ToUserID {
			t.Fatalf("Expected %v, got %v", expectedTransfers[i].ToUserID, savedTransfers[i].ToUserID)
		}
		if savedTransfers[i].Value.Cmp(expectedTransfers[i].Value) != 0 {
			t.Fatalf("Expected %v, got %v", expectedTransfers[i].Value, savedTransfers[i].Value)
		}
	}
}

func whenICreateTheTransfers(t *testing.T) {
	savedTransfers, err = transferService.UpsertTransfers(savedTracker.ID)
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

func whenIUpdateTheTransfers(t *testing.T) {
	savedTransfers, err = transferService.UpsertTransfers(savedTracker.ID)
}
