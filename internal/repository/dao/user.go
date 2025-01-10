package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicate = errors.New("邮箱冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDao interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	Insert(ctx context.Context, user User) error
}

type GORMUserDao struct {
	db *gorm.DB
}

func NewGORMUserDao(db *gorm.DB) UserDao {
	return &GORMUserDao{db: db}
}

func (dao *GORMUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (dao *GORMUserDao) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

func (dao *GORMUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}

func (dao *GORMUserDao) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now
	err := dao.db.Create(&user).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicate
		}
	}
	return err
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 唯一索引允许有多个空值
	// 但是不能有多个 ""
	Email    sql.NullString `gorm:"unique"`
	Password string
	Nickname string
	Phone    sql.NullString `gorm:"unique"`
	// 毫秒数而不是时间防止出现时区问题
	Ctime int64
	Utime int64
}
