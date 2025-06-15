package ws

import (
	"log"
	"time"

	"elano.fr/src/backend/dmx"
	"elano.fr/src/backend/storage"
)

func SetProjectStore(store storage.ProjectStore) {
	projectStore = store
}

func InitializeDMXController(portName string) error {
	dmxCtrlMu.Lock()
	defer dmxCtrlMu.Unlock()
	if dmxCtrl != nil {
		if err := dmxCtrl.Close(); err != nil {
			log.Printf("Error closing existing DMX controller: %v", err)
		}
	}
	ctrl, err := dmx.NewDMXController(portName)
	if err != nil {
		return err
	}
	dmxCtrl = ctrl
	if getClientCount() > 0 {
		go func() {
			time.Sleep(100 * time.Millisecond)
			startMonitoring()
		}()
	}
	return nil
}

func CloseDMXController() error {
	stopMonitoring()
	dmxCtrlMu.Lock()
	defer dmxCtrlMu.Unlock()
	if dmxCtrl != nil {
		err := dmxCtrl.Close()
		dmxCtrl = nil
		return err
	}
	return nil
}
