package api

import (
	"fmt"
	"net/http"
	"strings"
  "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"cems-dis/model"
	rs "cems-dis/server/response"
)

func (s ApiService) ListSensors(c *gin.Context) rs.Response {
	filters := []string{}
	if _, ok := c.Request.URL.Query()["enabled"]; ok {
		filters = append(filters, "(enabled = true)")
	}

	sensors, err := s.model.GetSensors(filters)
	if err != nil {
		return rs.Error(http.StatusInternalServerError, err.Error())
	}
	return rs.Success(sensors)
}

func (s ApiService) UpdateSensor(c *gin.Context) rs.Response {
	sensors := model.Sensors{}
	err := c.BindJSON(&sensors)
	if err != nil {
		return rs.Error(http.StatusBadRequest, "Invalid JSON body")
	}

	for _, ss := range sensors {
		ss = ss.Trim()
		if len(ss.Slot) == 0 {
			return rs.Error(http.StatusBadRequest, "ID slot sensor wajib diisi")
		}
		if !strings.Contains("|s01|s02|s03|s04|s05|s06|s07|s08|s09|s10|s11|s12|", "|" + ss.Slot + "|") {
			return rs.Error(http.StatusBadRequest, fmt.Sprintf("ID slot '%s' tidak valid", ss.Slot))
		}
		if ss.Enabled {
			if len(ss.Code) == 0 {
				return rs.Error(http.StatusBadRequest, "Kode sensor wajib diisi")
			}
			if len(ss.Name) == 0 {
				return rs.Error(http.StatusBadRequest, "Nama sensor wajib diisi")
			}
		}
	}

	for _, ss := range sensors {
		ss = ss.Trim()
		err := s.model.DB.Where("slot = ?", ss.Slot).Save(&ss).Error
		if err != nil {
			log.Warningf("DB error: %s", err.Error())
			return rs.Error(http.StatusInternalServerError, "DB error")
		}
	}

	return rs.Success(nil)
}
