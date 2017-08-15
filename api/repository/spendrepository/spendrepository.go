package spendrepository

import (
	"database/sql"
	"errors"

	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
)

// ErrorNotFound ...to be used when object does not exist in repository
var ErrorNotFound = errors.New("Spend not found")

// ErrorCouldNotInsertSpend ...
var ErrorCouldNotInsertSpend = errors.New("Could not insert spend")

// ErrorCouldNotUpdateSpend ...
var ErrorCouldNotUpdateSpend = errors.New("Could not update spend")

// SpendRepository ...
type SpendRepository interface {
	GetByID(id int64) (model.Spend, error)
	GetForTrackerID(id int64) ([]model.Spend, error)
	Insert(spend model.Spend) (model.Spend, error)
	Update(id int64, spend model.Spend) (model.Spend, error)
	Delete(id int64) (bool, error)
	DeleteForTrackerID(id int64) (bool, error)
}

// PostgresSpendRepository ...
type PostgresSpendRepository struct {
	logger infrastructure.Logger
	db     *sql.DB
}

// NewPostgresSpendRepository ...
func NewPostgresSpendRepository(logger infrastructure.Logger,
	db *sql.DB) *PostgresSpendRepository {
	repository := PostgresSpendRepository{}
	repository.logger = logger
	repository.db = db
	return &repository
}

// GetByID ...
func (repository *PostgresSpendRepository) GetByID(id int64) (model.Spend, error) {

	repoSpend := model.Spend{}

	err := repository.db.QueryRow("SELECT \"ID\", \"TrackerID\", \"UserID\", \"Name\", \"DateCreated\", \"Value\", \"Currency\" FROM \"Spends\" WHERE \"ID\" = $1", id).
		Scan(&repoSpend.ID, &repoSpend.TrackerID, &repoSpend.UserID, &repoSpend.Name, &repoSpend.DateCreated, &repoSpend.Value, &repoSpend.Currency)

	switch {
	case err == sql.ErrNoRows:
		return model.Spend{}, err
	case err != nil:
		return model.Spend{}, err
	}

	return repoSpend, nil
}

// GetForTrackerID ...
func (repository *PostgresSpendRepository) GetForTrackerID(id int64) ([]model.Spend, error) {

	spendsForTracker := []model.Spend{}

	rows, err := repository.db.Query("SELECT \"ID\", \"TrackerID\", \"UserID\", \"Name\", \"DateCreated\", \"Value\", \"Currency\" FROM \"Spends\" WHERE \"TrackerID\" = $1", id)
	if err != nil {
		return []model.Spend{}, err
	}

	for rows.Next() {

		var spend model.Spend

		err = rows.Scan(&spend.ID, &spend.TrackerID, &spend.UserID, &spend.Name, &spend.DateCreated, &spend.Value, &spend.Currency)
		if err != nil {
			return []model.Spend{}, err
		}

		spendsForTracker = append(spendsForTracker, spend)
	}

	return spendsForTracker, nil
}

// Insert ...
func (repository *PostgresSpendRepository) Insert(spend model.Spend) (model.Spend, error) {

	var lastInsertID int64

	err := repository.
		db.
		QueryRow("INSERT INTO \"Spends\"(\"TrackerID\", \"UserID\", \"Name\", \"DateCreated\", \"Value\", \"Currency\") VALUES ($1, $2, $3, $4, $5, $6) RETURNING \"ID\"",
			spend.TrackerID, spend.UserID, spend.Name, spend.DateCreated, spend.Value, spend.Currency).Scan(&lastInsertID)
	if err != nil {
		return model.Spend{}, err
	}

	spend.ID = lastInsertID

	return spend, nil

}

// Update ...
func (repository *PostgresSpendRepository) Update(id int64, spend model.Spend) (model.Spend, error) {

	stmt, err := repository.db.Prepare("UPDATE \"Spends\" SET \"TrackerID\"= $1, \"UserID\"= $2, \"Name\"= $3, \"DateCreated\"= $4, \"Value\"= $5, \"Currency\"= $6 WHERE \"ID\" = $7")
	if err != nil {
		return model.Spend{}, err
	}

	_, err = stmt.Exec(spend.TrackerID, spend.UserID, spend.Name, spend.DateCreated, spend.Value, spend.Currency, id)
	if err != nil {
		return model.Spend{}, err
	}

	return spend, nil
}

// Delete ...
func (repository *PostgresSpendRepository) Delete(id int64) (bool, error) {

	stmt, err := repository.db.Prepare("DELETE FROM \"Spends\" where \"ID\"=$1")
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteForTrackerID ...
func (repository *PostgresSpendRepository) DeleteForTrackerID(id int64) (bool, error) {

	stmt, err := repository.db.Prepare("DELETE FROM \"Spends\" where \"TrackerID\"=$1")
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return false, err
	}

	return true, nil
}
