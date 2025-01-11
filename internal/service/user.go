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
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
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

// 通过Unique索引保证 check and set的原子性 如果加锁，会导致性能问题，通过都是从快路径出
func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	// 快路径
	if err != repository.ErrUserNotFound {
		// 找到用户或者存在错误
		return u, err
	}
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 可能会遇到主从延迟
	return svc.repo.FindByPhone(ctx, phone)
}
