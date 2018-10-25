package collector

import (
	"time"

	"github.com/bcicen/ctop/models"
)

const mockLog = "Cura ob pro qui tibi inveni dum qua fit donec amare illic mea, regem falli contexo pro peregrinorum heremo absconditi araneae meminerim deliciosas actionibus facere modico dura sonuerunt psalmi contra rerum, tempus mala anima volebant dura quae o modis."

type MockLogs struct {
	done chan bool
}

func (l *MockLogs) Stream() chan models.Log {
	logCh := make(chan models.Log)
	go func() {
		for {
			select {
			case <-l.done:
				break
			default:
				logCh <- models.Log{Timestamp: time.Now(), Message: mockLog}
				time.Sleep(250 * time.Millisecond)
			}
		}
	}()
	return logCh
}

func (l *MockLogs) Stop() { l.done <- true }
