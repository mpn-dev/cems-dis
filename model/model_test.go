package model

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ModelTestSuite struct {
	suite.Suite
	model							*Model
	sqlMock						sqlmock.Sqlmock
}


func (s *ModelTestSuite) SetupSuite() {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	s.model = &Model{DB: db}
	s.sqlMock = mock
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}
