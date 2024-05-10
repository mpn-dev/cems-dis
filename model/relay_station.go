package model

import (
  "time"
)

type RelayStation struct {
  Id            int64       `json:"id"          gorm:"primaryKey"`
  Name          string      `json:"name"        gorm:"size:30"`
  Protocol      string      `json:"protocol"    gorm:"size:30"`
  BaseURL       string      `json:"base_url"    gorm:"size:50"`
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
    CreatedAt:    s.CreatedAt.Format(DEFAULT_DATE_TIME_FORMAT), 
    UpdatedAt:    s.UpdatedAt.Format(DEFAULT_DATE_TIME_FORMAT), 
  }
}
