package api

import (
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"cems-dis/model"
	rs "cems-dis/server/response"
	"cems-dis/utils"
)

func (s ApiService) ListTransmissions(c *gin.Context) rs.Response {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	ts_1, _ := strconv.ParseInt(c.Query("ts1"), 10, 64)
	ts_2, _ := strconv.ParseInt(c.Query("ts2"), 10, 64)

	sql := s.model.DB.Model(&model.TransmissionTable{}).
		Select("transmissions.id, raw_data_id, station_id, name station_name, " + 
				"code, status, note, transmissions.created_at, transmissions.updated_at").
		Joins("JOIN relay_stations ON relay_stations.id = transmissions.station_id")
	sql = model.SetSearchKeywords(sql, []string{"relay_stations.name", "status", "note"}, c.Query("q"))

	if (ts_1 > 0) || (ts_2 > 0) {
		if ts_2 < ts_1 {
			return rs.Error(http.StatusBadRequest, "Tanggal akhir tidak boleh kurang dari tanggal awal")
		}
		tm1 := utils.TimeToString(time.Unix(ts_1, 0))
		tm2 := utils.TimeToString(time.Unix(ts_2, 0))
		sql.Where("(transmissions.created_at BETWEEN ? AND ?)", tm1, tm2)
	}

	paging, err := rs.NewPaging(sql, page, size)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}

	var transmissions []*model.Transmission
	if err := paging.Sql().Order("transmissions.created_at DESC").Find(&transmissions).Error; err != nil {
		log.Warningf("DB error: %s", err.Error())
		return rs.Error(http.StatusInternalServerError, "DB error")
	}

	list := []*model.TransmissionOut{}
	for _, t := range transmissions {
		list = append(list, t.Out())
	}
	return rs.Success(list).DefaultWithPaging(paging)
}
