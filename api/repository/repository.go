package repository

import (
    "database/sql"
    "fmt"
    // llalalal
    _ "github.com/lib/pq"
)

// NewDB ...
func NewDB(dataSourceName string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dataSourceName)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        fmt.Println(err)
        return nil, err
    }
    return db, nil
}