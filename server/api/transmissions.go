package api

import (
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"cems-dis/model"
	rs "cems-dis/server/response"
)

func (s ApiService) ListTransmissions(c *gin.Context) rs.Response {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	ts_1, _ := strconv.ParseInt(c.Query("ts1"), 10, 64)
	ts_2, _ := strconv.ParseInt(c.Query("ts2"), 10, 64)
	sql := s.model.DB.Model(&model.Transmission{})
	sql = model.SetSearchKeywords(
		sql, 
		[]string{"station", "protocol", "base_url", "username", "password", "status", "note"}, 
		c.Query("q"))

	if (ts_1 > 0) || (ts_2 > 0) {
		if ts_2 < ts_1 {
			return rs.Error(http.StatusBadRequest, "Tanggal akhir tidak boleh kurang dari tanggal awal")
		}
		tm1 := time.Unix(ts_1, 0).Format(model.DEFAULT_DATE_TIME_FORMAT)
		tm2 := time.Unix(ts_2, 0).Format(model.DEFAULT_DATE_TIME_FORMAT)
		sql.Where("(created_at BETWEEN ? AND ?)", tm1, tm2)
	}

	paging, err := rs.NewPaging(sql, page, size)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}

	var transmissions []*model.Transmission
	if err := paging.Sql().Order("created_at DESC").Find(&transmissions).Error; err != nil {
		log.Warningf("DB error: %s", err.Error())
		return rs.Error(http.StatusInternalServerError, "DB error")
	}

	list := []*model.TransmissionOut{}
	for _, t := range transmissions {
		list = append(list, t.Out())
	}
	return rs.Success(list).DefaultWithPaging(paging)
}
