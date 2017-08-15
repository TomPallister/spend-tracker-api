package trackervalidation

import (
	"errors"

	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
)

// ErrorInvalidName ...
var ErrorInvalidName = errors.New("Invalid name")

// ErrorInvalidAdminUserID ...
var ErrorInvalidAdminUserID = errors.New("Invalid admin user id")

// ErrorAdminUserIDIsDifferentToSubjectID ...
var ErrorAdminUserIDIsDifferentToSubjectID = errors.New("Admin user id is different to subject id")

// ErrorNoTrackerUsers ...
var ErrorNoTrackerUsers = errors.New("No tracker users")

// ErrorAdminUserNotInTrackerUsersList ...
var ErrorAdminUserNotInTrackerUsersList = errors.New("Admin user id not in tracker user ids")

// ErrorInvalidDateCreated ...
var ErrorInvalidDateCreated = errors.New("Invalid date created")

// ErrorTheTrackerDoesNotExist ...
var ErrorTheTrackerDoesNotExist = errors.New("Tracker does not exist")
 
// TrackerValidator ...
type TrackerValidator interface {
	IsValidCreateTracker(tracker model.Tracker, logger infrastructure.Logger, adminUser model.User) (bool, error)
	IsValidUpdateTracker(tracker model.Tracker, logger infrastructure.Logger, adminUser model.User, existingTracker model.Tracker) (bool, error)
	IsValidDeleteTracker(id int64, logger infrastructure.Logger, adminUser model.User, existingTracker model.Tracker) (bool, error)
}

// GoDutchTrackerValidator ...
type GoDutchTrackerValidator struct {
}

// NewGoDutchTrackerValidator ...
func NewGoDutchTrackerValidator() *GoDutchTrackerValidator {

	service := GoDutchTrackerValidator{}

	return &service
}

// IsValidCreateTracker ...
func (validator *GoDutchTrackerValidator) IsValidCreateTracker(tracker model.Tracker, logger infrastructure.Logger, adminUser model.User) (bool, error) {

	if len(tracker.Name) <= 0 {
		logger.Error("Error: ", ErrorInvalidName)
		return false, ErrorInvalidName
	}
	if tracker.AdminUserID <= 0 {
		logger.Error("Error: ", ErrorInvalidAdminUserID)
		return false, ErrorInvalidAdminUserID
	}

	if len(tracker.TrackerUserIDs) <= 0 {
		logger.Error("Error: ", ErrorNoTrackerUsers)
		return false, ErrorNoTrackerUsers
	}

	if !infrastructure.Ints64Contains(tracker.TrackerUserIDs, tracker.AdminUserID) {
		logger.Error("Error: ", ErrorAdminUserNotInTrackerUsersList)
		return false, ErrorAdminUserNotInTrackerUsersList
	}

	if tracker.DateCreated.IsZero() {
		logger.Error("Error: ", ErrorInvalidDateCreated)
		return false, ErrorInvalidDateCreated
	}

	if tracker.AdminUserID != adminUser.ID {
		logger.Error("Error: ", ErrorAdminUserIDIsDifferentToSubjectID)
		return false, ErrorAdminUserIDIsDifferentToSubjectID
	}

	return true, nil
}

// IsValidUpdateTracker ...
func (validator *GoDutchTrackerValidator) IsValidUpdateTracker(tracker model.Tracker, logger infrastructure.Logger, adminUser model.User, existingTracker model.Tracker) (bool, error) {

	if len(tracker.Name) <= 0 {
		logger.Error("Error: ", ErrorInvalidName)
		return false, ErrorInvalidName
	}
	if tracker.AdminUserID <= 0 {
		logger.Error("Error: ", ErrorInvalidAdminUserID)
		return false, ErrorInvalidAdminUserID
	}

	if len(tracker.TrackerUserIDs) <= 0 {
		logger.Error("Error: ", ErrorNoTrackerUsers)
		return false, ErrorNoTrackerUsers
	}

	if !infrastructure.Ints64Contains(tracker.TrackerUserIDs, tracker.AdminUserID) {
		logger.Error("Error: ", ErrorAdminUserNotInTrackerUsersList)
		return false, ErrorAdminUserNotInTrackerUsersList
	}

	if tracker.DateCreated.IsZero() {
		logger.Error("Error: ", ErrorInvalidDateCreated)
		return false, ErrorInvalidDateCreated
	}

	if tracker.AdminUserID != adminUser.ID {
		logger.Error("Error: ", ErrorAdminUserIDIsDifferentToSubjectID)
		return false, ErrorAdminUserIDIsDifferentToSubjectID
	}

	if tracker.ID != existingTracker.ID {
		logger.Error("Error: ", ErrorTheTrackerDoesNotExist)
		return false, ErrorTheTrackerDoesNotExist
	}

	return true, nil
}

// IsValidDeleteTracker ...
func (validator *GoDutchTrackerValidator) IsValidDeleteTracker(id int64, logger infrastructure.Logger, adminUser model.User, existingTracker model.Tracker) (bool, error) {

	if existingTracker.AdminUserID != adminUser.ID {
		logger.Error("Error: ", ErrorAdminUserIDIsDifferentToSubjectID)
		return false, ErrorAdminUserIDIsDifferentToSubjectID
	}

	if id != existingTracker.ID {
		logger.Error("Error: ", ErrorTheTrackerDoesNotExist)
		return false, ErrorTheTrackerDoesNotExist
	}

	return true, nil
}
