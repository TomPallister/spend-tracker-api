package route

import (
	"encoding/base64"
	"net/http"
	"os"

	"github.com/TomPallister/godutch-api/api/environment"
	"github.com/TomPallister/godutch-api/api/handler"
	"github.com/TomPallister/godutch-api/api/handler/spendhandler"
	"github.com/TomPallister/godutch-api/api/handler/spendsummarieshandler"
	"github.com/TomPallister/godutch-api/api/handler/trackerhandler"
	"github.com/TomPallister/godutch-api/api/handler/transferhandler"
	"github.com/TomPallister/godutch-api/api/handler/userhandler"
	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// GetRouter ...
func GetRouter(env *environment.Env) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			decoded, err := base64.URLEncoding.DecodeString(os.Getenv("AUTH0_CLIENT_SECRET"))
			if err != nil {
				return nil, err
			}
			return decoded, nil
		},
	})

	router.
		PathPrefix("/.well-known/acme-challenge/").
		Handler(http.StripPrefix("/.well-known/acme-challenge/", http.FileServer(http.Dir("./.well-known/acme-challenge/"))))

	router.Handle("/secured/ping", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(handler.SecuredPingHandler)),
	))

	router.HandleFunc("/", handler.PingHandler)

	//POST USERS
	router.Handle("/api/v1/users", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(userhandler.CreateUserHandler(env))),
	)).
		Methods("POST")

	//INVITE USERS
	router.Handle("/api/v1/users/invite", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(userhandler.InviteUserHandler(env))),
	)).
		Methods("POST")

	//ACCEPT INVITE USERS
	router.Handle("/api/v1/users/accept/{cryptoEmail}", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(userhandler.AcceptInviteUserHandler(env))),
	)).
		Methods("POST")

	// GET USERS/ID
	router.Handle("/api/v1/users/{id}", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(userhandler.FindUserByIDHandler(env))),
	)).
		Methods("GET")

	// GET USERS - hack to return current identity and tracker users
	router.Handle("/api/v1/users", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(userhandler.FindUserBySubHandler(env))),
	)).
		Methods("GET")

	// GET TRACKERS
	router.Handle("/api/v1/trackers", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(trackerhandler.FindTrackersBySubHandler(env))),
	)).
		Methods("GET")

	// GET TRACKERS ID
	router.Handle("/api/v1/trackers/{id}", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(trackerhandler.FindTrackersByIDHandler(env))),
	)).
		Methods("GET")

	// POST TRACKERS
	router.Handle("/api/v1/trackers", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(trackerhandler.CreateTrackerHandler(env))),
	)).
		Methods("POST")

	// PUT TRACKERS
	router.Handle("/api/v1/trackers/{id}", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(trackerhandler.UpdateTrackerHandler(env))),
	)).
		Methods("PUT")

	// DELTE TRACKERS
	router.Handle("/api/v1/trackers/{id}", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(trackerhandler.DeleteTrackerHandler(env))),
	)).
		Methods("DELETE")

	// GET SPENDS
	router.Handle("/api/v1/spends", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(spendhandler.FindFindByTrackerIDHandler(env))),
	)).
		Methods("GET")

	// POST SPENDS
	router.Handle("/api/v1/spends", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(spendhandler.CreateSpendHandler(env))),
	)).
		Methods("POST")

	// PUT SPENDS
	router.Handle("/api/v1/spends/{id}", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(spendhandler.UpdateSpendHandler(env))),
	)).
		Methods("PUT")

	// DELTE SPENDS
	router.Handle("/api/v1/spends/{id}", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(spendhandler.DeleteSpendHandler(env))),
	)).
		Methods("DELETE")

	// GET TRANSFERS
	router.Handle("/api/v1/transfers", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(transferhandler.FindFindByTrackerIDHandler(env))),
	)).
		Methods("GET")

	// GET SPEND SUMMARIES
	router.Handle("/api/v1/spendsummaries", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(spendsummarieshandler.FindFindByTrackerIDHandler(env))),
	)).
		Methods("GET")

		// GET SPEND SUMMARIES
	router.Handle("/api/v1/spendSummaries", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.Handler(spendsummarieshandler.FindFindByTrackerIDHandler(env))),
	)).
		Methods("GET")

	return router
}
