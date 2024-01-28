package repositories

import (
	"database/sql"
	"hyperzoop/internal/core/entities"

	"github.com/lib/pq"
)

type UserPostgresRepository struct {
	db *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

func (r *UserPostgresRepository) Create(user *entities.User) (*entities.User, error) {
	row := r.db.QueryRow("INSERT INTO users (username, email, avatar, blocked) VALUES ($1, $2, $3, $4) RETURNING id, username, email, avatar, blocked, created_at, updated_at", user.Username, user.Email, user.Avatar, user.Blocked)
	return convertRowToUser(row)
}

func (r *UserPostgresRepository) FindByEmail(email string) (*entities.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, avatar, blocked, created_at, updated_at FROM users WHERE email = $1 LIMIT 1", email)
	return convertRowToUser(row)
}

func (r *UserPostgresRepository) FindById(id string) (*entities.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, avatar, blocked, created_at, updated_at FROM users WHERE id = $1 LIMIT 1", id)
	return convertRowToUser(row)
}

func (r *UserPostgresRepository) FindUserBySliceIds(ids []string) ([]*entities.User, error) {
	rows, err := r.db.Query("SELECT id, username, email, avatar, blocked, created_at, updated_at FROM users WHERE id = ANY($1)", pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return convertRowToUserSlice(rows)
}

func convertRowToUser(row *sql.Row) (*entities.User, error) {
	var user entities.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Avatar, &user.Blocked, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func convertRowToUserSlice(rows *sql.Rows) (users []*entities.User, err error) {
	for rows.Next() {
		user := &entities.User{}
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.Avatar, &user.Blocked, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	return
}
