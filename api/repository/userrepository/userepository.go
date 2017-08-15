package userrepository

import (
	"errors"

	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"database/sql"
)

// ErrorNotFound ...to be used when object does not exist in repository
var ErrorNotFound = errors.New("User not found")

// ErrorCouldNotInsertUser ...
var ErrorCouldNotInsertUser = errors.New("Could not insert user")

// UserRepository ...
type UserRepository interface {
	GetBySub(sub string) (model.User, error)
	GetByID(id int64) (model.User, error)
	GetByEmail(email string) (model.User, error)
	Insert(user model.User) (model.User, error)
	Update(id int64, user model.User) (model.User, error)
}

// PostgresUserRepository ..
type PostgresUserRepository struct {
	logger           infrastructure.Logger
	db 				 *sql.DB
}

// NewPostgresUserRepository ...
func NewPostgresUserRepository(logger infrastructure.Logger,
	db *sql.DB) *PostgresUserRepository {
	repository := PostgresUserRepository{}
	repository.logger = logger
	repository.db = db
	return &repository
}

// Insert ...
func (userRepository *PostgresUserRepository) Insert(user model.User) (model.User, error) {

	var lastInsertID int64
	err := userRepository.db.QueryRow("INSERT INTO \"Users\"(\"AuthenticationID\",\"DateCreated\",\"EmailAddress\", \"Name\") VALUES($1, $2, $3, $4) RETURNING \"ID\"",
		user.AuthenticationID, user.DateCreated, user.EmailAddress, user.Name).Scan(&lastInsertID)
	if err != nil {
		return model.User{}, err
	}		
			
	user.ID = lastInsertID
	
	return user, nil
}

// Update ...
func (userRepository *PostgresUserRepository) Update(id int64, user model.User) (model.User, error) {

    stmt, err := userRepository.db.Prepare("UPDATE \"Users\" SET \"Name\"= $1, \"DateCreated\"= $2, \"EmailAddress\"= $3, \"AuthenticationID\"= $4 WHERE \"ID\"= $5")
	if err != nil {
			return model.User{}, err
	}
	
    res, err := stmt.Exec(user.Name, user.DateCreated, user.EmailAddress, user.AuthenticationID, id)
	if err != nil {
		return model.User{}, err
	}
	
	_, err = res.RowsAffected()
	if err != nil {
		return model.User{}, err
	}
		
	return user, nil
}

// GetBySub ...
func (userRepository *PostgresUserRepository) GetBySub(sub string) (model.User, error) {

	repoUser := model.User{}

    err := userRepository.db.QueryRow("SELECT \"ID\", \"EmailAddress\", \"Name\", \"AuthenticationID\", \"DateCreated\" FROM \"Users\" WHERE \"AuthenticationID\" = $1", sub).
		Scan(&repoUser.ID, &repoUser.EmailAddress, &repoUser.Name, &repoUser.AuthenticationID, &repoUser.DateCreated)
    
	switch {
    case err == sql.ErrNoRows:
            return model.User{}, err
    case err != nil:
            return model.User{}, err
	}
    
	return repoUser, nil
}

// GetByID ...
func (userRepository *PostgresUserRepository) GetByID(id int64) (model.User, error) {
	
	repoUser := model.User{}

    err := userRepository.db.QueryRow("SELECT \"ID\", \"EmailAddress\", \"Name\", \"AuthenticationID\", \"DateCreated\" FROM \"Users\" WHERE \"ID\" = $1", id).
		Scan(&repoUser.ID, &repoUser.EmailAddress, &repoUser.Name, &repoUser.AuthenticationID, &repoUser.DateCreated)
    
	switch {
    case err == sql.ErrNoRows:
            return model.User{}, err
    case err != nil:
            return model.User{}, err
	}
    
	return repoUser, nil
}

// GetByEmail ...
func (userRepository *PostgresUserRepository) GetByEmail(email string) (model.User, error) {
 	
	repoUser := model.User{}

    err := userRepository.db.QueryRow("SELECT \"ID\", \"EmailAddress\", \"Name\", \"AuthenticationID\", \"DateCreated\" FROM \"Users\" WHERE \"EmailAddress\" = $1", email).
		Scan(&repoUser.ID, &repoUser.EmailAddress, &repoUser.Name, &repoUser.AuthenticationID, &repoUser.DateCreated)
    
	switch {
    case err == sql.ErrNoRows:
            return model.User{}, err
    case err != nil:
            return model.User{}, err
	}
    
	return repoUser, nil
}
