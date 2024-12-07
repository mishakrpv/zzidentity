package state

import (
	"github.com/jmoiron/sqlx"
	"github.com/zzidentity/zzidentity/pkg/state/model"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	tx := r.db.MustBegin()
	tx.MustExec("ISERT INTO users (user_id, email, password_hash) VALUES ($1, $2, $3)",
		user.ID, user.Email, user.PwdHash)
	return tx.Commit()
}
