package model

import (
	"time"
	"cems-dis/utils"
)

type TransmissionTable struct {
  Id              uint64        `gorm:"primaryKey"`
  RawDataId       uint64        `gorm:"index"`
  StationId       int64         `gorm:"size:30;index"`
  Code            int           `gorm:"index"`
  Status          string        `gorm:"size:20"`
  Note            string        `gorm:"size:100"`
  CreatedAt       time.Time     `gorm:"autoCreateTime;index"`
  UpdatedAt       time.Time     `gorm:"autoUpdateTime;index"`
}

type Transmission struct {
  Id              uint64        `json:"id"`
  RawDataId       uint64        `json:"raw_data_id"`
  StationId       int64         `json:"station_id"`
  StationName     string        `json:"station_name"`
  Protocol        string        `json:"protocol"`
  BaseURL         string        `json:"base_url"`
  Username        string        `json:"username"`
  Password        string        `json:"password"`
  Code            int           `json:"code"`
  Status          string        `json:"status"`
  Note            string        `json:"note"`
  CreatedAt       time.Time     `json:"created_at"`
  UpdatedAt       time.Time     `json:"updated_at"`
}

type TransmissionOut struct {
  Id              uint64        `json:"id"`
  RawDataId       uint64        `json:"raw_data_id"`
  StationId       int64         `json:"station_id"`
  StationName     string        `json:"station_name"`
  Code            int           `json:"code"`
  Status          string        `json:"status"`
  Note            string        `json:"note"`
  CreatedAt       string        `json:"created_at"`
  UpdatedAt       string        `json:"updated_at"`
}

func (t TransmissionTable) TableName() string {
	return "transmissions"
}

func (t *Transmission) Out() *TransmissionOut {
  return &TransmissionOut{
    Id:           t.Id, 
    RawDataId:    t.RawDataId, 
    StationId:    t.StationId, 
    StationName:  t.StationName, 
    Code:         t.Code, 
    Status:       t.Status, 
    Note:         t.Note, 
    CreatedAt:    utils.TimeToString(t.CreatedAt), 
    UpdatedAt:    utils.TimeToString(t.UpdatedAt), 
  }
}

func SupportedRelayProtocols() []string {
	return []string{"CEMS-MPN", "CEMS-KLHK"}
}


func (m *Model) updateTransmissionStatus(task *Transmission, code int, status string, note string) {
	task.Code = code
	task.Status = status
	task.Note = note
	m.DB.Model(&TransmissionTable{Id: task.Id}).Updates(map[string]interface{}{
		"code": 	code, 
		"status":	status, 
		"note": 	note, 
	})
}

func (m *Model) SetTransmissionPending(task *Transmission) {
	m.updateTransmissionStatus(task, 0, "Pending", "")
}

func (m *Model) SetTransmissionStarted(task *Transmission) {
	m.updateTransmissionStatus(task, 0, "Started", "")
}

func (m *Model) SetTransmissionSuccess(task *Transmission, code int, note string) {
	m.updateTransmissionStatus(task, code, "Success", note)
}

func (m *Model) SetTransmissionError(task *Transmission, code int, error string) {
	m.updateTransmissionStatus(task, code, "Error", error)
}
