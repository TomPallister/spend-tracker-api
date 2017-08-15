package spendrepository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/repository"
	"github.com/TomPallister/godutch-api/api/repository/spendrepository"
	"github.com/TomPallister/godutch-api/api/repository/trackerrepository"
	"github.com/TomPallister/godutch-api/api/repository/userrepository"
	"github.com/nu7hatch/gouuid"
	"github.com/shopspring/decimal"
)

var domainUser model.User
var domainTracker model.Tracker
var domainTrackers []model.Tracker
var userOne model.User
var userTwo model.User
var err error
var userRepository userrepository.UserRepository
var trackerRepository trackerrepository.TrackerRepository
var spendRepository spendrepository.SpendRepository
var domainUserID int64
var db *sql.DB
var domainTrackerID int64
var result bool
var domainSpend model.Spend
var domainSpends []model.Spend

func TestCanUpdateSpend(t *testing.T) {

	givenIHaveCleanDependencies(t)
	givenThereIsATrackerAndTwoUsers(t)

	spend := model.Spend{
		Currency:    domainTracker.Currency,
		DateCreated: time.Now(),
		Name:        "Cheese & Wine",
		TrackerID:   domainTracker.ID,
		UserID:      userOne.ID,
		Value:       decimal.NewFromFloat(10.99),
	}

	givenIHaveASpend(spend)
	givenISaveTheSpend(t)

	updateSpend := model.Spend{
		ID:          domainSpend.ID,
		Currency:    "$",
		DateCreated: domainSpend.DateCreated,
		Name:        "Ikea MALM",
		TrackerID:   domainTracker.ID,
		UserID:      userTwo.ID,
		Value:       decimal.NewFromFloat(876.99),
	}

	whenIUpdateTheSpend(updateSpend, t)
	thenTheSpendIsUpdated(updateSpend, t)
}

func TestCanDeleteByID(t *testing.T) {
	givenIHaveCleanDependencies(t)
	givenThereIsATrackerAndTwoUsers(t)

	spend := model.Spend{
		Currency:    domainTracker.Currency,
		DateCreated: time.Now(),
		Name:        "Cheese & Wine",
		TrackerID:   domainTracker.ID,
		UserID:      userOne.ID,
		Value:       decimal.NewFromFloat(10.99),
	}

	givenIHaveASpend(spend)
	givenISaveTheSpend(t)
	whenIDeleteTheSpend(t)
	thenTheSpendIsDeleted(t)
}

func TestCanGetByID(t *testing.T) {

	givenIHaveCleanDependencies(t)
	givenThereIsATrackerAndTwoUsers(t)

	spend := model.Spend{
		Currency:    domainTracker.Currency,
		DateCreated: time.Now(),
		Name:        "Cheese & Wine",
		TrackerID:   domainTracker.ID,
		UserID:      userOne.ID,
		Value:       decimal.NewFromFloat(10.99),
	}

	givenIHaveASpend(spend)
	givenISaveTheSpend(t)
	whenIGetTheSpend(t)
	thenTheFollowingSpendIsReturned(spend, t)
}

func TestCanGetForTrackerID(t *testing.T) {

	givenIHaveCleanDependencies(t)
	givenThereIsATrackerAndTwoUsers(t)

	spend := model.Spend{
		Currency:    domainTracker.Currency,
		DateCreated: time.Now(),
		Name:        "Cheese & Wine",
		TrackerID:   domainTracker.ID,
		UserID:      userOne.ID,
		Value:       decimal.NewFromFloat(10.99),
	}

	givenIHaveASpend(spend)
	givenISaveTheSpend(t)
	whenIGetTheSpendsForATracker(t)

	expected := []model.Spend{spend}

	thenTheFollowingSpendsAreReturned(expected, t)
}

func TestCanInsertSpend(t *testing.T) {
	givenIHaveCleanDependencies(t)
	givenThereIsATrackerAndTwoUsers(t)

	spend := model.Spend{
		Currency:    domainTracker.Currency,
		DateCreated: time.Now(),
		Name:        "Cheese",
		TrackerID:   domainTracker.ID,
		UserID:      userOne.ID,
		Value:       decimal.NewFromFloat(1.99),
	}

	givenIHaveASpend(spend)
	whenISaveTheSpend(t)
	thenTheSpendIsSaved(t)
}

func whenIUpdateTheSpend(spend model.Spend, t *testing.T) {
	domainSpend, err = spendRepository.Update(spend.ID, spend)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}

