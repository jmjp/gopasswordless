package dtos

type CreateFileInputDTO struct {
	Name         string `json:"name"`
	UserId       string `json:"user_id"`
	Organization string `json:"organization"`
	Size         int64  `json:"size"`
	Extension    string `json:"type"`
}
