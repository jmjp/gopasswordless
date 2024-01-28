package repositories

import (
	"database/sql"
	"hyperzoop/internal/core/entities"
)

type SessionPostgresRepository struct {
	db *sql.DB
}

func NewSessionPostgresRepository(db *sql.DB) *SessionPostgresRepository {
	return &SessionPostgresRepository{db: db}
}

func (r *SessionPostgresRepository) Create(session *entities.Session) (*entities.Session, error) {
	row := r.db.QueryRow("INSERT INTO sessions (user_id, valid_until) VALUES ($1, $2) RETURNING id, user_id, valid_until, user_agent, ip, latitude, longitude, city, region, country, isp, created_at, updated_at", session.UserId, session.ValidUntil)
	return convertRowToSession(row)
}

func (r *SessionPostgresRepository) All(userId string) ([]*entities.Session, error) {
	rows, err := r.db.Query("SELECT id, user_id, valid_until, user_agent, ip, latitude, longitude, city, region, country, isp, created_at, updated_at FROM sessions WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*entities.Session
	for rows.Next() {
		session, err := convertRowToSessionSlice(rows)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (r *SessionPostgresRepository) One(id string) (*entities.Session, error) {
	row := r.db.QueryRow("SELECT id, user_id, valid_until, user_agent, ip, latitude, longitude, city, region, country, isp, created_at, updated_at FROM sessions WHERE id = $1 LIMIT 1", id)
	return convertRowToSession(row)
}

func (r *SessionPostgresRepository) UpdateGeoLocation(session *entities.Session) error {
	_, err := r.db.Exec("UPDATE sessions set ip = $1, latitude = $2, longitude = $3, city = $4, region = $5, country = $6, isp = $7 WHERE id = $8", session.Ip, session.Latitude, session.Longitude, session.City, session.Region, session.Country, session.OrganizationName, session.Id)
	return err
}

func (r *SessionPostgresRepository) Update(session *entities.Session) error {
	_, err := r.db.Exec("UPDATE sessions set valid_until = $1, user_agent = $2 WHERE id = $3", session.ValidUntil, session.UserAgent, session.Id)
	return err
}

func (r *SessionPostgresRepository) Disconnect(sessionId string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE id = $1", sessionId)
	return err
}

func convertRowToSession(row *sql.Row) (*entities.Session, error) {
	var s entities.Session
	err := row.Scan(&s.Id, &s.UserId, &s.ValidUntil, &s.UserAgent, &s.Ip, &s.Latitude, &s.Longitude, &s.City, &s.Region, &s.Country, &s.OrganizationName, &s.CreatedAt, &s.UpdatedAt)
	return &s, err
}

func convertRowToSessionSlice(rows *sql.Rows) (*entities.Session, error) {
	var s entities.Session
	err := rows.Scan(&s.Id, &s.UserId, &s.ValidUntil, &s.UserAgent, &s.Ip, &s.Latitude, &s.Longitude, &s.City, &s.Region, &s.Country, &s.OrganizationName, &s.CreatedAt, &s.UpdatedAt)
	return &s, err
}
