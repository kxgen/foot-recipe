package user

type UpdatePasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
