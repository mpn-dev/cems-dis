package api

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"cems-dis/model"
)


func uidCheckSql() string {
	return `SELECT \* FROM "devices" WHERE \(uid = .+\) AND \(uid <> .+\) ORDER BY "devices"."uid" LIMIT .+`
}

func keyCheckSql() string {
	return `SELECT \* FROM "devices" WHERE \(api_key = .+\) AND \(uid <> .+\) ORDER BY "devices"."uid" LIMIT .+`
}

func devCheckResult() *sqlmock.Rows {
	// sqlmock won't work correctly if same rows reused 
	// in different test cases within the same function
	return sqlmock.NewRows([]string{"uid"}).AddRow("1001")
}

func devWriteResult() *sqlmock.Rows {
	// sqlmock won't work correctly if same rows reused 
	// in different test cases within the same function
	lat := 7.2
	lng := 101.2
	rows := sqlmock.NewRows([]string{"uid", "name", "latitude", "longitude", "api_key", "secret", "enabled", "created_at", "updated_at"})
	rows = rows.AddRow("1001", "Device #1", &lat, &lng, "key1", "sec1", true, time.Time{}, time.Time{})
	return rows
}

func examplePayload() string {
	return `{"uid": "1001", "name": "Device #1", "latitude": 7.2, "longitude": 101.2, "apikey": "key1", "secret": "sec1", "enabled": true}`
}

func exampleDevice() *model.Device {
	lat := 7.2
	lng := 101.2
	return &model.Device{
		UID:      	"1001", 
		Name:       "Device #1", 
		Latitude:   &lat, 
		Longitude:  &lng, 
		ApiKey:     "key1", 
		Secret:     "sec1", 
		Enabled:    true, 
		CreatedAt:	time.Time{}, 
		UpdatedAt:	time.Time{}, 
	}
}

func (s *ApiHandlerTestSuite) TestListDevices() {
	cntDevs := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(32)
	devices := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Device #1").
		AddRow(2, "Device #2")

	cases := []struct{
		name		string
		query		gin.H
		dbGet		bool
		paging	bool
		where		string
		limit		string
		cntArgs	[]driver.Value
		getArgs	[]driver.Value
		code		int
		mssg		interface{}
	}{
		{
			name:			"DB error", 
			query:		gin.H{}, 
			dbGet:		true, 
			paging:		false, 
			where:		"", 
			limit:		"", 
			cntArgs:	[]driver.Value{}, 
			getArgs:	[]driver.Value{}, 
			code:			http.StatusInternalServerError, 
			mssg:			"DB error", 
		}, 
		{
			name:			"success without paging, disabled only", 
			query:		gin.H{"disabled": ""}, 
			dbGet:		true, 
			paging:		false, 
			where:		` WHERE \(enabled = .+\)`, 
			limit:		"", 
			cntArgs:	[]driver.Value{}, 
			getArgs:	[]driver.Value{false}, 
			code:			http.StatusOK, 
			mssg:			nil, 
		}, 
		{
			name:			"success with paging, without page adjustment", 
			query:		gin.H{"page": "4", "size": "10"}, 
			dbGet:		true, 
			paging:		true, 
			where:		"", 
			limit:		" LIMIT .+ OFFSET .+", 
			cntArgs:	[]driver.Value{}, 
			getArgs:	[]driver.Value{10, 30}, 
			code:			http.StatusOK, 
			mssg:			nil, 
		}, 
		{
			name:			"success with paging, with page adjustment", 
			query:		gin.H{"page": "8", "size": "10"}, 
			dbGet:		true, 
			paging:		true, 
			where:		"", 
			limit:		" LIMIT .+", 
			cntArgs:	[]driver.Value{}, 
			getArgs:	[]driver.Value{10}, 
			code:			http.StatusOK, 
			mssg:			nil, 
		}, 
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			if len(c.query) > 0 {
				qry := url.Values{}
				for k, v := range c.query {
					d, _ := v.(string)
					qry.Add(k, d)
				}
				s.ctx.Request.URL.RawQuery = qry.Encode()
			}

			if c.dbGet {
				if c.paging {
					exp1 := s.sqlMock.ExpectQuery(fmt.Sprintf(`SELECT count\(\*\) FROM "devices"%s`, c.where))
					if len(c.cntArgs) > 0 {
						exp1 = exp1.WithArgs(c.cntArgs)
					}
					exp1.WillReturnRows(cntDevs)
				}

				exp2 := s.sqlMock.ExpectQuery(fmt.Sprintf(`SELECT \* FROM "devices"%s ORDER BY name%s`, c.where, c.limit))
				if len(c.getArgs) > 0 {
					exp2 = exp2.WithArgs(c.getArgs...)
				}
				if c.mssg == nil {
					exp2.WillReturnRows(devices)
				} else {
					exp2.WillReturnError(errors.New(c.mssg.(string)))
				}
			}

			resp := s.api.ListDevices(s.ctx)
			s.NoError(s.sqlMock.ExpectationsWereMet())
			s.Equal(c.code, resp.Code)
			s.Equal(c.mssg, resp.Error)
		})
	}
}

