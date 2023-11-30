package dto

type (
	UploadProfilePictureRequestBody struct {
		ImageURL string `json:"image_url" validate:"required,url"`
	}
	UploadProfilePicturePayload struct {
		ImageURL string
		UserID   int64
	}
)
