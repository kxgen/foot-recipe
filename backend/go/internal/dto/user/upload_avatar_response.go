package user

type UploadAvatarResponse struct {
	Success   bool   `json:"success"`
	// After Upload, Return and Show Avatar
	AvatarUrl string `json:"avatarUrl"`
}
