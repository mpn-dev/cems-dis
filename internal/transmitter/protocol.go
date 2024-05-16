package transmitter

import (
	"cems-dis/model"
)

type Protocol interface {
	Send(t model.Transmission, s model.RelayStation)
}
