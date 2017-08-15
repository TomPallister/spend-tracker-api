package userservice_test

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
	"github.com/TomPallister/godutch-api/api/infrastructure/encryption"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/repository/spendrepository"
	"github.com/TomPallister/godutch-api/api/repository/spendsummaryrepository"
	"github.com/TomPallister/godutch-api/api/repository/trackerrepository"
	"github.com/TomPallister/godutch-api/api/repository/transferrepository"
	"github.com/TomPallister/godutch-api/api/repository/userrepository"
)

var savedTracker model.Tracker
var savedUser model.User
var invitedUser model.User
var err error
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

func TestCanGetUser(t *testing.T) {
	givenThereAreCleanDependencies()
	givenThereIsAUserInTheRepository(model.User{ID: 0, Name: "Tom", AuthenticationID: "asdf"}, t)
	whenIGetAUserWithTheID("asdf", 1, t)
	thenTheFollowingIsReturned(model.User{ID: 1, Name: "Tom", AuthenticationID: "asdf"}, t)
}

func TestCanGetUserBySub(t *testing.T) {
	givenThereAreCleanDependencies()
	givenThereIsAUserInTheRepository(model.User{ID: 0, Name: "Tom", AuthenticationID: "asdf1"}, t)
	whenIGetAUserWithTheSub("asdf1", t)
	thenTheFollowingIsReturned(model.User{ID: 1, Name: "Tom", AuthenticationID: "asdf1"}, t)
}

func TestCanGetUserById(t *testing.T) {
	givenThereAreCleanDependencies()
	givenThereIsAUserInTheRepository(model.User{ID: 0, Name: "Tom", AuthenticationID: "asdf1"}, t)
	whenIGetAUserByIDForATracker(savedUser.ID, t)
	thenTheFollowingIsReturned(model.User{ID: 1, Name: "Tom", AuthenticationID: "asdf1"}, t)
}

func TestCanCreateUser(t *testing.T) {
	givenThereAreCleanDependencies()
	givenIHaveCreatedAUser("asd2", model.User{Name: "Tom", AuthenticationID: "asd2", EmailAddress: "email@", DateCreated: time.Now()}, t)
	whenIGetAUserWithTheID("asd2", 1, t)
	thenTheFollowingIsReturned(model.User{ID: 1, Name: "Tom", AuthenticationID: "asd2"}, t)
}

func TestCanInviteUser(t *testing.T) {
	givenThereAreCleanDependencies()
	givenIHaveCreatedAUser("asd2", model.User{Name: "Tom", AuthenticationID: "asd2", EmailAddress: "email@", DateCreated: time.Now()}, t)
	tracker := model.Tracker{
		AdminUserID:    savedUser.ID,
		Currency:       "£",
		Name:           "test",
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{savedUser.ID},
	}
	givenIHaveCreatedATracker(t, tracker)
	inviteUser := model.InviteUser{
		EmailAddress: "thomasgardham@googlemail.com",
		TrackerID:    savedTracker.ID,
	}
	whenIInviteTheUserWith(inviteUser, t)
	thenTheUserIsInvited(inviteUser.EmailAddress, t)
}

func TestCanAcceptInviteUser(t *testing.T) {
	givenThereAreCleanDependencies()
	givenIHaveCreatedAUser("asd2", model.User{Name: "Tom", AuthenticationID: "asd2", EmailAddress: "email@", DateCreated: time.Now()}, t)
	tracker := model.Tracker{
		AdminUserID:    savedUser.ID,
		Currency:       "£",
		Name:           "test",
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{savedUser.ID},
	}
	givenIHaveCreatedATracker(t, tracker)
	inviteUser := model.InviteUser{
		EmailAddress: "thomasgardham@googlemail.com",
		TrackerID:    savedTracker.ID,
	}
	whenIInviteTheUserWith(inviteUser, t)
	thenTheUserIsInvited(inviteUser.EmailAddress, t)

	userAcceptsInvite := model.User{
		AuthenticationID: "anysub",
		DateCreated:      time.Now(),
		EmailAddress:     inviteUser.EmailAddress,
		Name:             "Tom",
	}

	cryptoEmailAddress := encryption.Encrypt(inviteUser.EmailAddress)
	whenTheUserAcceptsTheInvite(t, userAcceptsInvite.AuthenticationID, userAcceptsInvite, cryptoEmailAddress)
	thenTheUserIsCreated(t)
}

func thenTheUserIsCreated(t *testing.T) {
	if savedUser.ID != invitedUser.ID {
		t.Fatalf("Expected %v, got %v", invitedUser.ID, savedUser.ID)
	}
}

func whenTheUserAcceptsTheInvite(t *testing.T, sub string, user model.User, cryptoEmail string) {
	invitedUser, err = userService.AcceptInvite(sub, user, cryptoEmail)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func givenIHaveCreatedATracker(t *testing.T, tracker model.Tracker) {
	savedTracker, err = trackerRepository.Insert(tracker)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func thenTheUserIsInvited(email string, t *testing.T) {
	if savedUser.ID < 1 {
		t.Fatalf("Invited user ID was %v", invitedUser)
	}
	if savedUser.EmailAddress != email {
		t.Fatalf("Expected %v, got %v", email, savedUser.EmailAddress)
	}
}

func whenIInviteTheUserWith(inviteUserModel model.InviteUser, t *testing.T) {
	savedUser, err = userService.InviteUser(savedUser.AuthenticationID, inviteUserModel, "some root url")
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func givenThereAreCleanDependencies() {
	emailService = &infrastructure.FakeEmailService{}
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

func givenIHaveCreatedAUser(sub string, user model.User, t *testing.T) {
	savedUser, err = userService.CreateUser(sub, user)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func givenThereIsAUserInTheRepository(newUser model.User, t *testing.T) {
	savedUser, err = userRepository.Insert(newUser)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func whenIGetAUserWithTheID(sub string, id int, t *testing.T) {
	savedUser, err = userService.FindByID(sub, id)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func whenIGetAUserWithTheSub(sub string, t *testing.T) {
	savedUser, err = userService.FindBySub(sub)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func whenIGetAUserByIDForATracker(id int, t *testing.T) {
	savedUser, err = userService.FindByIDForTracker(id)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func thenTheFollowingIsReturned(expectedUser model.User, t *testing.T) {
	if savedUser.ID != expectedUser.ID {
		t.Fatalf("Expected %v, got %v", expectedUser.ID, savedUser.ID)
	}
	if savedUser.Name != expectedUser.Name {
		t.Fatalf("Expected %v, got %v", expectedUser.Name, savedUser.Name)
	}
	if savedUser.AuthenticationID != expectedUser.AuthenticationID {
		t.Fatalf("Expected %v, got %v", expectedUser.AuthenticationID, savedUser.AuthenticationID)
	}
}
