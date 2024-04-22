package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
  "gorm.io/gorm"

	"cems-dis/model"
	rs "cems-dis/server/response"
)


func (s ApiService) DasLogin(c *gin.Context) rs.Response {
	login := &model.DeviceLogin{}
	err := c.BindJSON(login)
	if err != nil {
		return rs.Error(http.StatusBadRequest, "Invalid JSON body")
	}

	login.ApiKey = strings.Trim(login.ApiKey, " ")
	login.Secret = strings.Trim(login.Secret, " ")
	if len(login.ApiKey) == 0 {
		return rs.Error(http.StatusBadRequest, "API key wajib diisi")
	} else if len(login.Secret) == 0 {
		return rs.Error(http.StatusBadRequest, "Secret wajib diisi")
	}

	device := &model.Device{}
	err = s.model.DB.Where("api_key = ?", login.ApiKey).First(device).Error
	if (err != nil) && errors.Is(err, gorm.ErrRecordNotFound) {
		return rs.Error(http.StatusBadRequest, "API key tidak valid")
	} else if device.Secret != login.Secret {
		return rs.Error(http.StatusBadRequest, "Secret salah")
	}

	token, err := s.model.CreateDeviceLoginToken(device.Id)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}

	data := map[string]string{
		"login_token": token.LoginToken, 
		"refresh_token": token.RefreshToken, 
	}
	return rs.Success(data)
}

func (s ApiService) DasLoginByUid(c *gin.Context) rs.Response {
	// compatibility support only
	// for secure login, use DeviceLogin instead

	var login *struct {
		UID		string		`json:"uid"`
	}

	err := c.BindJSON(&login)
	if err != nil {
		return rs.Error(http.StatusBadRequest, "Invalid JSON body")
	}

	login.UID = strings.Trim(login.UID, " ")
	if len(login.UID) == 0 {
		return rs.Error(http.StatusBadRequest, "UID wajib diisi")
	}

	device := &model.Device{}
	err = s.model.DB.Where("uid = ?", login.UID).First(device).Error
	if (err != nil) && errors.Is(err, gorm.ErrRecordNotFound) {
		return rs.Error(http.StatusBadRequest, "UID tidak valid")
	}

	token, err := s.model.CreateDeviceLoginToken(device.Id)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}

	data := map[string]string{
		"login_token": token.LoginToken, 
		"refresh_token": token.RefreshToken, 
	}
	return rs.Success(data)
}

func (s ApiService) CemsRefreshToken(c *gin.Context) rs.Response {
	var request struct{
		RefreshToken		string	`json:"refresh_token"`
	}

	err := c.BindJSON(&request)
	if err != nil {
		return rs.Error(http.StatusBadRequest, "Invalid JSON body")
	}

	request.RefreshToken = strings.Trim(request.RefreshToken, " ")
	if len(request.RefreshToken) == 0 {
		return rs.Error(http.StatusBadRequest, "Refresh token wajib diisi")
	}

	deviceToken := &model.DeviceToken{}
	err = s.model.DB.Where("refresh_token = ?", request.RefreshToken).First(deviceToken).Error
	if (err != nil) && errors.Is(err, gorm.ErrRecordNotFound) {
		return rs.Error(http.StatusBadRequest, "Refresh token tidak valid")
	} else if deviceToken.RefreshExpiredAt.Before(time.Now()) {
		return rs.Error(http.StatusBadRequest, "Refresh token expired")
	}

	token, err := s.model.CreateDeviceLoginToken(deviceToken.DeviceId)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}

	data := map[string]string{
		"login_token": token.LoginToken, 
		"refresh_token": token.RefreshToken, 
	}
	return rs.Success(data)
}

