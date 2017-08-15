package userhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/TomPallister/godutch-api/api/environment"
	"github.com/TomPallister/godutch-api/api/handler"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/view"
)

// CreateUserHandler ...
func CreateUserHandler(env *environment.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subject, err := env.SubjectFinder.FindSubject(r, env.Logger)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusUnauthorized, w, env.Logger, err)
			return
		}
 
		var user view.User
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := r.Body.Close(); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := json.Unmarshal(body, &user); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		newUser, err := env.UserService.CreateUser(subject, user.User)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		userView := view.User{
			User: newUser,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(userView); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}

// InviteUserHandler ...
func InviteUserHandler(env *environment.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subject, err := env.SubjectFinder.FindSubject(r, env.Logger)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusUnauthorized, w, env.Logger, err)
			return
		}

		var inviteUser model.InviteUser
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := r.Body.Close(); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := json.Unmarshal(body, &inviteUser); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		goDutchRootURL := os.Getenv("GODUTCH_URL")
		if goDutchRootURL == "" {
			fmt.Println("Environment variable GODUTCH_URL is undefined.")
		}

		newUser, err := env.UserService.InviteUser(subject, inviteUser, goDutchRootURL)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		userView := view.User{
			User: newUser,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(userView); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}

// AcceptInviteUserHandler ...
func AcceptInviteUserHandler(env *environment.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subject, err := env.SubjectFinder.FindSubject(r, env.Logger)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusUnauthorized, w, env.Logger, err)
			return
		}

		var inviteUser view.User
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := r.Body.Close(); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := json.Unmarshal(body, &inviteUser); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		cryptoEmail, err := handler.GetCryptoEmailFromVars(r)
		newUser, err := env.UserService.AcceptInvite(subject, inviteUser.User, cryptoEmail)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		userView := view.User{
			User: newUser,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(userView); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}

// FindUserByIDHandler ...
func FindUserByIDHandler(env *environment.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subject, err := env.SubjectFinder.FindSubject(r, env.Logger)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusUnauthorized, w, env.Logger, err)
			return
		}

		id, err := handler.GetIDFromVARs(r)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		user, err := env.UserService.FindByID(subject, id)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		userView := view.User{
			User: user,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(userView); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}

// FindUserBySubHandler ...
func FindUserBySubHandler(env *environment.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subject, err := env.SubjectFinder.FindSubject(r, env.Logger)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusUnauthorized, w, env.Logger, err)
			return
		}

		trackerID := r.URL.Query()["trackerId"]

		if len(trackerID) <= 0 {

			user, err := env.UserService.FindBySub(subject)
			if err != nil {
				handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
				return
			}

			userView := view.User{
				User: user,
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(userView); err != nil {
				handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
				return
			}
			return
		}

		id, err := strconv.ParseInt(trackerID[0], 10, 64)
		if id > 0 {

			users, err := env.TrackerService.FindUsersForTracker(subject, id)
			if err != nil {
				handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
				return
			}

			viewUsers := view.Users{
				Users: users,
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(viewUsers); err != nil {
				handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
				return
			}
		}
	})
}
