package ports

import "hyperzoop/internal/core/entities"

type MagicLinkRepository interface {
	Create(link *entities.MagicLink) error
	FindValidByCode(code, cookie string) (*entities.MagicLink, error)
	Invalidate(code string) error
	Update(link *entities.MagicLink) error
}
