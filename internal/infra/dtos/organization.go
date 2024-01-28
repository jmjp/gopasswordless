package dtos

type UpdateOrganizationInputDTO struct {
	ID      string  `json:"id"`
	Name    *string `json:"name"`
	Phone   *string `json:"phone"`
	Address *string `json:"address"`
	Logo    *string `json:"logo"`
}
