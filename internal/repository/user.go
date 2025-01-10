package repository

import (
	"ShareSphere/V0/internal/domain"
	"ShareSphere/V0/internal/repository/dao"
	"context"
	"database/sql"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	Create(ctx context.Context, user domain.User) error
}

type CachedUserRepository struct {
	dao dao.UserDao
}

func NewUserRepository(dao dao.UserDao) UserRepository {
	return &CachedUserRepository{dao: dao}
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// SELECT * FROM `users` WHERE `email`=?
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entity2Domain(u), nil
}
func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entity2Domain(u), nil
}
func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return r.entity2Domain(u), nil
}
func (r *CachedUserRepository) Create(ctx context.Context, user domain.User) error {
	return r.dao.Insert(ctx, r.domain2Entity(user))
}

func (r *CachedUserRepository) entity2Domain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

func (r *CachedUserRepository) domain2Entity(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    sql.NullString{String: u.Email, Valid: u.Email != ""},
		Password: u.Password,
		Phone:    sql.NullString{String: u.Phone, Valid: u.Phone != ""},
		Ctime:    u.Ctime.UnixMilli(),
	}
}
