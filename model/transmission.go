package model

import (
	"time"
)

type Transmission struct {
	Id							int64					`json:"id"								gorm:"primaryKey"`
	RelayStationId	int64					`json:"relay_station_id"	gorm:"index"`
	RelayStation	  RelayStation	`json:"relay_station"`
	Code						int						`json:"code"							gorm:"index"`
	Error						string				`json:"error"							gorm:"size:100"`
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
