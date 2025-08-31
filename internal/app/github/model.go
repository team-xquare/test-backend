package github

import "time"

type Installation struct {
	ID             uint      `json:"id" db:"id"`
	InstallationID string    `json:"installation_id" db:"installation_id"`
	UserID         uint      `json:"user_id" db:"user_id"`
	AccountLogin   string    `json:"account_login" db:"account_login"`
	AccountType    string    `json:"account_type" db:"account_type"`
	Permissions    string    `json:"permissions" db:"permissions"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}