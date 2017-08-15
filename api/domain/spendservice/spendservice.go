package spendservice

import (
	"errors"
	"time"

	"github.com/TomPallister/godutch-api/api/domain/spendsummaryservice"
	"github.com/TomPallister/godutch-api/api/domain/trackerservice"
	"github.com/TomPallister/godutch-api/api/domain/transferservice"
	"github.com/TomPallister/godutch-api/api/domain/userservice"
	"github.com/TomPallister/godutch-api/api/domain/validation/spendvalidation"
	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/repository/spendrepository"
)

// ErrorCreateSpend ...
var ErrorCreateSpend = errors.New("Could not create Spend")

// ErrorUpdateSpend ...
var ErrorUpdateSpend = errors.New("Could not update Spend")

// ErrorUserDoesNotBelongToTracker ...
var ErrorUserDoesNotBelongToTracker = errors.New("User does not belong to tracker")

// SpendService ...
type SpendService interface {
	FindByTrackerID(sub string, id int64) ([]model.Spend, error)

	CreateSpend(sub string, spend model.Spend) (model.Spend, error)

	UpdateSpend(sub string, spend model.Spend) (model.Spend, error)

	DeleteSpend(sub string, id int64) (bool, error)
}

//GoDutchSpendService ...
type GoDutchSpendService struct {
	spendRepository     spendrepository.SpendRepository
	userService         userservice.UserService
	trackerService      trackerservice.TrackerService
	logger              infrastructure.Logger
	validator           spendvalidation.SpendValidator
	transferService     transferservice.TransferService
	spendSummaryService spendsummaryservice.SpendSummaryService
}

// NewGoDutchSpendService ...
func NewGoDutchSpendService(spendRepository spendrepository.SpendRepository,
	userService userservice.UserService,
	trackerService trackerservice.TrackerService,
	validator spendvalidation.SpendValidator,
	logger infrastructure.Logger,
	transferService transferservice.TransferService,
	spendSummaryService spendsummaryservice.SpendSummaryService) *GoDutchSpendService {

	service := GoDutchSpendService{}
	service.spendRepository = spendRepository
	service.trackerService = trackerService 
	service.userService = userService
	service.validator = validator
	service.logger = logger
	service.transferService = transferService
	service.spendSummaryService = spendSummaryService
	return &service
}

// CreateSpend ...
func (goDutchSpendService *GoDutchSpendService) CreateSpend(sub string,
	spend model.Spend) (model.Spend, error) {

	user, err := goDutchSpendService.userService.FindBySub(sub)
	if err != nil {
		return model.Spend{}, err
	}

	tracker, err := goDutchSpendService.trackerService.FindByID(sub, spend.TrackerID)
	if err != nil {
		return model.Spend{}, err
	}

	spend.DateCreated = time.Now()

	valid, err := goDutchSpendService.validator.IsValidCreateSpend(spend, goDutchSpendService.logger, user, tracker)
	if valid == false {
		return model.Spend{}, err
	}

	spend, err = goDutchSpendService.spendRepository.Insert(spend)
	if err != nil {
		return model.Spend{}, err
	}

	_, err = goDutchSpendService.transferService.UpsertTransfers(tracker.ID)
	if err != nil {
		return model.Spend{}, err
	}

	_, err = goDutchSpendService.spendSummaryService.UpsertSpendSummaries(tracker.ID)
	if err != nil {
		return model.Spend{}, err
	}

	return spend, nil
}

// UpdateSpend ...
func (goDutchSpendService *GoDutchSpendService) UpdateSpend(sub string,
	spend model.Spend) (model.Spend, error) {

	user, err := goDutchSpendService.userService.FindBySub(sub)
	if err != nil {
		return model.Spend{}, err
	}

	existingSpend, err := goDutchSpendService.spendRepository.GetByID(spend.ID)
	if err != nil {
		return model.Spend{}, err
	}

	tracker, err := goDutchSpendService.trackerService.FindByID(sub, spend.TrackerID)
	if err != nil {
		return model.Spend{}, err
	}

	valid, err := goDutchSpendService.validator.IsValidUpdateSpend(spend, goDutchSpendService.logger, user, existingSpend, tracker)
	if valid == false {
		return model.Spend{}, err
	}

	spend, err = goDutchSpendService.spendRepository.Update(spend.ID, spend)
	if err != nil {
		return model.Spend{}, err
	}

	_, err = goDutchSpendService.transferService.UpsertTransfers(tracker.ID)
	if err != nil {
		return model.Spend{}, err
	}

	_, err = goDutchSpendService.spendSummaryService.UpsertSpendSummaries(tracker.ID)
	if err != nil {
		return model.Spend{}, err
	}

	return spend, nil

}

// DeleteSpend ...
func (goDutchSpendService *GoDutchSpendService) DeleteSpend(sub string,
	id int64) (bool, error) {

	user, err := goDutchSpendService.userService.FindBySub(sub)
	if err != nil {
		return false, err
	}

	existingSpend, err := goDutchSpendService.spendRepository.GetByID(id)
	if err != nil {
		return false, err
	}

	valid, err := goDutchSpendService.validator.IsValidDeleteSpend(id, goDutchSpendService.logger, user, existingSpend)
	if valid == false {
		return false, err
	} 

	result, err := goDutchSpendService.spendRepository.Delete(existingSpend.ID)
	if err != nil {
		return false, err
	}

	_, err = goDutchSpendService.transferService.UpsertTransfers(existingSpend.TrackerID)
	if err != nil {
		return false, err 
	}

	_, err = goDutchSpendService.spendSummaryService.UpsertSpendSummaries(existingSpend.TrackerID)
	if err != nil {
		return false, err
	}

	return result, nil
}

// FindByTrackerID ...
func (goDutchSpendService *GoDutchSpendService) FindByTrackerID(sub string,
	id int64) ([]model.Spend, error) {

	user, err := goDutchSpendService.userService.FindBySub(sub)
	if err != nil {
		return []model.Spend{}, err
	}

	tracker, err := goDutchSpendService.trackerService.FindByID(user.AuthenticationID, id)
	if err != nil {
		return []model.Spend{}, err
	}

	if !infrastructure.Ints64Contains(tracker.TrackerUserIDs, user.ID) {
		return []model.Spend{}, ErrorUserDoesNotBelongToTracker
	}

	return goDutchSpendService.spendRepository.GetForTrackerID(tracker.ID)
}
