package transmitter

import (
	"fmt"
	"cems-dis/model"
)

type DummyProtocol struct {
	model		*model.Model
}

func dummySend() {
	for i := 0; i <= 4; i ++ {
		fmt.Printf("%d", i + 1)
		for j := 0; j <= 10000; j++ {
			for k := 0; k <= 200000; k++ {
				_ = 312987 * 92438 * 1287
			}
		}
	}
}

func (d *DummyProtocol) Send(task model.Transmission, station model.RelayStation) {
	task.Code = 0
	task.Error = ""
	task.Status = "Started"
	d.model.DB.Save(&task)
	fmt.Println("Transmitting data...")
	dummySend()
	fmt.Println("")
	task.Status = "Success"
	d.model.DB.Save(&task)
}

func NewDummyProtocol(model *model.Model) *DummyProtocol {
	return &DummyProtocol{
		model: model, 
	}
}
