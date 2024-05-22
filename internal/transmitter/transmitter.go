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

	log.Println("Transmitter started")

	for {
		if trx.stop {
			log.Println("Transmitter stopped")
			break
		}

		sleepSecs(1)
		var tasks []*model.Transmission
		err := trx.model.DB.Model(model.Transmission{}).Where("status = ?", "Pending").Order("updated_at").Limit(1).Find(&tasks).Error
		if err != nil {
			setResult(Error(nil, 0, fmt.Sprintf("Failed fetching transmission data: %s", err.Error())))
		}
		if len(tasks) > 0 {
			var protocol Protocol
			task := tasks[0]
			if task.Protocol == "DUMMY" {
				protocol = NewDummyProtocol(trx.model)
			} else if task.Protocol == "PING" {
				protocol = NewPingProtocol(trx.model)
			} else if task.Protocol == "CEMS-MPN" {
				protocol = NewCemsMpnProtocol(trx.model)
			// } else if task.Protocol == "CEMS-KLHK" {
			// 	protocol = NewKlhkProtocol(trx.model)
			}

			if protocol == nil {
				setResult(Error(task, 0, "Unsupported protocol"))
				continue
			}

			trx.model.SetTransmissionStarted(task)
			setResult(protocol.Send(task))
		}
	}
}

func Stop() {
	trx.stop = true
}

func setResult(result Result) {
	if result.IsSuccess() {
		trx.model.SetTransmissionSuccess(result.Task(), result.Code(), result.Note())
		log.Printf("transmitter.Start => %s", result.Info())
	} else {
		trx.model.SetTransmissionError(result.Task(), result.Code(), result.Note())
		log.Warningf("transmitter.Start => %s", result.Info())
	}
}