func (s *ApiHandlerTestSuite) TestGetDevice() {
	getSql := `SELECT \* FROM "devices" WHERE uid = .+ ORDER BY "devices"\."uid" LIMIT .+`
	device := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Device #1")

	cases := []struct{
		name 		string
		uid			string
		dbGet		bool
		dbErr		error
		code		int
		mssg		interface{}
	}{
		{
			name: 	"uid empty", 
			uid:		"", 
			dbGet:	false, 
			dbErr:	nil, 
			code:		http.StatusBadRequest, 
			mssg:		"UID tidak boleh kosong", 
		}, 
		{
			name: 	"uid contain only spaces", 
			uid:		"  ", 
			dbGet:	false, 
			dbErr:	nil, 
			code:		http.StatusBadRequest, 
			mssg:		"UID tidak boleh kosong", 
		}, 
		{
			name: 	"uid not in database", 
			uid:		"1024", 
			dbGet:	true, 
			dbErr:	gorm.ErrRecordNotFound, 
			code:		http.StatusNotFound, 
			mssg:		"UID tidak valid", 
		}, 
		{
			name: 	"database error", 
			uid:		"1024", 
			dbGet:	true, 
			dbErr:	errors.New("DB error"), 
			code:		http.StatusNotFound, 
			mssg:		"UID tidak valid", 
		}, 
		{
			name: 	"no error", 
			uid:		"1024", 
			dbGet:	true, 
			dbErr:	nil, 
			code:		http.StatusOK, 
			mssg:		nil, 
		}, 
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			s.ctx.Params = gin.Params{{Key: "uid", Value: c.uid}}

			if c.dbGet {
				if(c.dbErr != nil) {
					s.sqlMock.ExpectQuery(getSql).WithArgs(c.uid, 1).WillReturnError(c.dbErr)
				} else {
					s.sqlMock.ExpectQuery(getSql).WithArgs(c.uid, 1).WillReturnRows(device)
				}
			}

			resp := s.api.GetDevice(s.ctx)
			s.NoError(s.sqlMock.ExpectationsWereMet())
			s.Equal(c.code, resp.Code)
			s.Equal(c.mssg, resp.Error)
		})
	}
}

func (s *ApiHandlerTestSuite) TestInsertDevice() {
	// var nullDevice *model.Device

	cases := []struct{
		name		string
		body		string
		wSql		bool
		wErr		error
		code		int
		mssg		interface{}
		data		interface{}
	}{
		{
			name:		"extract device data failed", 
			body:		"", 
			wSql:		false, 
			wErr:		nil, 
			code:		http.StatusBadRequest, 
			mssg:		"Invalid JSON body", 
		}, 
		{
			name:		"database error", 
			body:		examplePayload(), 
			wSql:		true, 
			wErr:		errors.New("dummy database error"), 
			code:		http.StatusInternalServerError, 
			mssg:		"Unknown error", 
		}, 
		{
			name:		"success", 
			body:		examplePayload(), 
			wSql:		true, 
			wErr:		nil, 
			code:		http.StatusOK, 
			mssg:		nil, 
		}, 
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			s.ctx.Request.Body = io.NopCloser(strings.NewReader(c.body))
			if c.wSql {
				s.sqlMock.ExpectQuery(uidCheckSql()).WithArgs("1001", "", 1).WillReturnError(gorm.ErrRecordNotFound)
				s.sqlMock.ExpectQuery(keyCheckSql()).WithArgs("key1", "", 1).WillReturnError(gorm.ErrRecordNotFound)

				s.sqlMock.ExpectBegin()
				mock := s.sqlMock.ExpectExec(`INSERT INTO "devices" \(.+\) VALUES \(.+\)`)
				if c.wErr == nil {
					mock.WillReturnResult(sqlmock.NewResult(1, 1))
					s.sqlMock.ExpectCommit()
				} else {
					mock.WillReturnError(c.wErr)
					s.sqlMock.ExpectRollback()
				}
			}

			resp := s.api.InsertDevice(s.ctx)
			s.NoError(s.sqlMock.ExpectationsWereMet())
			s.Equal(c.code, resp.Code)
			s.Equal(c.mssg, resp.Error)
		})
	}
}

func (s *ApiHandlerTestSuite) TestUpdateDevice() {
}

func (s *ApiHandlerTestSuite) TestDeleteDevice() {
}

func (s *ApiHandlerTestSuite) TestGenerateDeviceSecret() {
}

