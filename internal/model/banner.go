package model

import "time"

type Banner struct {
	ID        int
	FeatureID int
	Content   string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
