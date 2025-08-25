package models

type Listener struct {
	ID   uint   `gorm:"primaryKey"`
	Hash string `gorm:"index:idx_listener_hash,unique"`

	// ChatID
	TelegramCID int64 `gorm:"not null"`
	// ThreadID
	TelegramTID int64 `gorm:"null"`

	// RegexPattern
	SearchTerm string `gorm:"type:text;not null"`
}
