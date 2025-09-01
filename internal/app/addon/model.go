package addon

import "time"

type Addon struct {
	ID        uint      `json:"id" db:"id"`
	ProjectID uint      `json:"project_id" db:"project_id"`
	Name      string    `json:"name" db:"name"`
	Type      string    `json:"type" db:"type"`
	Tier      string    `json:"tier" db:"tier"`
	Storage   string    `json:"storage" db:"storage"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}