func thenTheSpendIsUpdated(expected model.Spend, t *testing.T) {
	domainSpend, err = spendRepository.GetByID(domainSpend.ID)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if domainSpend.Currency != expected.Currency {
		t.Fatalf("Expected: %v, Received: %v", expected.Currency, domainSpend.Currency)
	}
	if domainSpend.Name != expected.Name {
		t.Fatalf("Expected: %v, Received: %v", expected.Name, domainSpend.Name)
	}
	if domainSpend.TrackerID != expected.TrackerID {
		t.Fatalf("Expected: %v, Received: %v", expected.TrackerID, domainSpend.TrackerID)
	}
	if domainSpend.UserID != expected.UserID {
		t.Fatalf("Expected: %v, Received: %v", expected.UserID, domainSpend.UserID)
	}
	if domainSpend.Value.String() != expected.Value.String() {
		t.Fatalf("Expected: %v, Received: %v", expected.Value, domainSpend.Value)
	}
}

func whenIDeleteTheSpend(t *testing.T) {
	result, err = spendRepository.Delete(domainSpend.ID)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}

func thenTheSpendIsDeleted(t *testing.T) {
	if result == false {
		t.Fatalf("The spend was not deleted")
	}
}

func thenTheFollowingSpendsAreReturned(expected []model.Spend, t *testing.T) {
	for i := 0; i < len(domainSpends); i++ {
		if domainSpends[i].Currency != expected[i].Currency {
			t.Fatalf("Expected: %v, Received: %v", expected[i].Currency, domainSpends[i].Currency)
		}
		if domainSpends[i].Name != expected[i].Name {
			t.Fatalf("Expected: %v, Received: %v", expected[i].Name, domainSpends[i].Name)
		}
		if domainSpends[i].TrackerID != expected[i].TrackerID {
			t.Fatalf("Expected: %v, Received: %v", expected[i].TrackerID, domainSpends[i].TrackerID)
		}
		if domainSpends[i].UserID != expected[i].UserID {
			t.Fatalf("Expected: %v, Received: %v", expected[i].UserID, domainSpends[i].UserID)
		}
		if domainSpends[i].Value.String() != expected[i].Value.String() {
			t.Fatalf("Expected: %v, Received: %v", expected[i].Value, domainSpends[i].Value)
		}
	}
}

func thenTheFollowingSpendIsReturned(expected model.Spend, t *testing.T) {
	if domainSpend.Currency != expected.Currency {
		t.Fatalf("Expected: %v, Received: %v", expected.Currency, domainSpend.Currency)
	}
	if domainSpend.Name != expected.Name {
		t.Fatalf("Expected: %v, Received: %v", expected.Name, domainSpend.Name)
	}
	if domainSpend.TrackerID != expected.TrackerID {
		t.Fatalf("Expected: %v, Received: %v", expected.TrackerID, domainSpend.TrackerID)
	}
	if domainSpend.UserID != expected.UserID {
		t.Fatalf("Expected: %v, Received: %v", expected.UserID, domainSpend.UserID)
	}
	if domainSpend.Value.String() != expected.Value.String() {
		t.Fatalf("Expected: %v, Received: %v", expected.Value, domainSpend.Value)
	}
}

func whenIGetTheSpendsForATracker(t *testing.T) {
	domainSpends, err = spendRepository.GetForTrackerID(domainTracker.ID)
	if err != nil {
		t.Fatal("Error: ", err)
	}
}

func whenIGetTheSpend(t *testing.T) {
	domainSpend, err = spendRepository.GetByID(domainSpend.ID)
	if err != nil {
		t.Fatal("Error: ", err)
	}
}

func whenISaveTheSpend(t *testing.T) {
	domainSpend, err = spendRepository.Insert(domainSpend)
	if err != nil {
		t.Fatal("Error: ", err)
	}
}

func givenISaveTheSpend(t *testing.T) {
	domainSpend, err = spendRepository.Insert(domainSpend)
	if err != nil {
		t.Fatal("Error: ", err)
	}
}

func thenTheSpendIsSaved(t *testing.T) {
	if domainSpend.ID == 0 {
		t.Fatal("The spend id was zero")
	}
}

func givenIHaveASpend(spend model.Spend) {
	domainSpend = spend
}

func givenThereIsATrackerAndTwoUsers(t *testing.T) {
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

	userOne = domainUser

	givenThereIsAUser(t)

	userTwo = domainUser

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

	givenIHaveADomainUser(user)
	whenISaveTheUser(t)
	thenTheUserIsSaved(t)
}

func givenIHaveCleanDependencies(t *testing.T) {
	db, err = repository.NewDB("postgres://godutch:password@localhost/godutch?sslmode=disable")
	userRepository = userrepository.NewPostgresUserRepository(infrastructure.ConsoleLogger{}, db)
	trackerRepository = trackerrepository.NewPostgresTrackerRepository(infrastructure.ConsoleLogger{}, db)
	spendRepository = spendrepository.NewPostgresSpendRepository(infrastructure.ConsoleLogger{}, db)
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
