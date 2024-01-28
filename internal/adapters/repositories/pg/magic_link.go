package repositories

import (
	"database/sql"
	"hyperzoop/internal/core/entities"
)

type MagicLinkPostgresRepository struct {
	db *sql.DB
}

func NewMagicLinkPostgresRepository(db *sql.DB) *MagicLinkPostgresRepository {
	return &MagicLinkPostgresRepository{db: db}
}

func (r *MagicLinkPostgresRepository) Create(link *entities.MagicLink) error {
	_, err := r.db.Exec("INSERT INTO magic_links (code, user_id, cookie, valid_until, used) VALUES ($1, $2, $3, $4, $5)", link.Code, link.UserId, link.Cookie, link.ValidUntil, link.Used)
	return err
}

func (r *MagicLinkPostgresRepository) FindValidByCode(code, cookie string) (*entities.MagicLink, error) {
	row := r.db.QueryRow("SELECT code, user_id, cookie, valid_until, used FROM magic_links WHERE code = $1 AND cookie = $2 AND valid_until > NOW() AND used = false", code, cookie)
	return convertRowToMagicLink(row)
}

func (r *MagicLinkPostgresRepository) Invalidate(code string) error {
	_, err := r.db.Exec("UPDATE magic_links SET used = true WHERE code = $1", code)
	return err
}

func (r *MagicLinkPostgresRepository) Update(link *entities.MagicLink) error {
	_, err := r.db.Exec("UPDATE magic_links SET used = $1 WHERE code = $2", link.Used, link.Code)
	return err
}
func convertRowToMagicLink(row *sql.Row) (*entities.MagicLink, error) {
	var link entities.MagicLink
	err := row.Scan(&link.Code, &link.UserId, &link.Cookie, &link.ValidUntil, &link.Used)
	return &link, err
}
