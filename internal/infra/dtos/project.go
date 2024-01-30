package dtos

type CreateProjectDTOInput struct {
	Name        string  `json:"name"`
	Status      *string `json:"status"`
	Description *string `json:"description"`
}

type UpdateProjectDTOInput struct {
	ID          int     `json:"id"`
	Name        *string `json:"name"`
	Status      *string `json:"status"`
	Description *string `json:"description"`
}
