package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
  "gorm.io/gorm"

	"cems-dis/model"
	rs "cems-dis/server/response"
)


func (s ApiService) GetRawDataById(c *gin.Context) rs.Response {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	record, err := s.model.GetRawDataById(id)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, "Invalid record ID")
	}
	if record == nil {
		return rs.Error(http.StatusNotFound, "Record ID not found")
	}
	sensors, err := s.model.GetActiveSensors()
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}
	return rs.Success(record.Out(sensors))
}

func (s ApiService) GetLatestData(c *gin.Context) rs.Response {
	res := s.getDeviceFromUrl(c)
	if res.IsError() {
		return res
	}
	sensors, err := s.model.GetActiveSensors()
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}
  device := res.Data.(*model.Device)
	rawData := &model.RawData{}
	err = s.model.DB.Where("uid = ?", device.UID).Order("timestamp DESC").Limit(1).First(rawData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs.Success(nil)
		}

		log.Warningf("DB error: %s", err.Error())
		return rs.Error(http.StatusInternalServerError, "DB error")
	}
	return rs.Success(rawData.Out(sensors))
}

func (s ApiService) GetChartData(c *gin.Context) rs.Response {
	res := s.getDeviceFromUrl(c)
	if res.IsError() {
		return res
	}
	rows := []model.RawData{}
	device := res.Data.(*model.Device)
	latest := &model.RawData{}
	err := s.model.DB.Where("uid = ?", device.UID).Order("timestamp DESC").Limit(1).First(latest).Error
	if err != nil {
		return rs.Success(rows)
	}
	ts1 := latest.Timestamp - (latest.Timestamp % 86400)
	ts2 := ts1 + 86400 - 1
	err = s.model.DB.Where("(uid = ?) AND (timestamp BETWEEN ? AND ?)", device.UID, ts1, ts2).Order("timestamp").Find(&rows).Error
	if err != nil {
		log.Warningf("DB error: %s", err.Error())
	}
	return rs.Success(rows)
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

	sensors, err := s.model.GetActiveSensors()
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}

	device := &model.Device{UID: uid}
	err = s.model.DB.Model(&model.Device{}).Find(device).Error
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

	paging, err := rs.NewPaging(sql, page, size)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}
	var records []*model.RawData
	if err = paging.Sql().Find(&records).Error; err != nil {
		log.Warningf("DB error: %s", err.Error())
	}
	list := []*model.RawDataOut{}
	for _, r := range records {
		list = append(list, r.Out(sensors))
	}
	return rs.Success(list).DefaultWithPaging(paging)
}

func (s ApiService) ListEmissionData(c *gin.Context) rs.Response {
	return rs.Success(nil)
}

func (s ApiService) ListPercentageData(c *gin.Context) rs.Response {
	return rs.Success(nil)
}