func (s ApiService) CemsPushData(c *gin.Context) rs.Response {
	loginToken := c.GetHeader("access_token")
	loginToken = strings.Trim(loginToken, " ")
	if len(loginToken) == 0 {
		return rs.Error(http.StatusBadRequest, "Missing access_token in header")
	}

	deviceToken := &model.DeviceToken{}
	err := s.model.DB.Where("login_token = ?", loginToken).Order("id DESC").First(deviceToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs.Error(http.StatusBadRequest, "Akses token tidak valid")
		} else {
			log.Warningf("DB error: %s", err.Error())
			return rs.Error(http.StatusInternalServerError, "DB error")
		}
	} else if deviceToken.LoginExpiredAt.Before(time.Now()) {
		return rs.Error(http.StatusBadRequest, "Akses token expired")
	}

	device, err := s.model.GetDeviceById(deviceToken.DeviceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs.Error(http.StatusBadRequest, "Device untuk akses token tersebut tidak ada di database")
		} else {
			log.Warningf("DB error: %s", err.Error())
			return rs.Error(http.StatusInternalServerError, "DB error")
		}
	}

	request := &model.RawDataIn{}
	if err := c.BindJSON(request); err != nil {
		return rs.Error(http.StatusBadRequest, "Invalid JSON body")
	} else if request.Timestamp == 0 {
		return rs.Error(http.StatusBadRequest, "Timestamp tidak valid")
	}

	record := model.NewRawData(device.Id, request)
	existing := &model.RawData{}
	s.model.DB.Model(existing).
		Where("(device_id = ?) AND (timestamp = ?)", record.DeviceId, record.Timestamp).
		First(existing)
	if existing.Id == 0 {
		err = s.model.DB.Create(record).Error
	} else {
		existing.Update(record)
		record = existing
		err = s.model.DB.Save(record).Error
	}

	if err != nil {
		log.Warningf(fmt.Sprintf("DB error: %s", errors.WithStack(err)))
		return rs.Error(http.StatusInternalServerError, "DB error")
	}

	return rs.Success(record.Out())
}

func (s ApiService) ListRawData(c *gin.Context) rs.Response {
	// todo: paging

	var records []model.RawData
	if err := s.model.DB.Find(&records).Error; err != nil {
		log.Warningf("DB error: %s", err.Error())
	}
	var list []*model.RawDataOut
	for _, r := range records {
		list = append(list, r.Out())
	}
	return rs.Success(list)
}

func (s ApiService) GetRawDataById(c *gin.Context) rs.Response {
	record := &model.RawData{}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return rs.Error(http.StatusBadRequest, "Invalid record ID")
	}
	if err := s.model.DB.Model(record).First(record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs.Error(http.StatusNotFound, "Record ID not found")
		} else {
			log.Warningf("DB error: %+v", errors.WithStack(err))
			return rs.Error(http.StatusInternalServerError, "DB error")
		}
	}
	return rs.Success(record.Out())
}

func (s ApiService) CemsRawData(c *gin.Context) rs.Response {
	deviceId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ts1, _ := strconv.ParseInt(c.Query("ts1"), 10, 64)
	ts2, _ := strconv.ParseInt(c.Query("ts2"), 10, 64)
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	if deviceId == 0 {
		return rs.Error(http.StatusBadRequest, "ID device tidak valid")
	}

	device := &model.Device{}
	err := s.model.DB.Model(device).Find(device).Error
	if (err != nil) && errors.Is(err, gorm.ErrRecordNotFound) {
		return rs.Error(http.StatusNotFound, "ID device tidak ada di database")
	}

	sql := s.model.DB.Model(&model.RawData{}).Order("timestamp DESC").Where("device_id = ?", deviceId)
	if ts1 > 0 {
		sql = sql.Where("(timestamp >= ?)", ts1)
	}
	if ts2 > 0 {
		sql = sql.Where("(timestamp <= ?)", ts2)
	}

	paging, err := rs.NewPagingFormatter(sql, size, page)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, "DB error")
	}
	var records []*model.RawData
	if err = paging.Sql().Find(&records).Error; err != nil {
		log.Warningf("DB error: %s", err.Error())
	}
	list := []*model.RawDataOut{}
	for _, r := range records {
		list = append(list, r.Out())
	}
	return rs.Success(list).UseFormatter(paging)
}

func (s ApiService) CemsEmissionData(c *gin.Context) rs.Response {
	return rs.Success(nil)
}

func (s ApiService) CemsPercentageData(c *gin.Context) rs.Response {
	return rs.Success(nil)
}
