package environment

import (
	"github.com/TomPallister/godutch-api/api/domain/spendservice"
	"github.com/TomPallister/godutch-api/api/domain/spendsummaryservice"
	"github.com/TomPallister/godutch-api/api/domain/trackerservice"
	"github.com/TomPallister/godutch-api/api/domain/transferservice"
	"github.com/TomPallister/godutch-api/api/domain/userservice"
	"github.com/TomPallister/godutch-api/api/infrastructure"
)

// Env ...
type Env struct {
	Logger              infrastructure.Logger
	UserService         userservice.UserService
	SubjectFinder       infrastructure.SubjectFinder
	TrackerService      trackerservice.TrackerService
	SpendService        spendservice.SpendService
	TransferService     transferservice.TransferService
	SpendSummaryService spendsummaryservice.SpendSummaryService
}
