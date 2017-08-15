package transferrepository

import (
	"database/sql"
	"errors"

	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/TomPallister/godutch-api/api/model"
)

// ErrorCouldNotFindTransfers ...
var ErrorCouldNotFindTransfers = errors.New("Could not find transfers")

// ErrorCouldNotInsertTransfers ...
var ErrorCouldNotInsertTransfers = errors.New("Could not insert transfers")

// TransferRepository ...
type TransferRepository interface {
	GetForTrackerID(id int64) ([]model.Transfer, error)
	Insert(transfers []model.Transfer) ([]model.Transfer, error)
	Delete(id int64) (bool, error)
}

// PostgresTransferRepository ...
type PostgresTransferRepository struct {
	logger infrastructure.Logger
	db     *sql.DB
}

// NewPostgresTransferRepository ...
func NewPostgresTransferRepository(logger infrastructure.Logger,
	db *sql.DB) *PostgresTransferRepository {
	repository := PostgresTransferRepository{}
	repository.logger = logger
	repository.db = db
	return &repository
}

// GetForTrackerID ...
func (repo *PostgresTransferRepository) GetForTrackerID(id int64) ([]model.Transfer, error) {

	transfers := []model.Transfer{}

	rows, err := repo.db.Query("SELECT \"ID\", \"ToUserID\", \"FromUserID\", \"TrackerID\", \"Currency\", \"Value\" FROM \"Transfers\" WHERE \"TrackerID\" = $1", id)
	if err != nil {
		return []model.Transfer{}, err
	}

	for rows.Next() {

		var transfer model.Transfer

		err = rows.Scan(&transfer.ID, &transfer.ToUserID, &transfer.FromUserID, &transfer.TrackerID, &transfer.Currency, &transfer.Value)
		if err != nil {
			return []model.Transfer{}, err
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

// Insert ..
func (repo *PostgresTransferRepository) Insert(transfers []model.Transfer) ([]model.Transfer, error) {

	for i := 0; i < len(transfers); i++ {
		var lastInsertID int64

		err := repo.
			db.
			QueryRow("INSERT INTO \"Transfers\"(\"ToUserID\", \"FromUserID\", \"TrackerID\", \"Currency\", \"Value\") VALUES ($1, $2, $3, $4, $5) RETURNING \"ID\"",
				transfers[i].ToUserID, transfers[i].FromUserID, transfers[i].TrackerID, transfers[i].Currency, transfers[i].Value).Scan(&lastInsertID)

		if err != nil {
			return []model.Transfer{}, err
		}

		transfers[i].ID = lastInsertID
	}

	return transfers, nil
}

// Delete ...
func (repo *PostgresTransferRepository) Delete(id int64) (bool, error) {

	stmt, err := repo.db.Prepare("DELETE FROM \"Transfers\" where \"TrackerID\"=$1")
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return false, err
	}

	return true, nil
}
