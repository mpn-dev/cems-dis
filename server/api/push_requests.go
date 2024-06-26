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

func (s ApiService) ListPushRequests(c *gin.Context) rs.Response {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	ts1, _ := strconv.ParseInt(c.Query("ts1"), 10, 64)
	ts2, _ := strconv.ParseInt(c.Query("ts2"), 10, 64)
	sql := s.model.DB.Model(&model.PushRequest{})
	sql = model.SetSearchKeywords(sql, []string{"uid", "ip_addr", "user_agent", "request", "status", "info"}, c.Query("q"))

	if (ts1 > 0) && (ts2 > 0) {
		if ts2 < ts1 {
			return rs.Error(http.StatusBadRequest, "Tanggal akhir tidak boleh kurang dari tanggal awal")
		}
		tm1 := utils.TimeToString(time.Unix(ts1, 0))
		tm2 := utils.TimeToString(time.Unix(ts2, 0))
		sql.Where("(created_at BETWEEN ? AND ?)", tm1, tm2)
	}

	paging, err := rs.NewPaging(sql, page, size)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}

	var requests []*model.PushRequest
	if err := paging.Sql().Order("created_at DESC").Find(&requests).Error; err != nil {
		log.Warningf("DB error: %s", err.Error())
		return rs.Error(http.StatusInternalServerError, "DB error")
	}

	list := []*model.PushRequestOut{}
	for _, r := range requests {
		list = append(list, r.Out())
	}

	return rs.Success(list).DefaultWithPaging(paging)
}
