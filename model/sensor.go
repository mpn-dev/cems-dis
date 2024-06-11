package model

import (
  "errors"
	"fmt"
  "strings"
  "time"
	log "github.com/sirupsen/logrus"
  "cems-dis/utils"
)

type Sensor struct {
  Slot        string      `json:"slot"          gorm:"size:10;primaryKey"`
  Code        string      `json:"code"          gorm:"size:20;index"`
  Name        string      `json:"name"          gorm:"size:30;index"`
  Unit        string      `json:"unit"          gorm:"size:10;index"`
  Enabled     bool        `json:"enabled"       gorm:"index"`
  CreatedAt   time.Time   `json:"created_at"    gorm:"autoCreateTime;index"`
  UpdatedAt   time.Time   `json:"updated_at"    gorm:"autoUpdateTime"`
}

type SensorOut struct {
  Sensor
  CreatedAt   string      `json:"created_at"`
  UpdatedAt   string      `json:"updated_at"`
}

type Sensors []Sensor

func (s Sensor) Trim() Sensor {
  s.Slot = strings.Trim(s.Slot, " ")
  s.Code = strings.Trim(s.Code, " ")
  s.Name = strings.Trim(s.Name, " ")
  s.Unit = strings.Trim(s.Unit, " ")
  return s
}

func (s Sensor) NameAndUnit() string {
	unit := ""
	if len(s.Unit) > 0 {
		unit = fmt.Sprintf(" (%s)", s.Unit)
	}
	return fmt.Sprintf("%s%s", s.Name, unit)
}

func (s Sensor) Out() SensorOut {
  return SensorOut{
    Sensor:     s, 
    CreatedAt:  utils.TimeToString(s.CreatedAt), 
    UpdatedAt:  utils.TimeToString(s.UpdatedAt), 
  }
}

func (ss Sensors) SlotMap() map[string]Sensor {
  result := map[string]Sensor{}
  for _, s := range ss {
    result[s.Slot] = s
  }
  return result
}

func (ss Sensors) CodeMap() map[string]Sensor {
  result := map[string]Sensor{}
  for _, s := range ss {
    result[s.Code] = s
  }
  return result
}

func (m *Model) GetSensors(filters []string) (Sensors, error) {
  sensors := Sensors{}
  sql := m.DB.Order("slot")
  for _, f := range filters {
    sql.Where(f)
  }
  err := sql.Find(&sensors).Error
  if err != nil {
    log.Warningf("DB error: %s", err.Error())
    return nil, errors.New("DB error")
  }
  return sensors, nil
}

func (m *Model) GetAllSensors() (Sensors, error) {
  return m.GetSensors([]string{})
}

func (m *Model) GetActiveSensors() (Sensors, error) {
  return m.GetSensors([]string{"(enabled = true)"})
}
