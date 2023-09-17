package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func NewMySQLDBGormMock() (sqlmock.Sqlmock, *gorm.DB) {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal("error in new connection", zap.Error(err))
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}))
	if err != nil {
		log.Fatal("error in open connection", zap.Error(err))
	}

	return sqlMock, db
}
