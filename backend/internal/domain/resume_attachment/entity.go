package resume_attachment

import (
	"time"
)

type ResumeAttachment struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	VersionID   string `gorm:"type:uuid;not null;index"`
	FileName    string `gorm:"type:varchar(255);not null"`
	ContentType string `gorm:"type:varchar(100);not null"`
	FileSize    int64  `gorm:"not null"`
	Data        []byte `gorm:"type:bytea;not null"`
	CreatedAt   time.Time
}
