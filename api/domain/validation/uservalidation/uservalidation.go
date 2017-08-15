package uservalidation

import (
	"errors"
	"strings"

	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/repository/userrepository"
)

// ErrorInvalidName ...
var ErrorInvalidName = errors.New("Invalid name")

// ErrorInvalidEmail ...
var ErrorInvalidEmail = errors.New("Invalid email")

// ErrorInvalidAuthenticationID ...
var ErrorInvalidAuthenticationID = errors.New("Invalid authentication ID")

// ErrorInvalidDateCreated ...
var ErrorInvalidDateCreated = errors.New("Invalid date created")

// ErrorEmailAlreadyInUse ...
var ErrorEmailAlreadyInUse = errors.New("Email already in use.")

// ErrorAuthIDAlreadyInUse ...
var ErrorAuthIDAlreadyInUse = errors.New("AuthID already in use.")

// UserValidator ...
type UserValidator interface {
	IsValidCreateUser(user model.User, logger infrastructure.Logger, userRepo userrepository.UserRepository) (bool, error)
	IsValidInviteUser(user model.User, logger infrastructure.Logger, userRepo userrepository.UserRepository) (bool, error)
}

// GoDutchUserValidator ...
type GoDutchUserValidator struct {
}

// NewGoDutchUserValidator ...
func NewGoDutchUserValidator() *GoDutchUserValidator {

	service := GoDutchUserValidator{}

	return &service
}

// IsValidInviteUser ...
func (validator *GoDutchUserValidator) IsValidInviteUser(user model.User, logger infrastructure.Logger, userRepo userrepository.UserRepository) (bool, error) {

	if len(user.EmailAddress) <= 0 {
		logger.Error("Error: ", ErrorInvalidEmail)
		return false, ErrorInvalidEmail
	}

	if !strings.Contains(user.EmailAddress, "@") {
		logger.Error("Error: ", ErrorInvalidEmail)
		return false, ErrorInvalidEmail
	}

	if user.DateCreated.IsZero() {
		logger.Error("Error: ", ErrorInvalidDateCreated)
		return false, ErrorInvalidDateCreated
	}

	userAlreadyExisits, _ := userRepo.GetByEmail(user.EmailAddress)
	if userAlreadyExisits.ID > 0 {
		logger.Error("Error: ", ErrorEmailAlreadyInUse)
		return false, ErrorEmailAlreadyInUse
	}

	return true, nil
}

// IsValidCreateUser ...
func (validator *GoDutchUserValidator) IsValidCreateUser(user model.User, logger infrastructure.Logger, userRepo userrepository.UserRepository) (bool, error) {

	if len(user.Name) <= 0 {
		logger.Error("Error: ", ErrorInvalidName)
		return false, ErrorInvalidName
	}

	if len(user.EmailAddress) <= 0 {
		logger.Error("Error: ", ErrorInvalidEmail)
		return false, ErrorInvalidEmail
	}

	if !strings.Contains(user.EmailAddress, "@") {
		logger.Error("Error: ", ErrorInvalidEmail)
		return false, ErrorInvalidEmail
	}

	if len(user.AuthenticationID) <= 0 {
		logger.Error("Error: ", ErrorInvalidAuthenticationID)
		return false, ErrorInvalidAuthenticationID
	}

	if user.DateCreated.IsZero() {
		logger.Error("Error: ", ErrorInvalidDateCreated)
		return false, ErrorInvalidDateCreated
	}

	userAlreadyExisits, _ := userRepo.GetByEmail(user.EmailAddress)
	if userAlreadyExisits.ID > 0 {
		logger.Error("Error: ", ErrorEmailAlreadyInUse)
		return false, ErrorEmailAlreadyInUse
	}

	userAlreadyExisits, _ = userRepo.GetBySub(user.AuthenticationID)
	if userAlreadyExisits.ID > 0 {
		logger.Error("Error: ", ErrorAuthIDAlreadyInUse)
		return false, ErrorAuthIDAlreadyInUse
	}

	return true, nil
}