func (s *ApiHandlerTestSuite) TestExtractDeviceData() {
	cases := []struct{
		name				string
		uidParam		string
		body				string
		uidToCheck	interface{}
		keyToCheck	interface{}
		uidTaken		bool
		keyTaken		bool
		device			*model.Device
		error				error
	}{
		{
			name:				"json parsing failed", 
			uidParam:		"", 
			body:				`ABCD-PQRS`, 
			uidToCheck:	nil, 
			keyToCheck:	nil, 
			uidTaken:		false, 
			keyTaken:		false, 
			device:			nil, 
			error:			errors.New("Invalid JSON body"), 
		}, 
		{
			name:				"uid empty", 
			uidParam:		"", 
			body:				`{"uid": "  ", "name": "Device #1", "latitude": 7.2, "longitude": 101.2, "apikey": "key1", "secret": "sec1", "enabled": true}`, 
			uidToCheck:	nil, 
			keyToCheck:	nil, 
			uidTaken:		false, 
			keyTaken:		false, 
			device:			nil, 
			error:			errors.New("UID wajib diisi"), 
		}, 
		{
			name:				"uid already taken", 
			uidParam:		"", 
			body:				examplePayload(), 
			uidToCheck:	"1001", 
			keyToCheck:	nil, 
			uidTaken:		true, 
			keyTaken:		false, 
			device:			nil, 
			error:			errors.New("UID '1001' sudah digunakan"), 
		}, 
		{
			name:				"name empty", 
			uidParam:		"", 
			body:				`{"uid": "1001", "name": "  ", "latitude": 7.2, "longitude": 101.2, "apikey": "key1", "secret": "sec1", "enabled": true}`, 
			uidToCheck:	"1001", 
			keyToCheck:	nil, 
			uidTaken:		false, 
			keyTaken:		false, 
			device:			nil, 
			error:			errors.New("Nama device wajib diisi"), 
		}, 
		{
			name:				"apikey empty", 
			uidParam:		"", 
			body:				`{"uid": "1001", "name": "Device #1", "latitude": 7.2, "longitude": 101.2, "apikey": "  ", "secret": "sec1", "enabled": true}`, 
			uidToCheck:	"1001", 
			keyToCheck:	nil, 
			uidTaken:		false, 
			keyTaken:		false, 
			device:			nil, 
			error:			errors.New("API key wajib diisi"), 
		}, 
		{
			name:				"apikey already taken", 
			uidParam:		"", 
			body:				examplePayload(), 
			uidToCheck:	"1001", 
			keyToCheck:	"key1", 
			uidTaken:		false, 
			keyTaken:		true, 
			device:			nil, 
			error:			errors.New("API key 'key1' sudah digunakan"), 
		}, 
		{
			name:				"secret empty", 
			uidParam:		"", 
			body:				`{"uid": "1001", "name": "Device #1", "latitude": 7.2, "longitude": 101.2, "apikey": "key1", "secret": "   ", "enabled": true}`, 
			uidToCheck:	"1001", 
			keyToCheck:	"key1", 
			uidTaken:		false, 
			keyTaken:		false, 
			device:			nil, 
			error:			errors.New("Secret wajib diisi"), 
		}, 
		{
			name:				"success without uid", 
			uidParam:		"", 
			body:				examplePayload(), 
			uidToCheck:	"1001", 
			keyToCheck:	"key1", 
			uidTaken:		false, 
			keyTaken:		false, 
			device:			exampleDevice(), 
			error:			nil, 
		}, 
		{
			name:				"success and string fields automatically trimmed", 
			uidParam:		"", 
			body:				`{"uid": "  1001  ", "name": "  Device #1  ", "latitude": 7.2, "longitude": 101.2, "apikey": "  key1  ", "secret": "  sec1  ", "enabled": true}`, 
			uidToCheck:	"1001", 
			keyToCheck:	"key1", 
			uidTaken:		false, 
			keyTaken:		false, 
			device:			exampleDevice(), 
			error:			nil, 
		}, 
		{
			name:				"success with uid", 
			uidParam:		"1001", 
			body:				examplePayload(), 
			uidToCheck:	"1001", 
			keyToCheck:	"key1", 
			uidTaken:		false, 
			keyTaken:		false, 
			device:			exampleDevice(), 
			error:			nil, 
		}, 
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			if len(c.uidParam) > 0 {
				s.ctx.Params = gin.Params{{Key: "uid", Value: c.uidParam}}
			}

			s.ctx.Request.Body = io.NopCloser(strings.NewReader(c.body))

			if c.uidToCheck != nil {
				mock := s.sqlMock.ExpectQuery(uidCheckSql()).WithArgs(c.uidToCheck.(string), c.uidParam, 1)
				if c.uidTaken {
					mock.WillReturnRows(devCheckResult())
				} else {
					mock.WillReturnError(gorm.ErrRecordNotFound)
				}
			}

			if c.keyToCheck != nil {
				mock := s.sqlMock.ExpectQuery(keyCheckSql()).WithArgs(c.keyToCheck.(string), c.uidParam, 1)
				if c.keyTaken {
					mock.WillReturnRows(devCheckResult())
				} else {
					mock.WillReturnError(gorm.ErrRecordNotFound)
				}
			}

			dev, err := s.api.extractDeviceData(s.ctx)
			s.NoError(s.sqlMock.ExpectationsWereMet())
			s.Equal(c.device, dev)
			s.Equal(c.error, err)
		})
	}
}
