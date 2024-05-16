package transmitter

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"cems-dis/model"
)

type transmitter struct {
	model		*model.Model
	stop		bool
}

var trx *transmitter

func Init(model *model.Model) {
	if trx == nil {
		trx = &transmitter{
			model: 	model, 
			stop:		false, 
		}
	}
}

func Start() {
	trx.stop = false

	for {
		if trx.stop {
			fmt.Println("Transmitter stopped")
			break
		}

		sleepSecs(1)
		var tasks []*model.Transmission
		err := trx.model.DB.Model(model.Transmission{}).Where("status = ?", "Pending").Order("updated_at").Limit(1).Find(&tasks).Error
		if err != nil {
			log.Warningf("transmitter.Start => Error fetching transmission data: %s", err.Error())
		}
		if len(tasks) > 0 {
			task := *tasks[0]
			station, _ := trx.model.GetRelayStationById(task.RelayStationId)
			if station == nil {
				trx.model.SetTransmissionError(task, 0, "Invalid relay station ID")
				continue
			}
			if !station.Enabled {
				trx.model.SetTransmissionError(task, 0, "Station disabled")
				continue
			}

			var protocol Protocol
			if station.Protocol == "DUMMY" {
				protocol = NewDummyProtocol(trx.model)
			} else if station.Protocol == "PING" {
				protocol = NewPingProtocol(trx.model)
			} else if station.Protocol == "CEMS-MPN" {
				protocol = NewCemsMpnProtocol(trx.model)
			// } else if task.Protocol == "CEMS-KLHK" {
			// 	protocol = NewKlhkProtocol(trx.model)
			}

			if protocol == nil {
				trx.model.SetTransmissionError(task, 0, "Unsupported protocol")
				continue
			}

			fmt.Printf("[Send task #%d]\n", task.Id)
			protocol.Send(task, *station)
		}
	}
}

func Stop() {
	trx.stop = true
}
