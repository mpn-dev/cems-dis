package model

import (
  "time"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
  "gorm.io/gorm"
  "cems-dis/utils"
)

type RelayStation struct {
  Id            int64       `json:"id"          gorm:"primaryKey"`
  Name          string      `json:"name"        gorm:"size:30"`
  Protocol      string      `json:"protocol"    gorm:"size:30"`
  BaseURL       string      `json:"base_url"    gorm:"size:50"`
  Username      string      `json:"username"    gorm:"size:30"`
  Password      string      `json:"password"    gorm:"size:30"`
  Enabled       bool        `json:"enabled"`
  CreatedAt     time.Time   `json:"created_at"  gorm:"autoCreateTime"`
  UpdatedAt     time.Time   `json:"updated_at"  gorm:"autoUpdateTime"`
}

type RelayStationOut struct {
  RelayStation
  CreatedAt     string      `json:"created_at"`
  UpdatedAt     string      `json:"updated_at"`
}

func (s *RelayStation) Out() *RelayStationOut {
  return &RelayStationOut{
    RelayStation: *s, 
    CreatedAt:    utils.TimeToString(s.CreatedAt), 
    UpdatedAt:    utils.TimeToString(s.UpdatedAt), 
  }
}

func (m *Model) GetRelayStationById(id int64) (*RelayStation, error) {
  station := &RelayStation{}
  if err := m.DB.Model(station).First(station, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      return nil, nil
    }
    log.Warningf("DB error: %s", err.Error())
    return nil, errors.New("DB error")
  }
  return station, nil
}
