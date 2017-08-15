package spendsummaryservice

import (
	"errors"

	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/repository/spendrepository"
	"github.com/TomPallister/godutch-api/api/repository/spendsummaryrepository"
	"github.com/TomPallister/godutch-api/api/repository/trackerrepository"
	"github.com/TomPallister/godutch-api/api/repository/userrepository"
	"github.com/shopspring/decimal"
)

// ErrorCreateSpendSummaries ...
var ErrorCreateSpendSummaries = errors.New("Could not create spend summaries")

// ErrorFindSpendSummaries ...
var ErrorFindSpendSummaries = errors.New("Could not find spend summaries")

// SpendSummaryService ...Dont call this from anything that hasnt already been authenticated and authorised
type SpendSummaryService interface {
	UpsertSpendSummaries(trackerID int64) ([]model.SpendSummary, error)
	FindSpendSummariesForTrackerID(sub string, trackerID int64) ([]model.SpendSummary, error)
	DeleteSpendSummaries(trackerID int64) (bool, error)
}

// GoDutchSpendSummaryService ...
type GoDutchSpendSummaryService struct {
	spendRepository        spendrepository.SpendRepository
	spendSummaryRepository spendsummaryrepository.SpendSummaryRepository
	trackerRepository      trackerrepository.TrackerRepository
	userRepository         userrepository.UserRepository
}

// NewGoDutchSpendSummaryService ...
func NewGoDutchSpendSummaryService(spendRepository spendrepository.SpendRepository,
	spendSummaryRepository spendsummaryrepository.SpendSummaryRepository,
	trackerRepository trackerrepository.TrackerRepository,
	userRepository userrepository.UserRepository) *GoDutchSpendSummaryService {
	service := GoDutchSpendSummaryService{}
	service.spendRepository = spendRepository
	service.spendSummaryRepository = spendSummaryRepository
	service.trackerRepository = trackerRepository
	service.userRepository = userRepository
	return &service
}

// UpsertSpendSummaries ...
func (service *GoDutchSpendSummaryService) UpsertSpendSummaries(trackerID int64) ([]model.SpendSummary, error) {

	spendSummaries, err := makeSpendSummaries(trackerID, service)
	if err != nil {
		return nil, err
	}

	existingSpendSummaries, err := service.spendSummaryRepository.GetForTrackerID(trackerID)
	if err != nil {
		return nil, err
	}

	if len(existingSpendSummaries) <= 0 {
		spendSummaries, err = service.spendSummaryRepository.Insert(spendSummaries)
		if err != nil {
			return nil, err
		}

		return spendSummaries, nil
	}

	summariesToUpdate := matchNewSummariesWithOld(existingSpendSummaries, spendSummaries)

	spendSummaries, err = service.spendSummaryRepository.Update(summariesToUpdate)
	if err != nil {
		return nil, err
	}

	return spendSummaries, nil
}

func matchNewSummariesWithOld(existing []model.SpendSummary, updated []model.SpendSummary) []model.SpendSummary {

	updatedMatchedWithExisting := []model.SpendSummary{}

	for _, uS := range updated {
		matched := false
		for _, eS := range existing {
			if uS.TrackerID == eS.TrackerID && uS.UserID == eS.UserID {
				eS.Value = uS.Value
				eS.Currency = uS.Currency
				updatedMatchedWithExisting = append(updatedMatchedWithExisting, eS)
				matched = true
				break
			}
		}
		if matched {
			continue
		}
		updatedMatchedWithExisting = append(updatedMatchedWithExisting, uS)
	}

	return updatedMatchedWithExisting
}

// DeleteSpendSummaries ..
func (service *GoDutchSpendSummaryService) DeleteSpendSummaries(trackerID int64) (bool, error) {
	return service.spendSummaryRepository.Delete(trackerID)
}

// FindSpendSummariesForTrackerID ...
func (service *GoDutchSpendSummaryService) FindSpendSummariesForTrackerID(sub string, trackerID int64) ([]model.SpendSummary, error) {

	user, err := service.userRepository.GetBySub(sub)
	if err != nil {
		return []model.SpendSummary{}, err
	}

	tracker, err := service.trackerRepository.GetByID(trackerID)
	if err != nil {
		return []model.SpendSummary{}, err
	}

	if !infrastructure.Ints64Contains(tracker.TrackerUserIDs, user.ID) {
		return []model.SpendSummary{}, ErrorFindSpendSummaries
	}

	spendSummaries, err := service.spendSummaryRepository.GetForTrackerID(trackerID)
	if err != nil {
		return []model.SpendSummary{}, nil
	}
	return spendSummaries, nil
}

func makeSpendSummaries(trackerID int64, service *GoDutchSpendSummaryService) ([]model.SpendSummary, error) {
	spends, err := service.spendRepository.GetForTrackerID(trackerID)
	if err != nil {
		return []model.SpendSummary{}, err
	}

	tracker, err := service.trackerRepository.GetByID(trackerID)
	if err != nil {
		return []model.SpendSummary{}, err
	}

	users, err := getUsersForTracker(tracker.TrackerUserIDs, service.userRepository)
	if err != nil {
		return []model.SpendSummary{}, err
	}

	spendSummaries := []model.SpendSummary{}

	for _, u := range users {
		totalSpend := getUsersTotalSpend(spends, u.ID)
		spendSummary := model.SpendSummary{
			TrackerID:    tracker.ID,
			UserID:       u.ID,
			Value:        totalSpend,
			Currency:     tracker.Currency,
		}
		spendSummaries = append(spendSummaries, spendSummary)
	}

	return spendSummaries, nil
}

func getUsersForTracker(trackerUserIDs []int64, repo userrepository.UserRepository) ([]model.User, error) {
	var users = []model.User{}

	for i := 0; i < len(trackerUserIDs); i++ {
		usr, err := repo.GetByID(trackerUserIDs[i])
		if err != nil {
			return users, nil
		}
		users = append(users, usr)
	}

	return users, nil
}

func getUsersTotalSpend(spends []model.Spend, userID int64) decimal.Decimal {
	total := decimal.NewFromFloat(0)

	for _, s := range spends {
		if s.UserID == userID {
			total = total.Add(s.Value)
		}
	}

	return total
}
