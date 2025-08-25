package driver

import (
	"net/url"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//# go:build orm-sqllite

func init() {
	builders["sqlite"] = sqliteBuilder
}

func sqliteBuilder(c *url.URL) (gorm.Dialector, error) {
	return sqlite.Open(c.Path), nil
}
