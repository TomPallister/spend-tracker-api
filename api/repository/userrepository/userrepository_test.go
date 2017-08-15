package userrepository_test

import (
	"testing"
	"time"

	"database/sql"

	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/repository"
	"github.com/TomPallister/godutch-api/api/repository/userrepository"
	"github.com/nu7hatch/gouuid"
)

var domainUser model.User
var err error
var userRepository userrepository.UserRepository
var domainUserID int64
var db *sql.DB

func TestCanInsertUser(t *testing.T) {

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

func TestCanGetUserBySub(t *testing.T) {

	u, _ := uuid.NewV4()

	user := model.User{
		Name:             "Laura",
		AuthenticationID: "21312fsdf" + u.String(),
		DateCreated:      time.Now(),
		EmailAddress:     "laura" + u.String() + "@laz.com",
	}

	givenIHaveCleanDependencies(t)
	givenIHaveADomainUser(user)
	givenISaveTheUser(t)
	whenIGetTheUserBySub(t)
	thenTheUserIsReturned(t)
}

func TestCanGetUserByID(t *testing.T) {

	u, _ := uuid.NewV4()

	user := model.User{
		Name:             "Laura",
		AuthenticationID: "21312fsdf" + u.String(),
		DateCreated:      time.Now(),
		EmailAddress:     "laura" + u.String() + "@laz.com",
	}

	givenIHaveCleanDependencies(t)
	givenIHaveADomainUser(user)
	givenISaveTheUser(t)
	whenIGetTheUserByID(t)
	thenTheUserIsReturned(t)
}

func TestCanGetUserByEmailAddress(t *testing.T) {

	u, _ := uuid.NewV4()

	user := model.User{
		Name:             "Laura",
		AuthenticationID: "21312fsdf" + u.String(),
		DateCreated:      time.Now(),
		EmailAddress:     "laura" + u.String() + "@laz.com",
	}

	givenIHaveCleanDependencies(t)
	givenIHaveADomainUser(user)
	givenISaveTheUser(t)
	whenIGetTheUserByEmail(t)
	thenTheUserIsReturned(t)
}

func TestCanUpdateUser(t *testing.T) {

	u, _ := uuid.NewV4()

	user := model.User{
		Name:             "Laura",
		AuthenticationID: "21312fsdf" + u.String(),
		DateCreated:      time.Now(),
		EmailAddress:     "laura" + u.String() + "@laz.com",
	}

	givenIHaveCleanDependencies(t)
	givenIHaveADomainUser(user)
	givenISaveTheUser(t)
	givenIGetTheUserByID(t)

	u, _ = uuid.NewV4()

	updateUser := model.User{
		ID:               domainUser.ID,
		AuthenticationID: "newSub" + u.String(),
		EmailAddress:     "newEmail" + u.String(),
		DateCreated:      domainUser.DateCreated,
		Name:             "newName",
	}

	whenIUpdateTheUser(updateUser, t)
	thenTheUserIsUpdated(updateUser, t)
}

func thenTheUserIsUpdated(expected model.User, t *testing.T) {
	domainUser, err = userRepository.GetByID(domainUser.ID)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}

	if domainUser.AuthenticationID != expected.AuthenticationID {
		t.Fatalf("Expected: %v, Received: %v", expected.AuthenticationID, domainUser.AuthenticationID)
	}

	if domainUser.EmailAddress != expected.EmailAddress {
		t.Fatalf("Expected: %v, Received: %v", expected.EmailAddress, domainUser.EmailAddress)
	}

	if domainUser.Name != expected.Name {
		t.Fatalf("Expected: %v, Received: %v", expected.Name, domainUser.Name)
	}
}

func whenIUpdateTheUser(user model.User, t *testing.T) {
	domainUser, err = userRepository.Update(user.ID, user)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func whenIGetTheUserBySub(t *testing.T) {
	domainUser, err = userRepository.GetBySub(domainUser.AuthenticationID)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func whenIGetTheUserByID(t *testing.T) {
	domainUser, err = userRepository.GetByID(domainUser.ID)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}
func givenIGetTheUserByID(t *testing.T) {
	domainUser, err = userRepository.GetByID(domainUser.ID)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func whenIGetTheUserByEmail(t *testing.T) {
	domainUser, err = userRepository.GetByEmail(domainUser.EmailAddress)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}

func thenTheUserIsReturned(t *testing.T) {
	if domainUser.ID == 0 {
		t.Fatalf("Expected: %v, Received: %v", domainUserID, domainUser.ID)
	}
}

func givenIHaveCleanDependencies(t *testing.T) {
	db, err = repository.NewDB("postgres://godutch:password@localhost/godutch?sslmode=disable")
	userRepository = userrepository.NewPostgresUserRepository(infrastructure.ConsoleLogger{}, db)
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

func givenISaveTheUser(t *testing.T) {
	domainUser, err = userRepository.Insert(domainUser)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
	domainUserID = domainUser.ID
}

func thenTheUserIsSaved(t *testing.T) {
	if domainUser.ID <= 0 {
		t.Fatalf("The user ID was %v", domainUser.ID)
	}
}
