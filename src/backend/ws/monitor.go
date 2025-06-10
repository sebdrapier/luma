package ws

import (
	"log"
	"time"

	"elano.fr/src/backend/driver"
)

func startMonitoring() {
	monitorMu.Lock()
	defer monitorMu.Unlock()
	if monitoring {
		return
	}
	monitoring = true
	monitorStop = make(chan struct{})
	monitorTicker = time.NewTicker(100 * time.Millisecond)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in monitoring: %v", r)
			}
		}()
		for {
			select {
			case <-monitorStop:
				return
			case <-monitorTicker.C:
				broadcastDMXState()
			}
		}
	}()
}

func stopMonitoring() {
	monitorMu.Lock()
	defer monitorMu.Unlock()
	if !monitoring {
		return
	}
	monitoring = false
	close(monitorStop)
	monitorStop = nil
	monitorTicker.Stop()
	monitorTicker = nil
}

func broadcastDMXState() {
	dmxCtrlMu.RLock()
	ctrl := dmxCtrl
	dmxCtrlMu.RUnlock()
	if ctrl == nil {
		return
	}
	channels, err := ctrl.GetAllChannels()
	if err != nil {
		return
	}
	var channelStates []ChannelState
	for i := 0; i < driver.DMXChannels; i++ {
		if channels[i] > 0 {
			channelStates = append(channelStates, ChannelState{Address: i + 1, Value: channels[i]})
		}
	}
	presetMu.RLock()
	activePreset := activePresetID
	presetMu.RUnlock()
	showMu.Lock()
	var activeShow string
	var showStep int
	var showLoop bool
	if currentShow != nil {
		activeShow = currentShow.id
		showStep = currentShow.currentStep
		showLoop = currentShow.loop
	}
	showMu.Unlock()
	state := DMXState{
		Channels:       channelStates,
		ActivePresetID: activePreset,
		ActiveShowID:   activeShow,
		ShowStep:       showStep,
		ShowLoop:       showLoop,
		Timestamp:      time.Now().UnixMilli(),
	}
	select {
	case broadcast <- Message{Type: "dmx_update", Payload: mustMarshal(state)}:
	default:
	}
}
