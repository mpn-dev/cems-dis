package model

import (
	"time"
)

type Transmission struct {
  Id              uint64        `json:"id"                gorm:"primaryKey"`
  RawDataId       uint64        `json:"raw_data_id"       gorm:"index"`
  Station         string        `json:"station"           gorm:"size:30;index"`
  Protocol        string        `json:"protocol"          gorm:"size:30;index"`
  BaseURL         string        `json:"base_url"          gorm:"size:50;index"`
  Username        string        `json:"username"          gorm:"size:30;index"`
  Password        string        `json:"password"          gorm:"size:30;index"`
  Code            int           `json:"code"              gorm:"index"`
  Status          string        `json:"status"            gorm:"size:20"`
  Note            string        `json:"note"              gorm:"size:100"`
  CreatedAt       time.Time     `json:"created_at"        gorm:"autoCreateTime;index"`
  UpdatedAt       time.Time     `json:"updated_at"        gorm:"autoUpdateTime;index"`
}

type TransmissionOut struct {
	Transmission
	CreatedAt				string				`json:"created_at"`
	UpdatedAt				string				`json:"updated_at"`
}


func (t *Transmission) Out() *TransmissionOut {
	return &TransmissionOut{
		Transmission:	*t, 
		CreatedAt: 		t.CreatedAt.Format(DEFAULT_DATE_TIME_FORMAT), 
		UpdatedAt:		t.UpdatedAt.Format(DEFAULT_DATE_TIME_FORMAT), 
	}
}

func SupportedRelayProtocols() []string {
	return []string{"CEMS-MPN", "CEMS-KLHK"}
}


func (m *Model) SetTransmissionPending(task *Transmission) {
	task.Code = 0
	task.Note = ""
	task.Status = "Started"
	m.DB.Save(&task)
}

func (m *Model) SetTransmissionStarted(task *Transmission) {
	task.Code = 0
	task.Note = ""
	task.Status = "Started"
	m.DB.Save(&task)
}

func (m *Model) SetTransmissionSuccess(task *Transmission, code int, note string) {
	task.Code = code
	task.Note = note
	task.Status = "Success"
	m.DB.Save(&task)
}

func (m *Model) SetTransmissionError(task *Transmission, code int, error string) {
	task.Code = code
	task.Note = error
	task.Status = "Error"
	m.DB.Save(&task)
}
