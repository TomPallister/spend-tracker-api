package trackerrepository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/repository"
	"github.com/TomPallister/godutch-api/api/repository/trackerrepository"
	"github.com/TomPallister/godutch-api/api/repository/userrepository"
	"github.com/nu7hatch/gouuid"
)

var domainUser model.User
var domainTracker model.Tracker
var domainTrackers []model.Tracker
var err error
var userRepository userrepository.UserRepository
var trackerRepository trackerrepository.TrackerRepository
var domainUserID int64
var db *sql.DB
var domainTrackerID int64
var result bool

func TestCanUpdateTracker(t *testing.T) {
	givenIHaveCleanDependencies(t)
	givenThereIsAUser(t)

	tracker := model.Tracker{
		AdminUserID:    domainUser.ID,
		Currency:       "£",
		DateCreated:    time.Now(),
		Name:           "weezey spends",
		TrackerUserIDs: []int64{domainUser.ID},
	}

	givenIHaveATracker(tracker)
	givenISaveTheTracker(t)

	userOne := domainUser

	givenThereIsAUser(t)

	userTwo := domainUser

	updatedTracker := model.Tracker{
		ID:             domainTracker.ID,
		AdminUserID:    userOne.ID,
		Currency:       "$",
		DateCreated:    domainTracker.DateCreated,
		Name:           "NewName",
		TrackerUserIDs: []int64{userOne.ID, userTwo.ID},
	}

	whenIUpdateTheTracker(updatedTracker, t)
	thenTheTrackerIsUpdated(updatedTracker, t)
}

func TestCanDeleteTrackerByID(t *testing.T) {
	givenIHaveCleanDependencies(t)
	givenThereIsAUser(t)

	tracker := model.Tracker{
		AdminUserID:    domainUser.ID,
		Currency:       "£",
		DateCreated:    time.Now(),
		Name:           "weezey spends",
		TrackerUserIDs: []int64{domainUser.ID},
	}

	givenIHaveATracker(tracker)
	givenISaveTheTracker(t)
	whenIDeleteTheTracker(t)
	thenTheTrackerIsDeleted(t)
}

func TestCanGetTrackersByUserID(t *testing.T) {
	givenIHaveCleanDependencies(t)
	givenThereIsAUser(t)

	tracker := model.Tracker{
		AdminUserID:    domainUser.ID,
		Currency:       "£",
		DateCreated:    time.Now(),
		Name:           "weezey spends",
		TrackerUserIDs: []int64{domainUser.ID},
	}

	givenIHaveATracker(tracker)
	givenISaveTheTracker(t)
	whenIGetTheTrackersForAUser(t)

	expected := []model.Tracker{tracker}

	thenTheFollowingTrackersAreReturned(expected, t)
}

func TestCanGetTrackerByID(t *testing.T) {
	givenIHaveCleanDependencies(t)
	givenThereIsAUser(t)

	tracker := model.Tracker{
		AdminUserID:    domainUser.ID,
		Currency:       "£",
		DateCreated:    time.Now(),
		Name:           "weezey spends",
		TrackerUserIDs: []int64{domainUser.ID},
	}

	givenIHaveATracker(tracker)
	givenISaveTheTracker(t)
	whenIGetTheTracker(t)
	thenTheFollowingIsReturned(tracker, t)
}

func TestCanInsertTracker(t *testing.T) {

	givenIHaveCleanDependencies(t)
	givenThereIsAUser(t)

	tracker := model.Tracker{
		AdminUserID:    domainUser.ID,
		Currency:       "£",
		DateCreated:    time.Now(),
		Name:           "weezey spends",
		TrackerUserIDs: []int64{domainUser.ID},
	}

	givenIHaveATracker(tracker)
	whenISaveTheTracker(t)
	thenTheTrackerIsSaved(t)
}

func whenIUpdateTheTracker(updated model.Tracker, t *testing.T) {
	domainTracker, err = trackerRepository.Update(updated.ID, updated)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}

func thenTheTrackerIsUpdated(expected model.Tracker, t *testing.T) {
	domainTracker, err = trackerRepository.GetByID(domainTracker.ID)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if domainTracker.AdminUserID != expected.AdminUserID {
		t.Fatalf("Expected: %v, Received: %v", expected.AdminUserID, domainTracker.AdminUserID)
	}
	if domainTracker.Currency != expected.Currency {
		t.Fatalf("Expected: %v, Received: %v", expected.Currency, domainTracker.Currency)
	}

	if domainTracker.Name != expected.Name {
		t.Fatalf("Expected: %v, Received: %v", expected.Name, domainTracker.Name)
	}

	for i := 0; i < len(domainTracker.TrackerUserIDs); i++ {
		if domainTracker.TrackerUserIDs[i] != expected.TrackerUserIDs[i] {
			t.Fatalf("Expected: %v, Received: %v", expected.TrackerUserIDs[i], domainTracker.TrackerUserIDs[i])
		}
	}
}

