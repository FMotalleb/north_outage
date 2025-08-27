package models

type Listener struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Hash string `gorm:"index:idx_listener_hash,unique" json:"unique_hash"`

	// ChatID
	TelegramCID int64 `gorm:"not null" json:"-"`
	// ThreadID
	TelegramTID int64 `gorm:"null" json:"-"`

	// RegexPattern
	SearchTerm string `gorm:"type:text;not null" json:"search_term"`
	City       string `gorm:"size:255;not null"  json:"city"`
}
