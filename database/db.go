package database

import (
	"gorm.io/gorm"

	"github.com/fmotalleb/north_outage/database/driver"
)

func NewDB(connection string) (*gorm.DB, error) {
	var conn gorm.Dialector
	var db *gorm.DB
	var err error
	if conn, err = driver.MakeConnection(connection); err != nil {
		return nil, err
	}
	if db, err = gorm.Open(conn, &gorm.Config{}); err != nil {
		return nil, err
	}
	return db, nil
}
