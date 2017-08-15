package trackerservice_test

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
)

var savedTrackerUsers []model.User
var newTracker model.Tracker
var savedTracker model.Tracker
var result bool
var err error
var trackersForUser []model.Tracker
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

func TestCanFindTrackersForUserId(t *testing.T) {

	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
	}

	givenThereAreCleanDependecies()
	givenThereIsAUserWithTheSubAndID("sub", 1)
	givenATrackerAlreadyExistsInTheRepository(tracker)
	whenIGetTheTrackersForAUser("sub")
	thenTheTrackersForTheUserAreReturned(t)
}

func TestCanFindUsersForTracker(t *testing.T) {
	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
	}

	givenThereAreCleanDependecies()
	givenThereIsAUserWithTheSubAndID("sub", 1)
	givenATrackerAlreadyExistsInTheRepository(tracker)
	whenIGetTheUsersForATracker("sub", t)
	thenTheTrackersForTheUserAreReturned(t)
}

func TestCanFindTrackerForTrackerId(t *testing.T) {

	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
	}

	expectedTracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
		ID:             1,
	}

	givenThereAreCleanDependecies()
	givenThereIsAUserWithTheSubAndID("sub", 1)
	givenATrackerAlreadyExistsInTheRepository(tracker)
	whenIGetTheTrackersForAUser("sub")
	thenTheFollowingIsReturned(expectedTracker, t)
}

func TestCanCreateTracker(t *testing.T) {

	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
	}

	expectedTracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
		ID:             1,
	}

	givenThereAreCleanDependecies()
	givenThereIsAUserWithTheSubAndID("sub", 1)
	givenIHaveATracker(tracker)
	whenICreateTheTracker("sub", t)
	thenTheFollowingIsReturned(expectedTracker, t)
}

func TestCanUpdateTracker(t *testing.T) {

	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
	}

	givenThereAreCleanDependecies()
	givenThereIsAUserWithTheSubAndID("sub", 1)
	givenATrackerAlreadyExistsInTheRepository(tracker)

	updatedTracker := model.Tracker{
		Name:           "Tom and Laura 2",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
		ID:             savedTracker.ID,
	}

	whenIUpdateTheTracker("sub", updatedTracker, t)
	thenTheFollowingIsReturned(updatedTracker, t)
}

func TestCanDeleteTracker(t *testing.T) {

	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
	}

	givenThereAreCleanDependecies()
	givenThereIsAUserWithTheSubAndID("sub", 1)
	givenATrackerAlreadyExistsInTheRepository(tracker)
	whenIDeleteTheTracker("sub", savedTracker.ID, t)
	thenTheTrackerIsDeleted(t)
}

func givenThereAreCleanDependecies() {
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

func thenTheTrackersForTheUserAreReturned(t *testing.T) {
	if len(trackersForUser) != 1 {
		t.Fatalf("No trackers were returned for the user")
	}
}

func thenTheUsersForTheTrackerAreReturned(t *testing.T) {
	if len(savedTrackerUsers) != 1 {
		t.Fatalf("No users were returned for the tracker")
	}
}

func whenIGetTheTrackersForAUser(sub string) {
	trackersForUser, err = trackerService.FindByUser(sub)
}

func whenIGetTheUsersForATracker(sub string, t *testing.T) {
	savedTrackerUsers, err = trackerService.FindUsersForTracker(sub, savedTracker.ID)
	if err != nil {
		t.Fatalf("Error was %v", err)
	}
}

func whenIGetTheTrackerByID(sub string, id int) {
	savedTracker, err = trackerService.FindByID(sub, id)
}

func givenThereIsAUserWithTheSubAndID(sub string, id int) {
	userRepository.Insert(model.User{
		AuthenticationID: sub,
		ID:               id,
	})
}

func givenIHaveATracker(tracker model.Tracker) {
	newTracker = tracker
}

func givenATrackerAlreadyExistsInTheRepository(tracker model.Tracker) {
	savedTracker, err = trackerRepository.Insert(tracker)
}

func whenICreateTheTracker(sub string, t *testing.T) {
	savedTracker, err = trackerService.CreateTracker(sub, newTracker)
	if err != nil {
		t.Fatalf("Error was %v", err)
	}
}

func whenIUpdateTheTracker(sub string, tracker model.Tracker, t *testing.T) {
	savedTracker, err = trackerService.UpdateTracker(sub, tracker)
	if err != nil {
		t.Fatalf("Error was %v", err)
	}
}

func whenIDeleteTheTracker(sub string, id int, t *testing.T) {
	result, err = trackerService.DeleteTracker(sub, id)
	if err != nil {
		t.Fatalf("Error was %v", err)
	}
}

func thenTheTrackerIsDeleted(t *testing.T) {
	if result == false {
		t.Fatalf("The tracker was not deleted")
	}
}

func thenTheFollowingIsReturned(expected model.Tracker, t *testing.T) {
	if savedTracker.ID != expected.ID {
		t.Fatalf("Expected %v, got %v", expected.ID, savedTracker.ID)
	}
	if savedTracker.Name != expected.Name {
		t.Fatalf("Expected %v, got %v", expected.Name, savedTracker.Name)
	}

}
