package model

import (
  "time"
  "cems-dis/utils"
)

type PushRequest struct {
  Id            uint64			`json:"id"          gorm:"primaryKey"`
  IpAddr        string			`json:"ip_addr"     gorm:"size:15"`
  UserAgent     string      `json:"user_agent"`
  DEV           string			`json:"uid"         gorm:"column:uid;size:32;index"`
  Request       string      `json:"request"`
  Status        string			`json:"status"      gorm:"size:20;index"`
  Info          string			`json:"info"`
  CreatedAt     time.Time		`json:"created_at"  gorm:"autoCreateTime;index"`
}

type PushRequestOut struct {
  PushRequest
  CreatedAt     string      `json:"created_at"`
}

func (r *PushRequest) Out() *PushRequestOut {
  return &PushRequestOut{
    PushRequest:  *r, 
    CreatedAt:    utils.TimeToString(r.CreatedAt), 
  }
}
