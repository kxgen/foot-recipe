package auth

type AuthResponse struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}
