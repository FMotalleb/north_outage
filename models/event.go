package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Event struct {
	ID      uint   `gorm:"primaryKey"`
	Hash    string `gorm:"index:idx_event_hash,unique"`
	City    string `gorm:"size:255;not null"`
	Address string `gorm:"type:text;not null"`

	Start time.Time
	End   time.Time
}

func (e *Event) ResetHash() {
	data := fmt.Sprintf("%s|%s|%d|%d",
		e.City,
		e.Address,
		e.Start.UnixNano(),
		e.End.UnixNano(),
	)

	hash := sha256.Sum256([]byte(data))
	e.Hash = hex.EncodeToString(hash[:])
}
