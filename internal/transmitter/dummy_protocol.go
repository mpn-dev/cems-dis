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
	fmt.Println("")
}

func (p *DummyProtocol) Send(task *model.Transmission) Result {
	dummySend()
	return Success(task, 0, "")
}

func NewDummyProtocol(model *model.Model) *DummyProtocol {
	return &DummyProtocol{
		model: model, 
	}
}
