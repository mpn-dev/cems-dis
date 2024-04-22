package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"cems-dis/model"
	rs "cems-dis/server/response"
	"cems-dis/utils"
)


func (s ApiService) ListDevices(c *gin.Context) rs.Response {
	var devices []model.Device

	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	size, _ := strconv.ParseInt(c.Query("size"), 10, 64)

	sql := s.model.DB.Model(&model.Device{}).Order("name")
	sql = model.SetSearchKeywords(sql, []string{"uid", "name", "api_key", "secret"}, c.Query("q"))

	if _, ok := c.GetQuery("disabled"); ok {
		sql.Where("(enabled = ?)", false)
	}

	paging, err := rs.NewPagingFormatter(sql, int(size), int(page))
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}
	err = paging.Sql().Find(&devices).Error
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}
	list := []*model.DeviceOut{}
	for _, d := range devices {
		list = append(list, d.Out())
	}
	return rs.Success(list).UseFormatter(paging)
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
	tmpDevice, err := s.extractDeviceData(c)
	if err != nil {
		return rs.Error(http.StatusBadRequest, err.Error())
	}
	newDevice := resp.Data.(*model.Device).Copy()
	newDevice.Update(tmpDevice)
	err = s.model.DB.Save(newDevice).Error
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}
	return rs.Success(newDevice)
}

func (s ApiService) DeleteDevice(c *gin.Context) rs.Response {
	id, err := s.getDeviceIdFromUrl(c)
	if err != nil {
		return rs.Error(http.StatusBadRequest, err.Error())
	}
	if !s.model.IsDeviceExist(id) {
		return rs.Error(http.StatusNotFound, fmt.Sprintf("ID device '%d' tidak ada di database", id))
	}
	err = s.model.DeleteDeviceById(id)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}
	return rs.Success(nil)
}

func (s ApiService) GenerateDeviceSecret(c *gin.Context) rs.Response {
	return rs.Success(utils.GenerateRandomString(32))
}

func (s ApiService) getDeviceIdFromUrl(c *gin.Context) (uint64, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id <= 0 {
		return 0, errors.New(fmt.Sprintf("ID device '%d' tidak valid", id))
	}
	return id, nil
}

func (s ApiService) getDeviceFromUrl(c *gin.Context) rs.Response {
	id, err := s.getDeviceIdFromUrl(c)
	if err != nil {
		return rs.Error(http.StatusBadRequest, err.Error())
	}
	device, err := s.model.GetDeviceById(id)
	if err != nil {
		return rs.Error(http.StatusNotFound, err.Error())
	}
	return rs.Success(device)
}

func (s ApiService) extractDeviceData(c *gin.Context) (*model.Device, error) {
	id, _ := s.getDeviceIdFromUrl(c)
	device := &model.Device{}
	err := c.BindJSON(device)
	if err != nil {
		log.Warningf(fmt.Sprintf("Error extracting device data: %s\n", err.Error()))
		return nil, errors.New("Invalid JSON body")
	}

	device.Trim()
	if len(device.UID) == 0 {
		return nil, errors.New("UID wajib diisi")
	} else if s.model.IsDeviceUidTaken(device.UID, id) {
		return nil, errors.New(fmt.Sprintf("UID '%s' sudah digunakan", device.UID))
	} else if len(device.Name) == 0 {
		return nil, errors.New("Nama device wajib diisi")
	} else if len(device.ApiKey) == 0 {
		return nil, errors.New("Api key wajib diisi")
	} else if s.model.IsDeviceApiKeyTaken(device.ApiKey, id) {
		return nil, errors.New(fmt.Sprintf("Api key '%s' sudah digunakan", device.ApiKey))
	} else if len(device.Secret) == 0 {
		return nil, errors.New("Secret wajib diisi")
	}

	return device, nil
}
