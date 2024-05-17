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
