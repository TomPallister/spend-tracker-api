package uservalidation_test

import (
	"testing"
	"time"

	"github.com/TomPallister/godutch-api/api/domain/userservice"
	"github.com/TomPallister/godutch-api/api/domain/validation/uservalidation"
	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/repository/userrepository"
)

var result bool
var logger = infrastructure.NilLogger{}
var newUser model.User
var err error
var userRepository *userrepository.InMemoryUserRepository
var userValidator = uservalidation.NewGoDutchUserValidator()

func TestValidateInviteUserNoEmail(t *testing.T) {

	user := model.User{
		Name:             "Tom",
		AuthenticationID: "ID",
		DateCreated:      time.Now()}

	givenIHaveACreateUserCommand(user)
	whenIValidateTheInviteUserCommand()
	thenTheCreateUserCommandIsRejectedWithError(uservalidation.ErrorInvalidEmail, t)
}

func TestValidateInviteUserInvalidEmail(t *testing.T) {

	user := model.User{
		Name:             "Tom",
		AuthenticationID: "ID",
		EmailAddress:     "email",
		DateCreated:      time.Now()}

	givenIHaveACreateUserCommand(user)
	whenIValidateTheInviteUserCommand()
	thenTheCreateUserCommandIsRejectedWithError(uservalidation.ErrorInvalidEmail, t)
}

func TestValidateInviteUserNoDateCreated(t *testing.T) {

	user := model.User{
		Name:             "",
		AuthenticationID: "",
		EmailAddress:     "email@"}

	givenIHaveACreateUserCommand(user)
	whenIValidateTheInviteUserCommand()
	thenTheCreateUserCommandIsRejectedWithError(uservalidation.ErrorInvalidDateCreated, t)
}

func TestValidateInviteUserNoUserWithEmailAlreadyExists(t *testing.T) {

	user := model.User{
		Name:             "",
		AuthenticationID: "",
		EmailAddress:     "email@",
		DateCreated:      time.Now()}

	givenThereIsAUserInTheRepository(model.User{ID: 0, Name: "Tom", AuthenticationID: "jeff", EmailAddress: "email@"}, t)
	givenIHaveACreateUserCommand(user)
	whenIValidateTheInviteUserCommand()
	thenTheCreateUserCommandIsRejectedWithError(uservalidation.ErrorEmailAlreadyInUse, t)
}

func TestValidateInviteUser(t *testing.T) {

	user := model.User{
		Name:             "",
		AuthenticationID: "",
		EmailAddress:     "email123@",
		DateCreated:      time.Now()}

	givenIHaveACreateUserCommand(user)
	whenIValidateTheInviteUserCommand()
	thenTheCreateUserCommandIsAccepted(t)
}

func TestValidateCreateUserNoName(t *testing.T) {

	user := model.User{
		AuthenticationID: "ID",
		EmailAddress:     "emailAddress",
		DateCreated:      time.Now()}

	givenIHaveACreateUserCommand(user)
	whenIValidateTheCreateUserCommand()
	thenTheCreateUserCommandIsRejectedWithError(uservalidation.ErrorInvalidName, t)
}

func TestValidateCreateUserNoEmail(t *testing.T) {

	user := model.User{
		Name:             "Tom",
		AuthenticationID: "ID",
		DateCreated:      time.Now()}

	givenIHaveACreateUserCommand(user)
	whenIValidateTheCreateUserCommand()
	thenTheCreateUserCommandIsRejectedWithError(uservalidation.ErrorInvalidEmail, t)
}

func TestValidateCreateUserInvalidEmail(t *testing.T) {

	user := model.User{
		Name:             "Tom",
		AuthenticationID: "ID",
		EmailAddress:     "email",
		DateCreated:      time.Now()}

	givenIHaveACreateUserCommand(user)
	whenIValidateTheCreateUserCommand()
	thenTheCreateUserCommandIsRejectedWithError(uservalidation.ErrorInvalidEmail, t)
}

func TestValidateCreateUserNoAuthenticationID(t *testing.T) {

	user := model.User{
		Name:         "Tom",
		EmailAddress: "email@",
		DateCreated:  time.Now()}

	givenIHaveACreateUserCommand(user)
	whenIValidateTheCreateUserCommand()
	thenTheCreateUserCommandIsRejectedWithError(uservalidation.ErrorInvalidAuthenticationID, t)
}

func TestValidateCreateUserNoDateCreated(t *testing.T) {

	user := model.User{
		Name:             "Tom",
		AuthenticationID: "authid",
		EmailAddress:     "email@"}

	givenIHaveACreateUserCommand(user)
	whenIValidateTheCreateUserCommand()
	thenTheCreateUserCommandIsRejectedWithError(uservalidation.ErrorInvalidDateCreated, t)
}

func TestValidateCreateUserNoUserWithEmailAlreadyExists(t *testing.T) {

	user := model.User{
		Name:             "Tom",
		AuthenticationID: "authid",
		EmailAddress:     "email@",
		DateCreated:      time.Now()}

	givenThereIsAUserInTheRepository(model.User{ID: 0, Name: "Tom", AuthenticationID: "jeff", EmailAddress: "email@"}, t)
	givenIHaveACreateUserCommand(user)
	whenIValidateTheCreateUserCommand()
	thenTheCreateUserCommandIsRejectedWithError(uservalidation.ErrorEmailAlreadyInUse, t)
}

func TestValidateCreateUserNoUserWithAuthIdAlreadyExists(t *testing.T) {

	user := model.User{
		Name:             "Tom",
		AuthenticationID: "authid",
		EmailAddress:     "email@",
		DateCreated:      time.Now()}

	givenThereIsAUserInTheRepository(model.User{ID: 0, Name: "Tom", AuthenticationID: "authid"}, t)
	givenIHaveACreateUserCommand(user)
	whenIValidateTheCreateUserCommand()
	thenTheCreateUserCommandIsRejectedWithError(uservalidation.ErrorAuthIDAlreadyInUse, t)
}

func TestValidateCreateUser(t *testing.T) {

	user := model.User{
		Name:             "Tom",
		AuthenticationID: "authid123",
		EmailAddress:     "email123@",
		DateCreated:      time.Now()}

	givenIHaveACreateUserCommand(user)
	whenIValidateTheCreateUserCommand()
	thenTheCreateUserCommandIsAccepted(t)
}

func givenThereIsAUserInTheRepository(newUser model.User, t *testing.T) {
	userRepository = userrepository.NewInMemoryUserRepository()
	_, err := userRepository.Insert(newUser)
	if err != nil {
		t.Fatalf("There was an error %v", err)
	}
}
func whenIValidateTheCreateUserCommand() {
	result, err = userValidator.IsValidCreateUser(newUser, logger, userRepository)
}

func whenIValidateTheInviteUserCommand() {
	result, err = userValidator.IsValidInviteUser(newUser, logger, userRepository)
}

func thenTheCreateUserCommandIsRejectedWithError(e error, t *testing.T) {
	if err != e {
		t.Fatalf("Error should be %v but was %v", e, err)
	}
}

func thenTheCreateUserCommandIsAccepted(t *testing.T) {
	if err != nil {
		t.Fatalf("There was an error")
	}
	if result == false {
		t.Fatalf("There result was false and it should be true")
	}
}

func givenIHaveACreateUserCommand(user model.User) {
	newUser = user
}

func thenTheUserIsCreated(t *testing.T) {
	if err == userservice.ErrorCreateUser {
		t.Fatalf("Could not create the user")
	}
}
