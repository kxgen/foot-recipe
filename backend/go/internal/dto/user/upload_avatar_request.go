package user

type UploadAvatarRequest struct {
	UserId      int    `json:"userId"`
	Slug        string `json:"slug"`
	ImageName   string `json:"imageName"`
	Base64Image string `json:"base64Image"`
}
