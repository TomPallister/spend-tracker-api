package infrastructure

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// ErrorCouldNotFindSubjectClaim ...
var ErrorCouldNotFindSubjectClaim = errors.New("Could not find subject claim")

// ErrorSubjectWasNotAString ...
var ErrorSubjectWasNotAString = errors.New("Subject was not a string")

// SubjectFinder ...
type SubjectFinder interface {
	FindSubject(r *http.Request, logger Logger) (string, error)
}

// GorrilaAndJwtSubjectFinder ...
type GorrilaAndJwtSubjectFinder struct {
}

// NewGorrilaAndJwtSubjectFinder ...
func NewGorrilaAndJwtSubjectFinder() *GorrilaAndJwtSubjectFinder {

	service := GorrilaAndJwtSubjectFinder{}

	return &service
}

// FindSubject ...
func (g *GorrilaAndJwtSubjectFinder) FindSubject(r *http.Request, logger Logger) (string, error) {
	u := context.Get(r, "user")
	sub, ok := u.(*jwt.Token).Claims["sub"]
	if ok == false {
		return "", ErrorCouldNotFindSubjectClaim
	}
	if str, ok := sub.(string); ok {
		return str, nil
	}
	return "", ErrorSubjectWasNotAString
}
