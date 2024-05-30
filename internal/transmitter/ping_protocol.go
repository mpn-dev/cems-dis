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
	fmt.Printf("Code: %d, Body: %s\n", res.StatusCode, string(res_body))
}

func (p *PingProtocol) Send(task *model.Transmission) Result {
	sendPing()
	return Success(task, 200, "")
}

func NewPingProtocol(model *model.Model) *PingProtocol {
	return &PingProtocol{
		model: model, 
	}
}
