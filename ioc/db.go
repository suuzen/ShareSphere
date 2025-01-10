package ioc

import (
	"ShareSphere/V0/internal/repository/dao"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:13578101231a@tcp(localhost:3306)/share_sphere"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db

}
