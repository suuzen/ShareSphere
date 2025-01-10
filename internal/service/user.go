package service

import (
	"ShareSphere/V0/internal/domain"
	"ShareSphere/V0/internal/repository"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService interface {
	Login(ctx context.Context, email string, password string) (domain.User, error)
	SignUp(ctx context.Context, user domain.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	// 查找存在用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil

}
func (svc *userService) SignUp(ctx context.Context, user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	return svc.repo.Create(ctx, user)
}
