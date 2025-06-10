package ws

import (
	"sync"
	"time"

	"elano.fr/src/backend/driver"
	"elano.fr/src/backend/storage"
)

var (
	clients        sync.Map // *websocket.Conn -> *sync.Mutex
	broadcast      = make(chan Message, 100)
	dmxCtrl        *driver.DMXController
	dmxCtrlMu      sync.RWMutex
	currentShow    *ShowController
	showMu         sync.Mutex
	activePresetID string
	presetMu       sync.RWMutex
	projectStore   storage.ProjectStore
	monitorTicker  *time.Ticker
	monitoring     bool
	monitorMu      sync.Mutex
	monitorStop    chan struct{}
)
