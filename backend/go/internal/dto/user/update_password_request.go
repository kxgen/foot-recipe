package user

type UpdatePasswordRequest struct {
	NewPassword string `json:"newPassword"`
}
