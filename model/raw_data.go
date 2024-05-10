package model

import "time"

type RawData struct {
  Id                  uint64      `gorm:"primaryKey"`
  DEV                 string      `gorm:"column:uid;size:32;index"`
  Timestamp           uint64
  SO2                 *float64
  NOX                 *float64
  PM                  *float64
  H2S                 *float64    `gorm:"column:h2s"`
  Opacity             *float64
  Flow                *float64
  O2                  *float64
  Temperature         *float64
  Pressure            *float64
  Relayed             int
  CreatedAt           time.Time   `gorm:"autoCreateTime"`
  UpdatedAt           time.Time   `gorm:"autoUpdateTime"`
}

type RawDataIn struct {
  Timestamp           uint64      `json:"timestamp"`
  SO2                 *float64    `json:"so2"`
  NOX                 *float64    `json:"nox"`
  PM                  *float64    `json:"pm"`
  H2S                 *float64    `json:"h2s"`
  Opacity             *float64    `json:"opacity"`
  Flow                *float64    `json:"flow"`
  O2                  *float64    `json:"o2"`
  Temperature         *float64    `json:"temperature"`
  Pressure            *float64    `json:"pressure"`
}

type RawDataOut struct {
  Id                  uint64      `json:"id"`
  DEV                 string      `json:"uid"`
  Timestamp           uint64      `json:"timestamp"`
  SO2                 *float64    `json:"so2"`
  NOX                 *float64    `json:"nox"`
  PM                  *float64    `json:"pm"`
  H2S                 *float64    `json:"h2s"`
  Opacity             *float64    `json:"opacity"`
  Flow                *float64    `json:"flow"`
  O2                  *float64    `json:"o2"`
  Temperature         *float64    `json:"temperature"`
  Pressure            *float64    `json:"pressure"`
  Relayed             int         `json:"relayed"`
  CreatedAt           string      `json:"created_at"`
  UpdatedAt           string      `json:"updated_at"`
}


func (r *RawData) Out() *RawDataOut {
  return &RawDataOut{
    Id:               r.Id, 
    DEV:              r.DEV, 
    Timestamp:        r.Timestamp, 
    SO2:              r.SO2, 
    NOX:              r.NOX, 
    PM:               r.PM, 
    H2S:              r.H2S, 
    Opacity:          r.Opacity, 
    Flow:             r.Flow, 
    O2:               r.O2, 
    Temperature:      r.Temperature, 
    Pressure:         r.Pressure, 
    Relayed:          r.Relayed, 
    CreatedAt:        r.CreatedAt.Format(DEFAULT_DATE_TIME_FORMAT), 
    UpdatedAt:        r.UpdatedAt.Format(DEFAULT_DATE_TIME_FORMAT), 
  }
}

func (r *RawData) Update(f *RawData) {
  r.DEV               = f.DEV
  r.Timestamp         = f.Timestamp
  r.SO2               = f.SO2
  r.NOX               = f.NOX
  r.PM                = f.PM
  r.H2S               = f.H2S
  r.Opacity           = f.Opacity
  r.Flow              = f.Flow
  r.O2                = f.O2
  r.Temperature       = f.Temperature
  r.Pressure          = f.Pressure
  r.Relayed           = f.Relayed
}

func NewRawData(uid string, i *RawDataIn) *RawData {
  return &RawData{
    DEV:              uid, 
    Timestamp:        i.Timestamp, 
    SO2:              i.SO2, 
    NOX:              i.NOX, 
    PM:               i.PM, 
    H2S:              i.H2S, 
    Opacity:          i.Opacity, 
    Flow:             i.Flow, 
    O2:               i.O2, 
    Temperature:      i.Temperature, 
    Pressure:         i.Pressure, 
    Relayed:          0, 
  }
}
