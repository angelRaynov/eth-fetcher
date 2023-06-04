package model

// Credentials represents the request body for authentication
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the response format for authentication
type AuthResponse struct {
	Token string `json:"token"`
}