func whenIDeleteTheTracker(t *testing.T) {
	result, err = trackerRepository.Delete(domainTracker.ID)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}

func thenTheTrackerIsDeleted(t *testing.T) {
	if result == false {
		t.Fatalf("Tracker wasnt deleted")
	}
}

func thenTheFollowingTrackersAreReturned(expected []model.Tracker, t *testing.T) {

	if len(domainTrackers) == 0 {
		t.Fatalf("There were no trackers")
	}

	for i := 0; i < len(domainTrackers); i++ {

		if domainTrackers[i].AdminUserID != expected[i].AdminUserID {
			t.Fatalf("Expected: %v, Received: %v", expected[i].AdminUserID, domainTrackers[i].AdminUserID)
		}
		if domainTrackers[i].Currency != expected[i].Currency {
			t.Fatalf("Expected: %v, Received: %v", expected[i].Currency, domainTrackers[i].Currency)
		}

		if domainTrackers[i].Name != expected[i].Name {
			t.Fatalf("Expected: %v, Received: %v", expected[i].Name, domainTrackers[i].Name)
		}

		for y := 0; y < len(domainTrackers[i].TrackerUserIDs); y++ {
			if domainTrackers[i].TrackerUserIDs[y] != expected[i].TrackerUserIDs[y] {
				t.Fatalf("Expected: %v, Received: %v", expected[i].TrackerUserIDs[y], domainTrackers[i].TrackerUserIDs[y])
			}
		}
	}
}

func thenTheFollowingIsReturned(expected model.Tracker, t *testing.T) {
	if domainTracker.AdminUserID != expected.AdminUserID {
		t.Fatalf("Expected: %v, Received: %v", expected.AdminUserID, domainTracker.AdminUserID)
	}
	if domainTracker.Currency != expected.Currency {
		t.Fatalf("Expected: %v, Received: %v", expected.Currency, domainTracker.Currency)
	}

	if domainTracker.Name != expected.Name {
		t.Fatalf("Expected: %v, Received: %v", expected.Name, domainTracker.Name)
	}

	for i := 0; i < len(domainTracker.TrackerUserIDs); i++ {
		if domainTracker.TrackerUserIDs[i] != expected.TrackerUserIDs[i] {
			t.Fatalf("Expected: %v, Received: %v", expected.TrackerUserIDs[i], domainTracker.TrackerUserIDs[i])
		}
	}
}

func whenIGetTheTrackersForAUser(t *testing.T) {
	domainTrackers, err = trackerRepository.GetForUserID(domainUser.ID)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}

func whenIGetTheTracker(t *testing.T) {
	domainTracker, err = trackerRepository.GetByID(domainTracker.ID)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}

func thenTheTrackerIsSaved(t *testing.T) {
	if domainTracker.ID == 0 {
		t.Fatalf("The domain tracker ID is 0")
	}
}

func whenISaveTheTracker(t *testing.T) {
	domainTracker, err = trackerRepository.Insert(domainTracker)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}

func givenISaveTheTracker(t *testing.T) {
	domainTracker, err = trackerRepository.Insert(domainTracker)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}

func givenIHaveATracker(tracker model.Tracker) {
	domainTracker = tracker
}

func givenThereIsAUser(t *testing.T) {

	u, _ := uuid.NewV4()

	user := model.User{
		Name:             "Laura",
		AuthenticationID: "21312fsdf" + u.String(),
		DateCreated:      time.Now(),
		EmailAddress:     "laura" + u.String() + "@laz.com",
	}

	givenIHaveCleanDependencies(t)
	givenIHaveADomainUser(user)
	whenISaveTheUser(t)
	thenTheUserIsSaved(t)
}

func givenIHaveCleanDependencies(t *testing.T) {
	db, err = repository.NewDB("postgres://godutch:password@localhost/godutch?sslmode=disable")
	userRepository = userrepository.NewPostgresUserRepository(infrastructure.ConsoleLogger{}, db)
	trackerRepository = trackerrepository.NewPostgresTrackerRepository(infrastructure.ConsoleLogger{}, db)
}

func givenIHaveADomainUser(user model.User) {
	domainUser = user
}

func whenISaveTheUser(t *testing.T) {
	domainUser, err = userRepository.Insert(domainUser)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func thenTheUserIsSaved(t *testing.T) {
	if domainUser.ID <= 0 {
		t.Fatalf("The user ID was %v", domainUser.ID)
	}
}
