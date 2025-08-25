package driver

import (
	"net/url"
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//# go:build orm-sqllite

func init() {
	builders["sqlite"] = sqliteBuilder
}

func sqliteBuilder(c *url.URL) (gorm.Dialector, error) {
	fullPath := path.Join(c.Hostname(), c.Path)
	return sqlite.Open(fullPath), nil
}
