package addon

import "time"

type CreateAddonRequest struct {
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Tier    string `json:"tier" binding:"required"`
	Storage string `json:"storage"`
}

type UpdateAddonRequest struct {
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Tier    string `json:"tier" binding:"required"`
	Storage string `json:"storage"`
}

type AddonResponse struct {
	ID        uint      `json:"id"`
	ProjectID uint      `json:"project_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Tier      string    `json:"tier"`
	Storage   string    `json:"storage"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}