package models

import "time"

type Event struct {
	ID      uint   `gorm:"primaryKey"`
	Hash    string `gorm:"index:idx_event_hash,unique"`
	City    string `gorm:"size:255;not null"`
	Address string `gorm:"type:text;not null"`

	Start time.Time
	End   time.Time
}
