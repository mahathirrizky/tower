package models

import (
	"time"

	"gorm.io/gorm"
)

type TowerEvent struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TowerID     uint           `gorm:"index" json:"tower_id"` // Foreign key to the Tower
	EventType   string         `gorm:"type:varchar(50)" json:"event_type"` // e.g., "OwnershipChange", "Dismantled", "Relocation", "Created"
	Timestamp   time.Time      `json:"timestamp"`
	Description string         `gorm:"type:text" json:"description"` // Human-readable summary
	OldData     string         `gorm:"type:jsonb" json:"old_data"` // JSON string of relevant old data
	NewData     string         `gorm:"type:jsonb" json:"new_data"` // JSON string of relevant new data
	UserID      uint           `json:"user_id"` // Who performed the action (if applicable)
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
