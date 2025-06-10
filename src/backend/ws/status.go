package ws

import (
	"github.com/gofiber/contrib/websocket"
)

func handleGetStatus(c *websocket.Conn) {
	dmxCtrlMu.RLock()
	dmxInitialized := dmxCtrl != nil
	dmxCtrlMu.RUnlock()

	showMu.Lock()
	showRunning := currentShow != nil
	var showID string
	var showStep int
	var showLoop bool
	if currentShow != nil {
		showID = currentShow.id
		showStep = currentShow.currentStep
		showLoop = currentShow.loop
	}
	showMu.Unlock()

	presetMu.RLock()
	activePreset := activePresetID
	presetMu.RUnlock()

	monitorMu.Lock()
	isMonitoring := monitoring
	monitorMu.Unlock()

	status := map[string]interface{}{
		"dmx_initialized":   dmxInitialized,
		"show_running":      showRunning,
		"active_show_id":    showID,
		"show_step":         showStep,
		"show_loop":         showLoop,
		"active_preset_id":  activePreset,
		"monitoring":        isMonitoring,
		"connected_clients": getClientCount(),
	}

	writeJSON(c, Message{
		Type:    "status",
		Payload: mustMarshal(status),
	})
}
