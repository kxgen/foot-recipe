package recipe

type UploadRecipeImageRequest struct {
	ImageName   string `json:"imageName"`
	Base64Image string `json:"base64Image"`
	RecipeId    int    `json:"recipeId"`
	TempFolder  string `json:"tempFolder"`
	Filename    string `json:"filename"`
}

type UploadRecipeImageResponse struct {
	Success  bool   `json:"success"`
	ImageUrl string `json:"imageUrl"`
}

type ConfirmRecipeImagesRequest struct {
	RecipeId   int    `json:"recipeId"`
	TempFolder string `json:"tempFolder"`
}

type ConfirmRecipeImagesResponse struct {
	Success bool `json:"success"`
}

type DeleteRecipeFolderRequest struct {
	RecipeId int `json:"recipeId"`
}

type DeleteRecipeFolderResponse struct {
	Success bool `json:"success"`
}
