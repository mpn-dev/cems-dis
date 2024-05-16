package model

import (
	"time"
)

type Transmission struct {
	Id							uint64				`json:"id"								gorm:"primaryKey"`
	RawDataId				uint64				`json:"raw_data_id"				gorm:"index"`
	RelayStationId	int64					`json:"relay_station_id"	gorm:"index"`
	Code						int						`json:"code"							gorm:"index"`
	Info						string				`json:"info"							gorm:"size:100"`
	Status					string				`json:"status"						gorm:"size:20"`
	CreatedAt				time.Time			`													gorm:"autoCreateTime"`
	UpdatedAt				time.Time			`													gorm:"autoUpdateTime"`
}

type TransmissionOut struct {
	Transmission 		Transmission
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


func (m *Model) SetTransmissionPending(task Transmission) {
	task.Code = 0
	task.Info = ""
	task.Status = "Started"
	m.DB.Save(&task)
}

func (m *Model) SetTransmissionStarted(task Transmission) {
	task.Code = 0
	task.Info = ""
	task.Status = "Started"
	m.DB.Save(&task)
}

func (m *Model) SetTransmissionSuccess(task Transmission, info string) {
	task.Code = 200
	task.Info = info
	task.Status = "Success"
	m.DB.Save(&task)
}

func (m *Model) SetTransmissionError(task Transmission, code int, error string) {
	task.Code = code
	task.Info = error
	task.Status = "Error"
	m.DB.Save(&task)
}
