package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Event struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Hash    string `gorm:"index:idx_event_hash,unique" json:"unique_hash"`
	City    string `gorm:"size:255;not null" json:"city"`
	Address string `gorm:"type:text;not null" json:"address"`

	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	CreatedAt time.Time `json:"created_at"`
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
