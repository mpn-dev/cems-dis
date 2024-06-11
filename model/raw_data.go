package model

import (
  "time"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
  "gorm.io/gorm"
  "cems-dis/utils"
)

type SensorValues map[string]*float64

type RawData struct {
  Id                  uint64          `json:"id"          gorm:"primaryKey"`
  DEV                 string          `json:"uid"         gorm:"column:uid;size:32;index"`
  Timestamp           int64           `json:"timestamp"`
  S01                 *float64
  S02                 *float64
  S03                 *float64
  S04                 *float64
  S05                 *float64
  S06                 *float64
  S07                 *float64
  S08                 *float64
  S09                 *float64
  S10                 *float64
  S11                 *float64
  S12                 *float64
  CreatedAt           time.Time       `json:"created_at"  gorm:"autoCreateTime;index"`
  UpdatedAt           time.Time       `json:"updated_at"  gorm:"autoUpdateTime"`
}

type RawDataOut struct {
  Id                  uint64          `json:"id"`
  DEV                 string          `json:"uid"`
  Timestamp           int64           `json:"timestamp"`
  SensorValues                        `json:"values"`
  CreatedAt           string          `json:"created_at"`
  UpdatedAt           string          `json:"updated_at"`
}

type CemsPayload struct {
  UID                 string          `json:"uid"`
  Timestamp           int64           `json:"timestamp"`
  SensorValues                        `json:"values"`
}

func (v SensorValues) Value(key string) *float64 {
  if vx, ok := v[key]; ok {
    return vx
  }
  return nil
}

func (r *RawData) GetValues(ss Sensors) SensorValues {
  vals := SensorValues{
    "s01": r.S01, 
    "s02": r.S02, 
    "s03": r.S03, 
    "s04": r.S04, 
    "s05": r.S05, 
    "s06": r.S06, 
    "s07": r.S07, 
    "s08": r.S08, 
    "s09": r.S09, 
    "s10": r.S10, 
    "s11": r.S11, 
    "s12": r.S12, 
  }

  sv := SensorValues{}
  for _, s := range ss {
    if s.Enabled {
      sv[s.Code] = vals.Value(s.Slot)
    }
  }

  return sv
}

func (r *RawData) SetValues(ss Sensors, vv SensorValues) {
  vals := SensorValues{}
  for _, s := range ss {
    if s.Enabled {
      vals[s.Slot] = vv.Value(s.Code)
    }
  }

  r.S01 = vals.Value("s01")
  r.S02 = vals.Value("s02")
  r.S03 = vals.Value("s03")
  r.S04 = vals.Value("s04")
  r.S05 = vals.Value("s05")
  r.S06 = vals.Value("s06")
  r.S07 = vals.Value("s07")
  r.S08 = vals.Value("s08")
  r.S09 = vals.Value("s09")
  r.S10 = vals.Value("s10")
  r.S11 = vals.Value("s11")
  r.S12 = vals.Value("s12")
}

func (r *RawData) Out(sensors Sensors) *RawDataOut {
  return &RawDataOut{
    Id:               r.Id, 
    DEV:              r.DEV, 
    Timestamp:        r.Timestamp, 
    SensorValues:     r.GetValues(sensors), 
    CreatedAt:        utils.TimeToString(r.CreatedAt), 
    UpdatedAt:        utils.TimeToString(r.UpdatedAt), 
  }
}

func (r *RawData) CemsPayload(sensors Sensors) *CemsPayload {
  return &CemsPayload{
    UID:            r.DEV, 
    Timestamp:      r.Timestamp, 
    SensorValues:   r.GetValues(sensors), 
  }
}

func (m *Model) GetRawDataById(id uint64) (*RawData, error) {
	record := &RawData{}
	if err := m.DB.Model(record).First(record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			log.Warningf("raw_data.GetRawDataById => DB error: %+v", errors.WithStack(err))
			return nil, errors.New("DB error")
		}
	}

	return record, nil
}

func (m *Model) ParseRawData(data map[string]interface{}) (*RawData, error) {
  result := &RawData{}
  sensors, err := m.GetActiveSensors()
  if err != nil {
    return nil, err
  }

  if val, ok := data["uid"]; ok {
    result.DEV = val.(string)
  }
  if val, ok := data["timestamp"]; ok {
    if ts, ok := val.(float64); ok {
      result.Timestamp = int64(ts)
    }
  }

  values := SensorValues{}
  for _, s := range sensors {
    if val, ok := data[s.Code]; ok {
      if vx, ok := val.(float64); ok {
        values[s.Code] = &vx
      }
    }
  }

  result.SetValues(sensors, values)
  return result, nil
}
