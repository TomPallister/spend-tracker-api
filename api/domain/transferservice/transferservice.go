package transferservice

import (
	"errors"

	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/repository/spendrepository"
	"github.com/TomPallister/godutch-api/api/repository/trackerrepository"
	"github.com/TomPallister/godutch-api/api/repository/transferrepository"
	"github.com/TomPallister/godutch-api/api/repository/userrepository"
	"github.com/shopspring/decimal"
)

// ErrorCreateTransfers ...
var ErrorCreateTransfers = errors.New("Could not create transfers")

// ErrorFindTransfers ...
var ErrorFindTransfers = errors.New("Could not find transfers")

// TransferService ...Dont call this from anything that hasnt already been authenticated and authorised
type TransferService interface {
	UpsertTransfers(trackerID int64) ([]model.Transfer, error)
	FindTransfersForTrackerID(sub string, trackerID int64) ([]model.Transfer, error)
	DeleteTransfers(trackerID int64) (bool, error)
}

// GoDutchTransferService ...
type GoDutchTransferService struct {
	spendRepository    spendrepository.SpendRepository
	transferRepository transferrepository.TransferRepository
	trackerRepository  trackerrepository.TrackerRepository
	userRepository     userrepository.UserRepository
}

// NewGoDutchTransferService ...
func NewGoDutchTransferService(spendRepository spendrepository.SpendRepository,
	transferRepository transferrepository.TransferRepository,
	trackerRepository trackerrepository.TrackerRepository,
	userRepository userrepository.UserRepository) *GoDutchTransferService {
	service := GoDutchTransferService{}
	service.spendRepository = spendRepository
	service.transferRepository = transferRepository
	service.trackerRepository = trackerRepository
	service.userRepository = userRepository
	return &service
}

// UpsertTransfers ...
func (service *GoDutchTransferService) UpsertTransfers(trackerID int64) ([]model.Transfer, error) {

	transfers, err := makeTransfers(trackerID, service)
	if err != nil {
		return nil, err
	}

	existingTransfers, err := service.transferRepository.GetForTrackerID(trackerID)
	if err != nil {
		return nil, err
	}

	if len(existingTransfers) <= 0 {
		transfers, err = service.transferRepository.Insert(transfers)
		if err != nil {
			return nil, err
		}

		return transfers, nil
	}

	for _, t := range existingTransfers {
		_, err := service.transferRepository.Delete(t.TrackerID)
		if err != nil {
			return nil, err
		}
	}

	transfers, err = service.transferRepository.Insert(transfers)
	if err != nil {
		return nil, err
	}

	return transfers, nil
}

// DeleteTransfers ..
func (service *GoDutchTransferService) DeleteTransfers(trackerID int64) (bool, error) {
	return service.transferRepository.Delete(trackerID)
}

// FindTransfersForTrackerID ...
func (service *GoDutchTransferService) FindTransfersForTrackerID(sub string, trackerID int64) ([]model.Transfer, error) {

	user, err := service.userRepository.GetBySub(sub)
	if err != nil {
		return []model.Transfer{}, err
	}

	tracker, err := service.trackerRepository.GetByID(trackerID)
	if err != nil {
		return []model.Transfer{}, err
	}

	if !infrastructure.Ints64Contains(tracker.TrackerUserIDs, user.ID) {
		return []model.Transfer{}, ErrorFindTransfers
	}

	transfers, err := service.transferRepository.GetForTrackerID(trackerID)
	if err != nil {
		return []model.Transfer{}, nil
	}
	return transfers, nil
}

