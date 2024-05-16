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

	token, err := s.model.CreateDeviceLoginToken(device.UID)
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
		return rs.Error(http.StatusBadRequest, "UID tidak ada di database")
	}

	token, err := s.model.CreateDeviceLoginToken(device.UID)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}

	data := map[string]string{
		"login_token": token.LoginToken, 
		"refresh_token": token.RefreshToken, 
	}
	return rs.Success(data)
}

func (s ApiService) DasRefreshToken(c *gin.Context) rs.Response {
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

	token, err := s.model.CreateDeviceLoginToken(deviceToken.DEV)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}

	data := map[string]string{
		"login_token": token.LoginToken, 
		"refresh_token": token.RefreshToken, 
	}
	return rs.Success(data)
}

func (s ApiService) DasReceiveData(c *gin.Context) rs.Response {
	loginToken := c.GetHeader("access_token")
	loginToken = strings.Trim(loginToken, " ")
	if len(loginToken) == 0 {
		return rs.Error(http.StatusBadRequest, "Missing access_token in header")
	}

	deviceToken := &model.DeviceToken{}
	err := s.model.DB.Where("refresh_token = ?", loginToken).Order("id DESC").First(deviceToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs.Error(http.StatusBadRequest, "Akses token tidak valid")
		} else {
			log.Warningf("DB error: %s", err.Error())
			return rs.Error(http.StatusInternalServerError, "DB error")
		}
	} else if deviceToken.RefreshExpiredAt.Before(time.Now()) {
		return rs.Error(http.StatusBadRequest, "Akses token expired")
	}

	device, err := s.model.GetDeviceByUid(deviceToken.DEV)
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

	record := model.NewRawData(device.UID, request)
	existing := &model.RawData{}
	s.model.DB.Model(existing).
		Where("(uid = ?) AND (timestamp = ?)", record.DEV, record.Timestamp).
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

	s.queueDataTransmission(record.Id);

	return rs.Success(record.Out())
}

func (s ApiService) DasRelayData(d *model.RawData) error {
	return nil
}

func (s ApiService) GetRawDataById(c *gin.Context) rs.Response {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	record, err := s.model.GetRawDataById(id)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, "Invalid record ID")
	}
	if record == nil {
		return rs.Error(http.StatusNotFound, "Record ID not found")
	}
	return rs.Success(record.Out())
}

func (s ApiService) ListRawData(c *gin.Context) rs.Response {
	uid := strings.Trim(c.Param("uid"), " ")
	ts1, _ := strconv.ParseInt(c.Query("ts1"), 10, 64)
	ts2, _ := strconv.ParseInt(c.Query("ts2"), 10, 64)
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	if err := s.validateUidString(uid); err != nil {
		return rs.Error(http.StatusBadRequest, err.Error())
	}

	device := &model.Device{UID: uid}
	err := s.model.DB.Model(&model.Device{}).Find(device).Error
	if (err != nil) && errors.Is(err, gorm.ErrRecordNotFound) {
		return rs.Error(http.StatusNotFound, "ID device tidak ada di database")
	}

	sql := s.model.DB.Model(&model.RawData{}).Order("timestamp DESC").Where("uid = ?", uid)
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

func (s ApiService) ListEmissionData(c *gin.Context) rs.Response {
	return rs.Success(nil)
}

func (s ApiService) ListPercentageData(c *gin.Context) rs.Response {
	return rs.Success(nil)
}

func (s ApiService) queueDataTransmission(rawDataId uint64) {
	var stations []*model.RelayStation
	err := s.model.DB.Model(&model.RelayStation{}).Where("enabled = ?", true).Find(&stations).Error
	if err != nil {
		log.Warningf("DB error: %s", err.Error())
		return
	}
	for _, sta := range stations {
		trx := model.Transmission{
			RawDataId:			rawDataId, 
			RelayStationId:	sta.Id, 
			Code:						0, 
			Error:					"", 
			Status:					"Pending", 
		}
		s.model.DB.Save(&trx)
	}
}
