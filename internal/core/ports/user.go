package ports

import "hyperzoop/internal/core/entities"

type UserRepository interface {
	Create(user *entities.User) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	FindById(id string) (*entities.User, error)
	FindUserBySliceIds(ids []string) ([]*entities.User, error)
}
