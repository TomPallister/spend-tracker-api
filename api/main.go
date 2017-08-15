package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/TomPallister/godutch-api/api/domain/spendservice"
	"github.com/TomPallister/godutch-api/api/domain/spendsummaryservice"
	"github.com/TomPallister/godutch-api/api/domain/trackerservice"
	"github.com/TomPallister/godutch-api/api/domain/transferservice"
	"github.com/TomPallister/godutch-api/api/domain/userservice"
	"github.com/TomPallister/godutch-api/api/domain/validation/spendvalidation"
	"github.com/TomPallister/godutch-api/api/domain/validation/trackervalidation"
	"github.com/TomPallister/godutch-api/api/domain/validation/uservalidation"
	"github.com/TomPallister/godutch-api/api/environment"
	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/repository"
	"github.com/TomPallister/godutch-api/api/repository/spendrepository"
	"github.com/TomPallister/godutch-api/api/repository/spendsummaryrepository"
	"github.com/TomPallister/godutch-api/api/repository/trackerrepository"
	"github.com/TomPallister/godutch-api/api/repository/transferrepository"
	"github.com/TomPallister/godutch-api/api/repository/userrepository"
	"github.com/TomPallister/godutch-api/api/route"
	"github.com/codegangsta/negroni"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// ErrorDatabaseConnectionString ...
var ErrorDatabaseConnectionString = errors.New("Environment variable PGSQLCONNECTIONSTRING is undefined.")

// ErrorAPIPort ...
var ErrorAPIPort = errors.New("Environment variable API_PORT is undefined.")

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	startServer()
}

func startServer() {

	connectionString := os.Getenv("PGSQL_CONNECTIONSTRING")
	if connectionString == "" {
		panic(ErrorDatabaseConnectionString)
	}

	db, err := repository.NewDB(connectionString)
	if err != nil {
		panic(err)
	}
	var emailService = &infrastructure.SendGridEmailService{}
	var logger = infrastructure.ConsoleLogger{}
	var trackerRepository = trackerrepository.NewPostgresTrackerRepository(logger, db)
	var userRepository = userrepository.NewPostgresUserRepository(logger, db)
	var spendRepository = spendrepository.NewPostgresSpendRepository(logger, db)
	var transferRepository = transferrepository.NewPostgresTransferRepository(logger, db)
	var spendSummaryRepository = spendsummaryrepository.NewPostgresSpendSummaryRepository(logger, db)
	var spendSummaryService = spendsummaryservice.
		NewGoDutchSpendSummaryService(spendRepository, spendSummaryRepository, trackerRepository, userRepository)
	var transferService = transferservice.
		NewGoDutchTransferService(spendRepository, transferRepository, trackerRepository, userRepository)
	var userService = userservice.
		NewGoDutchUserService(userRepository, uservalidation.NewGoDutchUserValidator(), logger, emailService, trackerRepository)
	var trackerService = trackerservice.
		NewGoDutchTrackerService(trackerRepository, userService, logger, trackervalidation.NewGoDutchTrackerValidator(), transferService, spendSummaryService, spendRepository)
	var spendService = spendservice.
		NewGoDutchSpendService(spendRepository, userService, trackerService, spendvalidation.NewGoDutchSpendValidator(), logger, transferService, spendSummaryService)

	env := &environment.Env{
		Logger:              logger,
		UserService:         userService,
		SubjectFinder:       infrastructure.NewGorrilaAndJwtSubjectFinder(),
		TrackerService:      trackerService,
		SpendService:        spendService,
		TransferService:     transferService,
		SpendSummaryService: spendSummaryService,
	}

	router := route.GetRouter(env)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	})

	n := negroni.Classic()
	n.Use(c)
	n.UseHandler(router)

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		panic(ErrorAPIPort)
	}

	http.ListenAndServe(apiPort, n)
}
