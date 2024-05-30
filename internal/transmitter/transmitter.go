package transmitter

import (
	"errors"
	"fmt"
	"time"
	log "github.com/sirupsen/logrus"
	"cems-dis/config"
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
	setResult := func (result Result) error {
		if !result.IsSuccess() {
			trx.model.SetTransmissionError(result.Task(), result.Code(), result.Note())
			log.Warningf("transmitter.Start => %s", result.Info())
			return errors.New(result.Info())
		}

		trx.model.SetTransmissionSuccess(result.Task(), result.Code(), result.Note())
		log.Printf("transmitter.Start => %s", result.Info())
		return nil
	}

	trx.stop = false

	log.Println("Transmitter started")
	sleepSecs(5)

	for {
		if trx.stop {
			log.Println("Transmitter stopped")
			break
		}

		timeLimit := time.Now().Add(time.Duration(-config.RetransmitAfterSecs()) * time.Second)
		upd := trx.model.DB.Model(&model.Transmission{}).
			Where("((status = ?) OR (status = ?)) AND (updated_at < ?)", "Started", "Error", timeLimit).
			Update("status", "Pending")
		if err := upd.Error; err != nil {
			log.Warningf("DB error: %s", err.Error())
		}

		var tasks []*model.Transmission
		err := trx.model.DB.Model(model.TransmissionTable{}).
			Select("transmissions.id, raw_data_id, station_id, name station_name,  " + 
				"protocol,base_url, username, password, code, status, note, " + 
				"transmissions.created_at, transmissions.updated_at").
			Joins("JOIN relay_stations ON relay_stations.id = transmissions.station_id").
			Where("status = ?", "Pending").
			Order("updated_at").
			Limit(1).
			Find(&tasks).Error
		if err != nil {
			log.Warningf("DB error: %s", err.Error())
			setResult(Error(nil, 0, fmt.Sprintf("Failed fetching transmission data: %s", err.Error())))
		}
		if len(tasks) == 0 {
			sleepSecs(1)
			continue
		}

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

		fmt.Printf("[Transmit:%s#%d]\n", task.Protocol, task.Id)
		trx.model.SetTransmissionStarted(task)
		err = setResult(protocol.Send(task))
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}
		fmt.Println("Success")
	}
}

func Stop() {
	trx.stop = true
}
