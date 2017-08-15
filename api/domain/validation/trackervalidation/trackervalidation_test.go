package trackervalidation_test

import (
	"testing"
	"time"

	"github.com/TomPallister/godutch-api/api/domain/validation/trackervalidation"
	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
)

var result bool
var logger = infrastructure.NilLogger{}
var newTracker model.Tracker
var newUser model.User
var trackerValidator = trackervalidation.NewGoDutchTrackerValidator()
var err error

func TestValidateCreateTrackerNoName(t *testing.T) {

	tracker := model.Tracker{
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: make([]int, 1),
	}
	user := model.User{
		ID: 1,
	}

	givenThereIsAUser(user)
	givenIHaveATracker(tracker)
	whenIValidateTheTracker()
	thenTheCreateTrackerCommandIsRejectedWithError(trackervalidation.ErrorInvalidName, t)
}

func TestValidateCreateTrackerNoAdminUserID(t *testing.T) {

	tracker := model.Tracker{
		Name:           "Tom and Laura",
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
	}
	user := model.User{
		ID: 1,
	}

	givenThereIsAUser(user)
	givenIHaveATracker(tracker)
	whenIValidateTheTracker()
	thenTheCreateTrackerCommandIsRejectedWithError(trackervalidation.ErrorInvalidAdminUserID, t)

}

func TestValidateCreateTrackerNoTrackerUserIDs(t *testing.T) {

	var trackerUsers = []int{}

	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: trackerUsers,
	}
	user := model.User{
		ID: 1,
	}

	givenThereIsAUser(user)
	givenIHaveATracker(tracker)
	whenIValidateTheTracker()
	thenTheCreateTrackerCommandIsRejectedWithError(trackervalidation.ErrorNoTrackerUsers, t)

}

func TestValidateCreateTrackerAdminUserIDIsNotInTrackerUserIDs(t *testing.T) {

	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    12,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
	}
	user := model.User{
		ID: 1,
	}

	givenThereIsAUser(user)
	givenIHaveATracker(tracker)
	whenIValidateTheTracker()
	thenTheCreateTrackerCommandIsRejectedWithError(trackervalidation.ErrorAdminUserNotInTrackerUsersList, t)

}

func TestValidateTrackerNoDateCreated(t *testing.T) {

	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		TrackerUserIDs: []int{1},
	}
	user := model.User{
		ID: 1,
	}

	givenThereIsAUser(user)
	givenIHaveATracker(tracker)
	whenIValidateTheTracker()
	thenTheCreateTrackerCommandIsRejectedWithError(trackervalidation.ErrorInvalidDateCreated, t)
}

func TestValidateTrackerAdminUserIDIsDifferentToUserID(t *testing.T) {

	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		TrackerUserIDs: []int{1},
		DateCreated:    time.Now(),
	}
	user := model.User{
		ID: 12,
	}

	givenThereIsAUser(user)
	givenIHaveATracker(tracker)
	whenIValidateTheTracker()
	thenTheCreateTrackerCommandIsRejectedWithError(trackervalidation.ErrorAdminUserIDIsDifferentToSubjectID, t)
}

func TestValidateTracker(t *testing.T) {
	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
	}

	user := model.User{
		ID: 1,
	}

	givenThereIsAUser(user)
	givenIHaveATracker(tracker)
	whenIValidateTheTracker()
	thenTheCreateTrackerCommandIsAccepted(t)
}

func TestValidateUpdateTracker(t *testing.T) {
	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
	}

	user := model.User{
		ID: 1,
	}

	givenThereIsAUser(user)
	givenIHaveATracker(tracker)
	whenIValidateTheUpdateTracker()
	thenTheCreateTrackerCommandIsAccepted(t)
}

func TestValidateDeleteTracker(t *testing.T) {
	tracker := model.Tracker{
		Name:           "Tom and Laura",
		AdminUserID:    1,
		DateCreated:    time.Now(),
		TrackerUserIDs: []int{1},
	}

	user := model.User{
		ID: 1,
	}

	givenThereIsAUser(user)
	givenIHaveATracker(tracker)
	whenIValidateTheDeleteTracker()
	thenTheCreateTrackerCommandIsAccepted(t)
}

func givenThereIsAUser(user model.User) {
	newUser = user
}

func givenIHaveATracker(tracker model.Tracker) {
	newTracker = tracker
}

func whenIValidateTheTracker() {
	result, err = trackerValidator.IsValidCreateTracker(newTracker, logger, newUser)
}

func whenIValidateTheUpdateTracker() {
	result, err = trackerValidator.IsValidUpdateTracker(newTracker, logger, newUser, newTracker)
}

func whenIValidateTheDeleteTracker() {
	result, err = trackerValidator.IsValidDeleteTracker(newTracker.ID, logger, newUser, newTracker)
}

func thenTheCreateTrackerCommandIsRejectedWithError(e error, t *testing.T) {
	if err != e {
		t.Fatalf("Error should be %v but was %v", e, err)
	}
}

func thenTheCreateTrackerCommandIsAccepted(t *testing.T) {
	if err != nil {
		t.Fatalf("There was an error")
	}
	if result == false {
		t.Fatalf("There result was false and it should be true")
	}
}
