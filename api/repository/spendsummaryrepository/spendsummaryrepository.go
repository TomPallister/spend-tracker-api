package spendsummaryrepository

import (
	"errors"
	"database/sql"
	"github.com/TomPallister/godutch-api/api/model"
	"github.com/TomPallister/godutch-api/api/infrastructure"

)

// ErrorCouldNotFindSpendSummaries ...
var ErrorCouldNotFindSpendSummaries = errors.New("Could not find spend summaries")

// ErrorCouldNotInsertSummaries ...
var ErrorCouldNotInsertSummaries = errors.New("Could not insert summaries")

// SpendSummaryRepository ...
type SpendSummaryRepository interface {
	GetForTrackerID(id int64) ([]model.SpendSummary, error)
	Insert(spendSummaries []model.SpendSummary) ([]model.SpendSummary, error)
	Update(spendSummaries []model.SpendSummary) ([]model.SpendSummary, error)
	Delete(id int64) (bool, error)
}

// PostgresSpendSummaryRepository ...
type PostgresSpendSummaryRepository struct {
	logger           infrastructure.Logger
	db 				 *sql.DB
}

// NewPostgresSpendSummaryRepository ...
func NewPostgresSpendSummaryRepository(logger infrastructure.Logger,
	db *sql.DB) *PostgresSpendSummaryRepository {
	repository := PostgresSpendSummaryRepository{}
	repository.logger = logger
	repository.db = db
	return &repository
}

// GetForTrackerID ...
func (repo *PostgresSpendSummaryRepository) GetForTrackerID(id int64) ([]model.SpendSummary, error) {

	ssForTracker := []model.SpendSummary{}
	
	rows, err := repo.db.Query("SELECT \"ID\", \"TrackerID\", \"UserID\", \"Currency\", \"Value\" FROM \"SpendSummaries\" WHERE \"TrackerID\" = $1", id)
	if err != nil {
		return []model.SpendSummary{}, err
	}
			
	for rows.Next() {
		
		var ss model.SpendSummary

		err = rows.Scan(&ss.ID, &ss.TrackerID, &ss.UserID, &ss.Currency, &ss.Value)
		if err != nil { 
			return []model.SpendSummary{}, err
		}
		
		ssForTracker = append(ssForTracker, ss)
	}
	
	return ssForTracker, nil
}

// Insert ..
func (repo *PostgresSpendSummaryRepository) Insert(spendSummaries []model.SpendSummary) ([]model.SpendSummary, error) {

	for i := 0; i < len(spendSummaries); i++ {
		var lastInsertID int64
		
			err := repo. 
				db.
				QueryRow("INSERT INTO \"SpendSummaries\"( \"TrackerID\", \"UserID\", \"Currency\", \"Value\") VALUES ($1, $2, $3, $4) RETURNING \"ID\"",
				spendSummaries[i].TrackerID, spendSummaries[i].UserID, spendSummaries[i].Currency, spendSummaries[i].Value).Scan(&lastInsertID)
			
			if err != nil {
				return []model.SpendSummary{}, err
			}		
			
		spendSummaries[i].ID = lastInsertID
	}
		
	return spendSummaries, nil
}

// Update ...
func (repo *PostgresSpendSummaryRepository) Update(spendSummaries []model.SpendSummary) ([]model.SpendSummary, error) {

	for i := 0; i < len(spendSummaries); i++ {
		
		if spendSummaries[i].ID == 0 {
			
			var lastInsertID int64
		
			err := repo. 
				db.
				QueryRow("INSERT INTO \"SpendSummaries\"( \"TrackerID\", \"UserID\", \"Currency\", \"Value\") VALUES ($1, $2, $3, $4) RETURNING \"ID\"",
				spendSummaries[i].TrackerID, spendSummaries[i].UserID, spendSummaries[i].Currency, spendSummaries[i].Value).Scan(&lastInsertID)
			
			if err != nil {
				return []model.SpendSummary{}, err
			}		
			
			spendSummaries[i].ID = lastInsertID
		} else {
			
			stmt, err := repo.
				db.
				Prepare("UPDATE \"SpendSummaries\" SET \"TrackerID\"=$1, \"UserID\"=$2, \"Currency\"=$3, \"Value\"=$4 WHERE \"ID\" = $5")
				
			if err != nil { 
				return []model.SpendSummary{}, err
			}

			_, err = stmt.Exec(spendSummaries[i].TrackerID, spendSummaries[i].UserID, spendSummaries[i].Currency, spendSummaries[i].Value, spendSummaries[i].ID)
			if err != nil {
				return []model.SpendSummary{}, err
			}
		}	
	}
	return spendSummaries, nil
}


// Delete ...
func (repo *PostgresSpendSummaryRepository) Delete(id int64) (bool, error) {
    
	stmt, err := repo.db.Prepare("DELETE FROM \"SpendSummaries\" where \"TrackerID\"=$1")
    if err != nil {
		return false, err
	}

    _, err = stmt.Exec(id)
   	if err != nil {
		return false, err
	}
	
	return true, nil
}
