package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
  "gorm.io/gorm"

	"cems-dis/model"
	rs "cems-dis/server/response"
	"cems-dis/utils"
)


func (s ApiService) DasReceiveData(c *gin.Context) rs.Response {
	handler := func(svc ApiService, pr *model.PushRequest, devToken *model.DeviceToken) rs.Response {
		setErr := func(r *model.PushRequest, code int, msg string) rs.Response {
			r.Status = "Error"
			r.Info = msg
			return rs.Error(code, msg)
		}

		bodyRaw, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return setErr(pr, http.StatusBadRequest, "Error reading request body")
		}

		pr.Request = string(bodyRaw)
		device, err := s.model.GetDeviceByUid(devToken.DEV)
		if err != nil {
			return setErr(pr, http.StatusInternalServerError, err.Error())
		} else if device == nil {
			return setErr(pr, http.StatusBadRequest, "Device untuk bearer token tersebut tidak ada di database")
		}

		request := &model.RawDataIn{}
		if err := json.Unmarshal(bodyRaw, request); err != nil {
			return setErr(pr, http.StatusBadRequest, "Invalid JSON body")
		} else if request.Timestamp == 0 {
			return setErr(pr, http.StatusBadRequest, "Timestamp tidak valid")
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
			return setErr(pr, http.StatusInternalServerError, "DB error")
		}

		s.queueDataTransmission(record.Id);

		pr.Status = "Success"
		return rs.New(201, nil, record.Out()).UseDasFormatter()
	}

	token := utils.ParseBearerToken(c.GetHeader("Authorization"))
	deviceToken, _ := s.model.GetDeviceToken(token)
	pushRequest := &model.PushRequest{
		IpAddr: 		c.ClientIP(), 
		UserAgent:	c.Request.UserAgent(), 
		DEV: 				deviceToken.DEV, 
		Request:		"", 
		Status:			"", 
		Info:				"", 
	}

	if err := s.model.DB.Save(pushRequest).Error; err != nil {
		log.Warningf("DB error: %s", err.Error())
		return rs.Error(http.StatusInternalServerError, "DB error")
	}

	res := handler(s, pushRequest, deviceToken)
	if err := s.model.DB.Save(pushRequest).Error; err != nil {
		log.Warningf("DB error: %s", err.Error())
	}
	return res
}

func (s ApiService) DasLogin(c *gin.Context) rs.Response {
	//  this function is not used

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
		"access_token": token.LoginToken, 
	}
	return rs.Success(data)
}

func (s ApiService) DasRefreshToken(c *gin.Context) rs.Response {
	// this function is not used

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

func (s ApiService) queueDataTransmission(rawDataId uint64) {
	var stations []*model.RelayStation
	err := s.model.DB.Model(&model.RelayStation{}).Where("enabled = ?", true).Find(&stations).Error
	if err != nil {
		log.Warningf("DB error: %s", err.Error())
		return
	}
	for _, sta := range stations {
		trx := model.Transmission{
			RawDataId:	rawDataId, 
			Station:		sta.Name, 
			Protocol:		sta.Protocol, 
			BaseURL: 		sta.BaseURL, 
			Username:		sta.Username, 
			Password: 	sta.Password, 
			Code:				0, 
			Note:				"", 
			Status:			"Pending", 
		}
		s.model.DB.Save(&trx)
	}
}
