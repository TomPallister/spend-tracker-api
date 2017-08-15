package trackerservice

import (
	"errors"
	"time"

	"github.com/TomPallister/godutch-api/api/domain/spendsummaryservice"
	"github.com/TomPallister/godutch-api/api/domain/transferservice"
	"github.com/TomPallister/godutch-api/api/domain/userservice"
	"github.com/TomPallister/godutch-api/api/domain/validation/trackervalidation"
	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/repository/spendrepository"
	"github.com/TomPallister/godutch-api/api/repository/trackerrepository"
)

// ErrorCreateTracker ...
var ErrorCreateTracker = errors.New("Could not create tracker")

// ErrorUpdateTracker ...
var ErrorUpdateTracker = errors.New("Could not update tracker")

// ErrorYouDoNotHavePermissionsToSeeTheUsersOfThisTracker ...
var ErrorYouDoNotHavePermissionsToSeeTheUsersOfThisTracker = errors.New("You dont have permission to see these tracker users")

// TrackerService ...
type TrackerService interface {
	FindByUser(sub string) ([]model.Tracker, error)

	FindByID(sub string, id int64) (model.Tracker, error)

	CreateTracker(sub string, tracker model.Tracker) (model.Tracker, error)

	UpdateTracker(sub string, tracker model.Tracker) (model.Tracker, error)

	DeleteTracker(sub string, id int64) (bool, error)

	FindUsersForTracker(sub string, id int64) ([]model.User, error)
}

//GoDutchTrackerService ...
type GoDutchTrackerService struct {
	trackerRepository   trackerrepository.TrackerRepository
	userService         userservice.UserService
	logger              infrastructure.Logger
	validator           trackervalidation.TrackerValidator
	transferService     transferservice.TransferService
	spendSummaryService spendsummaryservice.SpendSummaryService
	spendRepository     spendrepository.SpendRepository
}

// NewGoDutchTrackerService ...
func NewGoDutchTrackerService(trackerRepository trackerrepository.TrackerRepository,
	userService userservice.UserService,
	logger infrastructure.Logger,
	validator trackervalidation.TrackerValidator,
	transferService transferservice.TransferService,
	spendSummaryService spendsummaryservice.SpendSummaryService,
	spendRepository spendrepository.SpendRepository) *GoDutchTrackerService {
	service := GoDutchTrackerService{}
	service.trackerRepository = trackerRepository
	service.userService = userService
	service.logger = logger
	service.validator = validator
	service.transferService = transferService
	service.spendSummaryService = spendSummaryService
	service.spendRepository = spendRepository
	return &service
}

// FindByUser ...
func (goDutchTrackerService *GoDutchTrackerService) FindByUser(sub string) ([]model.Tracker, error) {

	adminUser, err := goDutchTrackerService.userService.FindBySub(sub)
	if err != nil {
		return []model.Tracker{}, err
	}

	return goDutchTrackerService.trackerRepository.GetForUserID(adminUser.ID)
}

// FindUsersForTracker ...
func (goDutchTrackerService *GoDutchTrackerService) FindUsersForTracker(sub string, id int64) ([]model.User, error) {

	adminUser, err := goDutchTrackerService.userService.FindBySub(sub)
	if err != nil {
		return []model.User{}, err
	}

	tracker, err := goDutchTrackerService.trackerRepository.GetByID(id)
	if err != nil {
		return []model.User{}, err
	}

	if !infrastructure.Ints64Contains(tracker.TrackerUserIDs, adminUser.ID) {
		return []model.User{}, ErrorYouDoNotHavePermissionsToSeeTheUsersOfThisTracker
	}

	users := []model.User{}

	for _, u := range tracker.TrackerUserIDs {

		user, err := goDutchTrackerService.userService.FindByIDForTracker(u)
		if err != nil {
			return []model.User{}, err
		}
		users = append(users, user)
	}

	return users, nil
}

// FindByID ...
func (goDutchTrackerService *GoDutchTrackerService) FindByID(sub string, id int64) (model.Tracker, error) {

	_, err := goDutchTrackerService.userService.FindBySub(sub)
	if err != nil {
		return model.Tracker{}, err
	}

	return goDutchTrackerService.trackerRepository.GetByID(id)
}

// CreateTracker ...
func (goDutchTrackerService *GoDutchTrackerService) CreateTracker(sub string,
	tracker model.Tracker) (model.Tracker, error) {

	adminUser, err := goDutchTrackerService.userService.FindBySub(sub)
	if err != nil {
		return model.Tracker{}, err
	}

	tracker.DateCreated = time.Now()

	valid, err := goDutchTrackerService.validator.IsValidCreateTracker(tracker, goDutchTrackerService.logger, adminUser)
	if valid == false {
		return model.Tracker{}, err
	}

	tracker, err = goDutchTrackerService.trackerRepository.Insert(tracker)
	if err != nil {
		return model.Tracker{}, err
	}

	_, err = goDutchTrackerService.transferService.UpsertTransfers(tracker.ID)
	if err != nil {
		return model.Tracker{}, err
	}

	_, err = goDutchTrackerService.spendSummaryService.UpsertSpendSummaries(tracker.ID)
	if err != nil {
		return model.Tracker{}, err
	}

	return tracker, nil
}

// UpdateTracker ...
func (goDutchTrackerService *GoDutchTrackerService) UpdateTracker(sub string,
	tracker model.Tracker) (model.Tracker, error) {

	adminUser, err := goDutchTrackerService.userService.FindBySub(sub)
	if err != nil {
		return model.Tracker{}, err
	}

	existingTracker, err := goDutchTrackerService.trackerRepository.GetByID(tracker.ID)
	if err != nil {
		return model.Tracker{}, err
	}

	valid, err := goDutchTrackerService.validator.IsValidUpdateTracker(tracker, goDutchTrackerService.logger, adminUser, existingTracker)
	if valid == false {
		return model.Tracker{}, err
	}

	tracker, err = goDutchTrackerService.trackerRepository.Update(tracker.ID, tracker)
	if err != nil {
		return model.Tracker{}, err
	}

	_, err = goDutchTrackerService.transferService.UpsertTransfers(tracker.ID)
	if err != nil {
		return model.Tracker{}, err
	}

	_, err = goDutchTrackerService.spendSummaryService.UpsertSpendSummaries(tracker.ID)
	if err != nil {
		return model.Tracker{}, err
	}

	return tracker, nil
}

// DeleteTracker ...
func (goDutchTrackerService *GoDutchTrackerService) DeleteTracker(sub string,
	id int64) (bool, error) {

	adminUser, err := goDutchTrackerService.userService.FindBySub(sub)
	if err != nil {
		return false, err
	}

	existingTracker, err := goDutchTrackerService.trackerRepository.GetByID(id)
	if err != nil {
		return false, err
	}

	valid, err := goDutchTrackerService.validator.IsValidDeleteTracker(id, goDutchTrackerService.logger, adminUser, existingTracker)
	if valid == false {
		return false, err
	}

	_, err = goDutchTrackerService.transferService.DeleteTransfers(id)
	if err != nil {
		return false, err
	}

	_, err = goDutchTrackerService.spendSummaryService.DeleteSpendSummaries(id)
	if err != nil {
		return false, err
	}

	_, err = goDutchTrackerService.spendRepository.DeleteForTrackerID(id)
	if err != nil {
		return false, err
	}

	result, err := goDutchTrackerService.trackerRepository.Delete(id)
	if err != nil {
		return false, err
	}

	return result, nil
}
