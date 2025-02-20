package db

import (
	"backend/blog/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *configs.Configs,) *Db {
	db, err := gorm.Open(postgres.Open(conf.DSN), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	return &Db{
		db,
	}
}