func makeTransfers(trackerID int64, service *GoDutchTransferService) ([]model.Transfer, error) {
	spends, err := service.spendRepository.GetForTrackerID(trackerID)
	if err != nil {
		return []model.Transfer{}, err
	}

	tracker, err := service.trackerRepository.GetByID(trackerID)
	if err != nil {
		return []model.Transfer{}, err
	}

	users, err := getUsersForTracker(tracker.TrackerUserIDs, service.userRepository)
	if err != nil {
		return []model.Transfer{}, err
	}

	totalSpend := sumOfSpends(spends)

	countOfUsers := float64(len(users))

	eachUsersShare := totalSpend.Div(decimal.NewFromFloat(countOfUsers))

	tracksGroupedByUser := getTracksGroupedByUser(users, spends)

	eachUserOwes := getUserOwed(tracksGroupedByUser, eachUsersShare)

	transfers, err := getTransfersToSettleThisTracker(eachUserOwes, tracker.Currency, trackerID, users)
	if err != nil {
		return []model.Transfer{}, err
	}

	return transfers, nil
}

func getTransfersToSettleThisTracker(userOweds []userOwed,
	currency string, trackerID int64, users []model.User) ([]model.Transfer, error) {

	var transfers = []model.Transfer{}

	for _, userOwed := range userOweds {

		otherUsersOwed := getOtherUsers(userOweds, userOwed)

		for _, otherUserOwed := range otherUsersOwed {

			if userDoesntOweAnyMoney(otherUserOwed) {

				amountUserOwes := userOwed.TotalOwed

				for amountUserOwes.Cmp(decimal.NewFromFloat(0)) == +1 {

					amountThisUserCouldPayCurrentOwedUser := decimal.Max(otherUserOwed.TotalOwed, amountUserOwes)

					transfer := model.Transfer{
						FromUserID: userOwed.UserID,
						ToUserID:   otherUserOwed.UserID,
						Value:      amountThisUserCouldPayCurrentOwedUser,
						Currency:   currency,
						TrackerID:  trackerID,
					}

					transfers = append(transfers, transfer)

					amountUserOwes = amountUserOwes.Sub(amountThisUserCouldPayCurrentOwedUser)
				}
			}
		}
	}

	return transfers, nil
}

func userDoesntOweAnyMoney(userOwed userOwed) bool {
	if userOwed.TotalOwed.Cmp(decimal.NewFromFloat(0)) == -1 {
		return true
	}
	return false
}

func getOtherUsers(userOweds []userOwed, thisUser userOwed) []userOwed {

	var otherUserOweds = []userOwed{}

	for _, u := range userOweds {
		if u.UserID != thisUser.UserID {
			otherUserOweds = append(otherUserOweds, u)
		}
	}

	return otherUserOweds
}

func getUserOwed(tracksGroupedByUser map[int64][]model.Spend, userShare decimal.Decimal) []userOwed {

	var trackerUserOwed = []userOwed{}

	for k, v := range tracksGroupedByUser {
		totalOwed := userShare.Sub(sumOfSpends(v))
		userOwed := userOwed{
			UserID:    k,
			TotalOwed: totalOwed,
		}
		trackerUserOwed = append(trackerUserOwed, userOwed)
	}

	return trackerUserOwed
}

func getTracksGroupedByUser(users []model.User, spends []model.Spend) map[int64][]model.Spend {

	var tracksGroupedByUser = make(map[int64][]model.Spend)

	for _, u := range users {

		var userSpends = []model.Spend{}

		for _, s := range spends {
			if s.UserID == u.ID {
				userSpends = append(userSpends, s)
			}
		}

		if len(userSpends) <= 0 {
			s := model.Spend{
				UserID: u.ID,
				Value:  decimal.NewFromFloat(0),
			}
			userSpends = append(userSpends, s)
		}

		tracksGroupedByUser[u.ID] = userSpends
	}

	return tracksGroupedByUser
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

func sumOfSpends(spends []model.Spend) decimal.Decimal {

	total := decimal.NewFromFloat(0)

	for _, s := range spends {
		total = total.Add(s.Value)
	}

	return total
}

type userOwed struct {
	UserID    int64
	TotalOwed decimal.Decimal
}
