package transmitter

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"cems-dis/model"
)

type PingProtocol struct {
	model		*model.Model
}

func sendPing() {
	url := "http://mpn-monitoring.com/ping"
	body := []byte{}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	// req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	res_body, _ := io.ReadAll(res.Body)
	fmt.Printf("status code: %v\n", res.StatusCode)
	fmt.Printf("body: %s\n", string(res_body))
}

func (p *PingProtocol) Send(task model.Transmission, station model.RelayStation) {
	task.Code = 0
	task.Error = ""
	task.Status = "Started"
	p.model.DB.Save(&task)
	fmt.Println("Transmitting data...")
	sendPing()
	task.Status = "Success"
	p.model.DB.Save(&task)
}

func NewPingProtocol(model *model.Model) *PingProtocol {
	return &PingProtocol{
		model: model, 
	}
}
