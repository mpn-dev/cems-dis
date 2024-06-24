package api

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"cems-dis/model"
)

type ApiHandlerTestSuite struct {
	suite.Suite
	sqlMock				sqlmock.Sqlmock
	model					*model.Model
	api						ApiService
	rec						*httptest.ResponseRecorder
	ctx						*gin.Context
}


func (s *ApiHandlerTestSuite) setupSqlMock() {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})

	gin.SetMode(gin.TestMode)
	s.sqlMock = mock
	s.model = &model.Model{DB: db}
	s.api = ApiService{model: s.model}
	s.rec = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.rec)
	s.ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:		&url.URL{}, 
	}
}

func (s *ApiHandlerTestSuite) SetupSuite() {
}

func (s *ApiHandlerTestSuite) SetupTest() {
}

func (s *ApiHandlerTestSuite) SetupSubTest() {
	s.setupSqlMock()
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(ApiHandlerTestSuite))
}
