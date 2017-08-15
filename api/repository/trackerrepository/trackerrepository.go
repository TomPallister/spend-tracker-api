package trackerrepository

import (
	"errors"
	
	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
	"database/sql"

)

// ErrorNotFound ...to be used when object does not exist in repository
var ErrorNotFound = errors.New("Tracker not found")

// ErrorCouldNotInsertTracker ...
var ErrorCouldNotInsertTracker = errors.New("Could not insert tracker")

// ErrorCouldNotUpdateTracker ...
var ErrorCouldNotUpdateTracker = errors.New("Could not update tracker")

// TrackerRepository ...
type TrackerRepository interface { 
	GetByID(id int64) (model.Tracker, error)
	GetForUserID(id int64) ([]model.Tracker, error)
	Insert(tracker model.Tracker) (model.Tracker, error)
	Update(id int64, tracker model.Tracker) (model.Tracker, error)
	Delete(id int64) (bool, error)
}

// PostgresTrackerRepository ... 
type PostgresTrackerRepository struct {
	logger infrastructure.Logger
	db 	   *sql.DB
}

// NewPostgresTrackerRepository ...
func NewPostgresTrackerRepository(logger infrastructure.Logger,
	db *sql.DB) *PostgresTrackerRepository {
	repository := PostgresTrackerRepository{}
	repository.logger = logger
	repository.db = db
	return &repository
}

// GetByID ...
func (repository *PostgresTrackerRepository) GetByID(id int64) (model.Tracker, error) {
	
	var repoTracker model.Tracker

    err := repository.db.QueryRow("SELECT \"ID\", \"AdminUserID\", \"Name\", \"Currency\", \"DateCreated\" FROM \"Trackers\" WHERE \"ID\" = $1", id).
		Scan(&repoTracker.ID, &repoTracker.AdminUserID, &repoTracker.Name, &repoTracker.Currency, &repoTracker.DateCreated)
    
	switch {
    case err == sql.ErrNoRows:
            return model.Tracker{}, err
    case err != nil:
            return model.Tracker{}, err
	}
    
	rows, err := repository.db.Query("SELECT \"UserID\" FROM \"TrackerUsers\" WHERE \"TrackerID\" = $1", id)
	if err != nil {
		return model.Tracker{}, err
	}
	
	trackerUserIds := []int64{}
	
	for rows.Next() {
		
        var userID int64
  
        err = rows.Scan(&userID)
       	if err != nil {
			return model.Tracker{}, err
		}
		
		trackerUserIds = append(trackerUserIds, userID)
    }
	
	repoTracker.TrackerUserIDs = trackerUserIds
	 
	return repoTracker, nil
}

// GetForUserID ...
func (repository *PostgresTrackerRepository) GetForUserID(id int64) ([]model.Tracker, error) {
	
	rows, err := repository.db.Query("SELECT \"Trackers\".\"ID\", \"AdminUserID\", \"Name\", \"Currency\", \"DateCreated\" FROM \"Trackers\" INNER JOIN \"TrackerUsers\" ON \"Trackers\".\"ID\"=\"TrackerUsers\".\"TrackerID\" WHERE \"TrackerUsers\".\"UserID\" = $1", id)
	if err != nil {
		return []model.Tracker{}, err
	}
	
	trackersForUser := []model.Tracker{}
	
	for rows.Next() {
		
        var tracker model.Tracker
  
        err = rows.Scan(&tracker.ID, &tracker.AdminUserID, &tracker.Name, &tracker.Currency, &tracker.DateCreated)
       	if err != nil {
			return []model.Tracker{}, err
		}
		
		
		tURows, err := repository.db.Query("SELECT \"UserID\" FROM \"TrackerUsers\" WHERE \"TrackerID\" = $1", tracker.ID)
		if err != nil {
			return []model.Tracker{}, err
		}
		
		trackerUserIds := []int64{}
		
		for tURows.Next() {
			
			var userID int64
	
			err = tURows.Scan(&userID)
			if err != nil {
				return []model.Tracker{}, err
			}
			
			trackerUserIds = append(trackerUserIds, userID)
		}
		
		tracker.TrackerUserIDs = trackerUserIds
		
		trackersForUser = append(trackersForUser, tracker)
    }

	return trackersForUser, nil
}

// Insert ...
func (repository *PostgresTrackerRepository) Insert(tracker model.Tracker) (model.Tracker, error) {
	
	var lastInsertID int64
	err := repository.
		db.
		QueryRow("INSERT INTO \"Trackers\"(\"AdminUserID\", \"Name\", \"Currency\", \"DateCreated\") VALUES ($1, $2, $3, $4) RETURNING \"ID\"",
		tracker.AdminUserID, tracker.Name, tracker.Currency, tracker.DateCreated).Scan(&lastInsertID)
	if err != nil {
		return model.Tracker{}, err
	}		
			
	tracker.ID = lastInsertID
	
	for _, u := range tracker.TrackerUserIDs {
		
		stmt, err := repository.db.Prepare("INSERT INTO \"TrackerUsers\"(\"TrackerID\", \"UserID\") VALUES ($1, $2)")
		if err != nil {
			return model.Tracker{}, err
		}	
		
		_, err = stmt.Exec(tracker.ID, u)
		if err != nil {
			return model.Tracker{}, err
		}	
	}
		
	return tracker, nil
}

// Update ...
func (repository *PostgresTrackerRepository) Update(id int64, tracker model.Tracker) (model.Tracker, error) {
	
    stmt, err := repository.db.Prepare("UPDATE \"Trackers\" SET \"AdminUserID\"=$1, \"Name\"=$2, \"Currency\"=$3, \"DateCreated\"=$4 WHERE \"ID\" = $5")
	if err != nil {
		return model.Tracker{}, err
	}

    _, err = stmt.Exec(tracker.AdminUserID, tracker.Name, tracker.Currency, tracker.DateCreated, id)
	if err != nil {
		return model.Tracker{}, err
	}
	
	stmt, err = repository.db.Prepare("DELETE from \"TrackerUsers\" WHERE \"TrackerID\" = $1")
	if err != nil {
		return model.Tracker{}, err
	}
	
    _, err = stmt.Exec(id)
	if err != nil {
		return model.Tracker{}, err
	}	
	
	for _, u := range tracker.TrackerUserIDs {
		
		stmt, err := repository.db.Prepare("INSERT INTO \"TrackerUsers\"(\"TrackerID\", \"UserID\") VALUES ($1, $2)")
		if err != nil {
			return model.Tracker{}, err
		}	
		
		_, err = stmt.Exec(tracker.ID, u)
		if err != nil {
			return model.Tracker{}, err
		}	
	}
	
	return tracker, nil
}

// Delete ...
func (repository *PostgresTrackerRepository) Delete(id int64) (bool, error) {
	
    stmt, err := repository.db.Prepare("DELETE FROM \"TrackerUsers\" where \"TrackerID\"=$1")
    if err != nil {
		return false, err
	}

    _, err = stmt.Exec(id)
   	if err != nil {
		return false, err
	}
	
	 stmt, err = repository.db.Prepare("DELETE FROM \"Trackers\" where \"ID\"=$1")
    if err != nil {
		return false, err
	}

    _, err = stmt.Exec(id)
   	if err != nil {
		return false, err
	}
	
	return true, nil
}