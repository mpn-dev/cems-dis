package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"cems-dis/model"
	rs "cems-dis/server/response"
	"cems-dis/utils"
)


func (s ApiService) ListDevices(c *gin.Context) rs.Response {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	sql := s.model.DB.Model(&model.Device{}).Order("name")
	sql = model.SetSearchKeywords(sql, []string{"uid", "name", "api_key", "secret"}, c.Query("q"))

	if _, ok := c.GetQuery("disabled"); ok {
		sql.Where("(enabled = ?)", false)
	}

	var paging *rs.Paging
	var err error
	if (page > 0) || (size > 0) {
		paging, err = rs.NewPaging(sql, int(page), int(size))
		if err != nil {
			return rs.Error(http.StatusInternalServerError, err.Error())
		}
		sql = paging.Sql()
	}

	var devices []model.Device
	if err := sql.Find(&devices).Error; err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}

	list := []*model.DeviceOut{}
	for _, d := range devices {
		list = append(list, d.Out())
	}

	result := rs.Success(list)
	if paging != nil {
		result = result.DefaultWithPaging(paging)
	}
	return result
}

func (s ApiService) GetDevice(c *gin.Context) rs.Response {
	resp := s.getDeviceFromUrl(c)
	if resp.IsError() {
		return resp
	}
	device := resp.Data.(*model.Device)
	return rs.Success(device.Out())
}

func (s ApiService) InsertDevice(c *gin.Context) rs.Response {
	device, err := s.extractDeviceData(c)
	if err != nil {
		return rs.Error(http.StatusBadRequest, err.Error())
	}
	err = s.model.DB.Create(device).Error
	if err != nil {
		log.Warningf("Error in InsertDevice: %s, device: %+v\n", err.Error(), device)
		return rs.Error(http.StatusInternalServerError, "Unknown error")
	}
	return rs.Success(device)
}

func (s ApiService) UpdateDevice(c *gin.Context) rs.Response {
	resp := s.getDeviceFromUrl(c)
	if resp.IsError() {
		return resp
	}
	newDevice, err := s.extractDeviceData(c)
	if err != nil {
		return rs.Error(http.StatusBadRequest, err.Error())
	}
	oldDevice := resp.Data.(*model.Device)
	err = s.model.DB.Model(&model.Device{}).Where("uid = ?", oldDevice.UID).Updates(newDevice).Error
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}
	return rs.Success(newDevice)
}

func (s ApiService) DeleteDevice(c *gin.Context) rs.Response {
	uid, err := s.getDeviceUidFromUrl(c)
	if err != nil {
		return rs.Error(http.StatusBadRequest, err.Error())
	}
	if !s.model.IsDeviceExist(uid) {
		return rs.Error(http.StatusNotFound, fmt.Sprintf("UID '%s' tidak ada di database", uid))
	}
	err = s.model.DeleteDeviceByUid(uid)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}
	return rs.Success(nil)
}

func (s ApiService) GenerateDeviceSecret(c *gin.Context) rs.Response {
	return rs.Success(utils.GenerateRandomString(32))
}

func (s ApiService) validateUidString(uid string) error {
	if len(uid) == 0 {
		return errors.New("UID tidak boleh kosong")
	} else if len(uid) > model.MAX_DEVICE_UID_LENGTH {
		return errors.New(fmt.Sprintf("Panjang UID tidak boleh lebih dari %d karakter", model.MAX_DEVICE_UID_LENGTH))
	}
	return nil
}

func (s ApiService) getDeviceUidFromUrl(c *gin.Context) (string, error) {
	uid := strings.Trim(c.Param("uid"), " ")
	if err := s.validateUidString(uid); err != nil {
		return "", errors.New("UID tidak boleh kosong")
	}
	return uid, nil
}

func (s ApiService) getDeviceFromUrl(c *gin.Context) rs.Response {
	uid, err := s.getDeviceUidFromUrl(c)
	if err != nil {
		return rs.Error(http.StatusBadRequest, err.Error())
	}
	device, _ := s.model.GetDeviceByUid(uid)
	if device == nil {
		return rs.Error(http.StatusNotFound, "UID tidak valid")
	}
	return rs.Success(device)
}

func (s ApiService) extractDeviceData(c *gin.Context) (*model.Device, error) {
	oldUid, _ := s.getDeviceUidFromUrl(c)
	device := &model.Device{}
	err := c.BindJSON(device)
	if err != nil {
		log.Warningf(fmt.Sprintf("Error extracting device data: %s\n", err.Error()))
		return nil, errors.New("Invalid JSON body")
	}

	device.Trim()
	if len(device.UID) == 0 {
		return nil, errors.New("UID wajib diisi")
	} else if s.model.IsDeviceUidTaken(device.UID, oldUid) {
		return nil, errors.New(fmt.Sprintf("UID '%s' sudah digunakan", device.UID))
	} else if len(device.Name) == 0 {
		return nil, errors.New("Nama device wajib diisi")
	} else if len(device.ApiKey) == 0 {
		return nil, errors.New("Api key wajib diisi")
	} else if s.model.IsDeviceApiKeyTaken(device.ApiKey, oldUid) {
		return nil, errors.New(fmt.Sprintf("Api key '%s' sudah digunakan", device.ApiKey))
	} else if len(device.Secret) == 0 {
		return nil, errors.New("Secret wajib diisi")
	}

	return device, nil
}
