package userservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/TomPallister/godutch-api/api/domain/validation/uservalidation"
	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/infrastructure/encryption"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/repository/trackerrepository"
	"github.com/TomPallister/godutch-api/api/repository/userrepository"
	"github.com/nu7hatch/gouuid"
)

// ErrorCreateUser ...
var ErrorCreateUser = errors.New("Could not create user")

// ErrorAcceptInviteUser ...
var ErrorAcceptInviteUser = errors.New("Could not accept invite user")

// ErrorPermissionsToViewUser ...
var ErrorPermissionsToViewUser = errors.New("You do not have permission to view this user")

// ErrorDoNotHavePermissionToInvite ...
var ErrorDoNotHavePermissionToInvite = errors.New("You do not have permission to invite to thsi tracker")

// UserService ...
type UserService interface {
	FindBySub(sub string) (model.User, error)

	FindByID(sub string, id int64) (model.User, error)

	CreateUser(sub string,
		user model.User) (model.User, error)

	InviteUser(sub string, inviteUser model.InviteUser, rootURL string) (model.User, error)

	AcceptInvite(sub string, user model.User, cryptoText string) (model.User, error)

	FindByIDForTracker(id int64) (model.User, error)
}

// NewGoDutchUserService ...
func NewGoDutchUserService(userRepository userrepository.UserRepository,
	validator uservalidation.UserValidator,
	logger infrastructure.Logger,
	emailService infrastructure.EmailService,
	trackerRepository trackerrepository.TrackerRepository) *GoDutchUserService {

	service := GoDutchUserService{}
	service.userRepository = userRepository
	service.validator = validator
	service.logger = logger
	service.emailService = emailService
	service.trackerRepository = trackerRepository
	return &service
}

// GoDutchUserService ...
type GoDutchUserService struct {
	userRepository    userrepository.UserRepository
	validator         uservalidation.UserValidator
	logger            infrastructure.Logger
	emailService      infrastructure.EmailService
	trackerRepository trackerrepository.TrackerRepository
}

// FindBySub ...
func (godutchUserService *GoDutchUserService) FindBySub(sub string) (model.User, error) {
	return godutchUserService.userRepository.GetBySub(sub)
}

// FindByID ...
func (godutchUserService *GoDutchUserService) FindByID(sub string,
	id int64) (model.User, error) {

	currentUser, err := godutchUserService.userRepository.GetBySub(sub)
	if err != nil {
		return model.User{}, err
	}
	cuerrenUsersTrackers, err := godutchUserService.trackerRepository.GetForUserID(currentUser.ID)
	if err != nil {
		return model.User{}, err
	}

	userIdsTheCurrentUserCanSee := []int64{}

	for _, t := range cuerrenUsersTrackers {
		for _, u := range t.TrackerUserIDs {
			userIdsTheCurrentUserCanSee = append(userIdsTheCurrentUserCanSee, u)
		}
	}

	if !infrastructure.Ints64Contains(userIdsTheCurrentUserCanSee, id) {
		return model.User{}, ErrorPermissionsToViewUser
	}

	user, err := godutchUserService.userRepository.GetByID(id)
	return user, err
}

// FindByIDForTracker ...
func (godutchUserService *GoDutchUserService) FindByIDForTracker(id int64) (model.User, error) {

	user, err := godutchUserService.userRepository.GetByID(id)
	return user, err
}

// CreateUser ...
func (godutchUserService *GoDutchUserService) CreateUser(sub string,
	user model.User) (model.User, error) {

	user.DateCreated = time.Now()

	user.AuthenticationID = sub

	valid, err := godutchUserService.validator.IsValidCreateUser(user, godutchUserService.logger, godutchUserService.userRepository)
	if valid == false {
		return model.User{}, err
	}

	return godutchUserService.userRepository.Insert(user)
}

// InviteUser ...
func (godutchUserService *GoDutchUserService) InviteUser(sub string, inviteUser model.InviteUser,
	rootURL string) (model.User, error) {

	invitingUser, err := godutchUserService.FindBySub(sub)
	if err != nil {
		return model.User{}, err
	}

	tracker, err := godutchUserService.trackerRepository.GetByID(inviteUser.TrackerID)
	if err != nil {
		return model.User{}, err
	}

	if invitingUser.ID != tracker.AdminUserID {
		return model.User{}, ErrorDoNotHavePermissionToInvite
	}
	tempAuthID, _ := uuid.NewV4()

	invitedUser := model.User{
		AuthenticationID: tempAuthID.String(),
		DateCreated:      time.Now(),
		EmailAddress:     inviteUser.EmailAddress,
		Name:             "",
	}

	valid, err := godutchUserService.validator.IsValidInviteUser(invitedUser, godutchUserService.logger, godutchUserService.userRepository)
	if valid == false {
		return model.User{}, err
	}

	invitedUser, err = godutchUserService.userRepository.Insert(invitedUser)
	if err != nil {
		return model.User{}, err
	}

	tracker.TrackerUserIDs = append(tracker.TrackerUserIDs, invitedUser.ID)

	tracker, err = godutchUserService.trackerRepository.Update(tracker.ID, tracker)
	if err != nil {
		return model.User{}, err
	}

	ciphertext := encryption.Encrypt(invitedUser.EmailAddress)

	result, err := godutchUserService.emailService.SendEmail(invitedUser.EmailAddress,
		fmt.Sprintf("You have been invited to GoDutch by %v", invitingUser.EmailAddress),
		fmt.Sprintf("go to %vacceptinvite/%v to sign up..", rootURL, ciphertext),
		"computer@godutch.money", "GoDutch")

	if result == false {
		return model.User{}, err
	}

	return invitedUser, nil

}

// AcceptInvite ...
func (godutchUserService *GoDutchUserService) AcceptInvite(sub string, user model.User,
	cryptoEmail string) (model.User, error) {

	emailAddress := encryption.Decrypt(cryptoEmail)

	invitedUser, err := godutchUserService.userRepository.GetByEmail(emailAddress)
	if err != nil {
		return model.User{}, err
	}

	if emailAddress != invitedUser.EmailAddress {
		return model.User{}, ErrorAcceptInviteUser
	}

	invitedUser.AuthenticationID = sub
	invitedUser.Name = user.Name

	if len(user.EmailAddress) > 0 {
		invitedUser.EmailAddress = user.EmailAddress
	}

	invitedUser, err = godutchUserService.userRepository.Update(invitedUser.ID, invitedUser)
	if err != nil {
		return model.User{}, err
	}

	return invitedUser, nil
}
