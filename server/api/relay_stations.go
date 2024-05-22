package api

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
  "gorm.io/gorm"
	"cems-dis/model"
	rs "cems-dis/server/response"
)


func (s ApiService) RelayStationProtocols(c *gin.Context) rs.Response {
	return rs.Success(model.SupportedRelayProtocols())
}

func (s ApiService) ListRelayStation(c *gin.Context) rs.Response {
	page, _ := strconv.Atoi(c.Param("page"))
	size, _ := strconv.Atoi(c.Param("size"))
	sql := s.model.DB.Model(&model.RelayStation{})
	sql = model.SetSearchKeywords(sql, []string{"name", "protocol"}, c.Param("q"))
	paging, err := rs.NewPaging(sql, page, size)
	if err != nil {
		return rs.Error(http.StatusBadRequest, err.Error())
	}
	var stations []*model.RelayStation
	if err := paging.Sql().Find(&stations).Error; err != nil {
		log.Warningf("DB error: %s", err.Error())
		return rs.Error(http.StatusInternalServerError, "DB error")
	}
	list := []*model.RelayStationOut{}
	for _, r := range stations {
		list = append(list, r.Out())
	}
	return rs.Success(list).DefaultWithPaging(paging)
}

func (s ApiService) GetRelayStation(c *gin.Context) rs.Response {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	station := &model.RelayStation{}
	if err := s.model.DB.Model(station).First(station, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs.Error(http.StatusNotFound, "ID relay station tidak ada di database")
		} else {
			log.Warningf("DB error: %s", err.Error())
			return rs.Error(http.StatusInternalServerError, "DB error")
		}
	}
	return rs.Success(station.Out())
}

func (s ApiService) InsertRelayStation(c *gin.Context) rs.Response {
	resp := s.extractRelayStation(c)
	if resp.IsError() {
		return resp
	}
	station := resp.Data.(*model.RelayStation)
	station.Id = 0
	if err := s.model.DB.Create(station).Error; err != nil {
		log.Warningf("DB error: %s", err.Error())
		return rs.Error(http.StatusInternalServerError, "DB error")
	}
	return rs.Success(station)
}

func (s ApiService) UpdateRelayStation(c *gin.Context) rs.Response {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id <= 0 {
		return rs.Error(http.StatusNotFound, "ID relay station tidak valid")
	}
	resp := s.extractRelayStation(c)
	if resp.IsError() {
		return resp
	}
	station := resp.Data.(*model.RelayStation)
	station.Id = id
	if err := s.model.DB.Save(station).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs.Error(http.StatusNotFound, "ID relay station tidak ada di database")
		} else {
			log.Warningf("DB error: %s", err.Error())
			return rs.Error(http.StatusInternalServerError, "DB error")
		}
	}
	return rs.Success(station)
}

func (s ApiService) DeleteRelayStation(c *gin.Context) rs.Response {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := s.model.DB.Delete(&model.RelayStation{Id: id}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rs.Error(http.StatusNotFound, "ID relay station tidak ada di database")
		} else {
			log.Warningf("DB error: %s", err.Error())
			return rs.Error(http.StatusInternalServerError, "DB error")
		}
	}
	return rs.Success(nil)
}

func (s ApiService) extractRelayStation(c *gin.Context) rs.Response {
	station := &model.RelayStation{Enabled: true}
	if err := c.BindJSON(station); err != nil {
		return rs.Error(http.StatusBadRequest, "Bad JSON body")
	}
	station.Name = strings.Trim(station.Name, " ")
	station.Protocol = strings.Trim(station.Protocol, " ")
	station.BaseURL = strings.Trim(station.BaseURL, " ")
	if len(station.Name) == 0 {
		return rs.Error(http.StatusBadRequest, "Nama wajib diisi")
	} else if len(station.Protocol) == 0 {
		return rs.Error(http.StatusBadRequest, "Protokol wajib diisi")
	} else if !slices.Contains(model.SupportedRelayProtocols(), station.Protocol) {
		return rs.Error(http.StatusBadRequest, "Nama protocol tidak dikenal")
	} else if len(station.BaseURL) == 0 {
		return rs.Error(http.StatusBadRequest, "Base URL wajib diisi")
	}
	return rs.Success(station)
}
