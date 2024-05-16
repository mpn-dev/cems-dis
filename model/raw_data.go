package model

import (
  "time"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
  "gorm.io/gorm"
)

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
  RawData
  CreatedAt           string      `json:"created_at"`
  UpdatedAt           string      `json:"updated_at"`
}

type SensorValues map[string]*float64

type CemsPayload struct {
	UID					string				`json:"uid"`
	Timestamp		uint64				`json:"timestamp"`
	Values			SensorValues	`json:"values"`
}


func (r *RawData) Out() *RawDataOut {
  return &RawDataOut{
    RawData:          *r, 
    CreatedAt:        r.CreatedAt.Format(DEFAULT_DATE_TIME_FORMAT), 
    UpdatedAt:        r.UpdatedAt.Format(DEFAULT_DATE_TIME_FORMAT), 
  }
}

func (r *RawData) CemsPayload() *CemsPayload {
  return &CemsPayload{
    UID:        r.DEV, 
    Timestamp:  r.Timestamp, 
    Values:     SensorValues{
      "so2":          r.SO2, 
      "nox":          r.NOX, 
      "pm":           r.PM, 
      "h2s":          r.H2S, 
      "opacity":      r.Opacity, 
      "flow":         r.Flow, 
      "o2":           r.O2, 
      "temperature":  r.Temperature, 
      "pressure":     r.Pressure, 
    }, 
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
}

func (m *Model) GetRawDataById(id uint64) (*RawDataOut, error) {
	record := &RawData{}
	if err := m.DB.Model(record).First(record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			log.Warningf("raw_data.GetRawDataById => DB error: %+v", errors.WithStack(err))
			return nil, errors.New("DB error")
		}
	}
	return record.Out(), nil
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
  }
}
