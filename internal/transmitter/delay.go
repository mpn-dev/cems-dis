package transmitter

import (
	"time"
)

func delaySometime() {
	for i := 0; i <= 4; i ++ {
		for j := 0; j <= 10000; j++ {
			for k := 0; k <= 200000; k++ {
				_ = 312987 * 92438 * 1287
			}
		}
	}
}

func sleepSecs(secs int) {
	time.Sleep(time.Duration(secs) * time.Second)
}

func sleepMillisecs(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
