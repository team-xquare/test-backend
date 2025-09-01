package addon

import "time"

type CreateAddonRequest struct {
	Name    string      `json:"name" binding:"required"`
	Type    string      `json:"type" binding:"required"`
	Tier    string      `json:"tier" binding:"required"`
	Storage string      `json:"storage"`
	Config  interface{} `json:"config"`
}

type UpdateAddonRequest struct {
	Name    string      `json:"name" binding:"required"`
	Type    string      `json:"type" binding:"required"`
	Tier    string      `json:"tier" binding:"required"`
	Storage string      `json:"storage"`
	Config  interface{} `json:"config"`
}

type AddonResponse struct {
	ID        uint        `json:"id"`
	ProjectID uint        `json:"project_id"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Tier      string      `json:"tier"`
	Storage   string      `json:"storage"`
	Config    interface{} `json:"config"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}