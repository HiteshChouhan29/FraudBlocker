package models

// PhoneNumber represents a blocked phone number
type PhoneNumber struct {
	ID          int    `json:"id"`
	Number      string `json:"number"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}
