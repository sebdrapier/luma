package ws

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"elano.fr/src/backend/dmx"
	"elano.fr/src/backend/models"
	"github.com/gofiber/contrib/websocket"
)

func handleWebSocket(c *websocket.Conn) {
	clients.Store(c, &sync.Mutex{})
	clientID := c.Locals("id")
	log.Printf("WebSocket client connected: %v", clientID)
	sendCurrentState(c)
	if getClientCount() == 1 {
		startMonitoring()
	}
	defer func() {
		clients.Delete(c)
		log.Printf("WebSocket client disconnected: %v", clientID)
		c.Close()
		if getClientCount() == 0 {
			stopMonitoring()
		}
	}()
	c.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-pingTicker.C:
				if err := writeMessage(c, websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()
	for {
		_, msgData, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		var msg Message
		if err := json.Unmarshal(msgData, &msg); err != nil {
			sendError(c, "parse_error", "Invalid message format", err.Error())
			continue
		}
		go handleMessage(c, msg)
	}
}

func handleMessage(c *websocket.Conn, msg Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in handleMessage: %v", r)
			sendError(c, "internal_error", "Internal server error", "")
		}
	}()
	switch msg.Type {
	case "apply_preset":
		handleApplyPreset(c, msg.Payload)
	case "run_show":
		handleRunShow(c, msg.Payload)
	case "stop_show":
		handleStopShow(c)
	case "update_channel":
		handleUpdateChannel(c, msg.Payload)
	case "blackout":
		handleBlackout(c)
	case "get_status":
		handleGetStatus(c)
	case "get_dmx_state":
		handleGetDMXState(c)
	case "start_monitoring":
		handleStartMonitoring(c)
	case "stop_monitoring":
		handleStopMonitoring(c)
	case "get_project_config":
		handleGetProjectConfig(c)
	default:
		sendError(c, "unknown_type", "Unknown message type", msg.Type)
	}
}

func sendCurrentState(c *websocket.Conn) {
	handleGetDMXState(c)
	handleGetProjectConfig(c)
}

func handleApplyPreset(c *websocket.Conn, payload json.RawMessage) {
	var idPayload struct {
		PresetID string `json:"preset_id"`
	}
	var preset PresetPayload
	var presetID string
	if err := json.Unmarshal(payload, &idPayload); err == nil && idPayload.PresetID != "" {
		if projectStore != nil {
			if project := projectStore.Get(); project != nil {
				for _, p := range project.Presets {
					if p.ID == idPayload.PresetID {
						preset = make(PresetPayload)
						for _, ch := range p.Channels {
							preset[strconv.Itoa(ch.DMXAddress)] = int(ch.Value)
						}
						presetID = p.ID
						break
					}
				}
			}
		}
		if preset == nil {
			sendError(c, "preset_not_found", "Preset not found", idPayload.PresetID)
			return
		}
	} else {
		if err := json.Unmarshal(payload, &preset); err != nil {
			sendError(c, "invalid_payload", "Invalid preset payload", err.Error())
			return
		}
	}
	dmxCtrlMu.RLock()
	ctrl := dmxCtrl
	dmxCtrlMu.RUnlock()
	if ctrl == nil {
		sendError(c, "dmx_error", "DMX controller not initialized", "")
		return
	}
	if err := ctrl.Blackout(); err != nil {
		sendError(c, "dmx_error", "Failed to blackout", err.Error())
		return
	}
	channels := make(map[int]byte)
	for addrStr, val := range preset {
		addr, err := strconv.Atoi(addrStr)
		if err != nil || addr < 1 || addr > 512 || val < 0 || val > 255 {
			sendError(c, "invalid_payload", "Invalid channel data", addrStr)
			return
		}
		channels[addr] = byte(val)
	}
	if err := ctrl.SetChannels(channels); err != nil {
		sendError(c, "dmx_error", "Failed to set channels", err.Error())
		return
	}
	presetMu.Lock()
	activePresetID = presetID
	presetMu.Unlock()
	showMu.Lock()
	if currentShow != nil {
		currentShow.cancel()
		currentShow = nil
	}
	showMu.Unlock()
	broadcast <- Message{Type: "preset_applied", Payload: mustMarshal(map[string]interface{}{"preset_id": presetID, "channels": preset})}
}

func handleRunShow(c *websocket.Conn, payload json.RawMessage) {
	var idPayload struct {
		ShowID string `json:"show_id"`
		Loop   bool   `json:"loop"`
	}
	var show ShowPayload
	var showID string
	var showModel *models.Show
	if err := json.Unmarshal(payload, &idPayload); err == nil && idPayload.ShowID != "" {
		if projectStore != nil {
			if project := projectStore.Get(); project != nil {
				for _, s := range project.Shows {
					if s.ID == idPayload.ShowID {
						copy := s
						showModel = &copy
						showID = s.ID
						show.Loop = idPayload.Loop
						show.Steps = make([]ShowStep, len(s.Steps))
						for i, step := range s.Steps {
							for _, p := range project.Presets {
								if p.ID == step.PresetID {
									preset := make(PresetPayload)
									for _, ch := range p.Channels {
										preset[strconv.Itoa(ch.DMXAddress)] = int(ch.Value)
									}
									show.Steps[i] = ShowStep{Preset: preset, DelayMs: step.DelayMS, FadeMs: step.FadeMS}
									break
								}
							}
						}
						break
					}
				}
			}
		}
		if showModel == nil {
			sendError(c, "show_not_found", "Show not found", idPayload.ShowID)
			return
		}
	} else {
		if err := json.Unmarshal(payload, &show); err != nil {
			sendError(c, "invalid_payload", "Invalid show payload", err.Error())
			return
		}
	}
	if len(show.Steps) == 0 {
		sendError(c, "invalid_show", "Show must have at least one step", "")
		return
	}
	dmxCtrlMu.RLock()
	ctrl := dmxCtrl
	dmxCtrlMu.RUnlock()
	if ctrl == nil {
		sendError(c, "dmx_error", "DMX controller not initialized", "")
		return
	}
	showMu.Lock()
	if currentShow != nil {
		currentShow.cancel()
		showMu.Unlock()
		time.Sleep(50 * time.Millisecond)
		showMu.Lock()
	}
	ctx, cancel := context.WithCancel(context.Background())
	currentShow = &ShowController{cancel: cancel, id: showID, currentStep: 0, showData: showModel, loop: show.Loop}
	showMu.Unlock()
	presetMu.Lock()
	activePresetID = ""
	presetMu.Unlock()
	broadcast <- Message{Type: "show_started", Payload: mustMarshal(map[string]interface{}{"show_id": showID, "steps": len(show.Steps), "loop": show.Loop})}
	go runShowSequence(ctx, ctrl, show, showID)
}

func handleStopShow(c *websocket.Conn) {
	showMu.Lock()
	if currentShow != nil {
		currentShow.cancel()
		currentShow = nil
	}
	showMu.Unlock()
	writeJSON(c, Message{Type: "show_stopped", Payload: json.RawMessage("{}")})
}

func handleUpdateChannel(c *websocket.Conn, payload json.RawMessage) {
	var u ChannelUpdatePayload
	if err := json.Unmarshal(payload, &u); err != nil {
		sendError(c, "invalid_payload", "Invalid channel update payload", err.Error())
		return
	}
	if u.DMXAddress < 1 || u.DMXAddress > 512 || u.Value < 0 || u.Value > 255 {
		sendError(c, "invalid_payload", "Channel update out of range", "")
		return
	}
	dmxCtrlMu.RLock()
	ctrl := dmxCtrl
	dmxCtrlMu.RUnlock()
	if ctrl == nil {
		sendError(c, "dmx_error", "DMX controller not initialized", "")
		return
	}
	if err := ctrl.SetChannel(u.DMXAddress, byte(u.Value)); err != nil {
		sendError(c, "dmx_error", "Failed to set channel", err.Error())
		return
	}
	presetMu.Lock()
	activePresetID = ""
	presetMu.Unlock()
	broadcast <- Message{Type: "channel_update", Payload: payload}
}

func handleBlackout(c *websocket.Conn) {
	dmxCtrlMu.RLock()
	ctrl := dmxCtrl
	dmxCtrlMu.RUnlock()
	if ctrl == nil {
		sendError(c, "dmx_error", "DMX controller not initialized", "")
		return
	}
	if err := ctrl.Blackout(); err != nil {
		sendError(c, "dmx_error", "Failed to blackout", err.Error())
		return
	}
	presetMu.Lock()
	activePresetID = ""
	presetMu.Unlock()
	showMu.Lock()
	if currentShow != nil {
		currentShow.cancel()
		currentShow = nil
	}
	showMu.Unlock()
	broadcast <- Message{Type: "blackout", Payload: json.RawMessage("{}")}
}

func handleGetDMXState(c *websocket.Conn) {
	dmxCtrlMu.RLock()
	ctrl := dmxCtrl
	dmxCtrlMu.RUnlock()
	if ctrl == nil {
		sendError(c, "dmx_error", "DMX controller not initialized", "")
		return
	}
	channels, err := ctrl.GetAllChannels()
	if err != nil {
		sendError(c, "dmx_error", "Failed to get channels", err.Error())
		return
	}
	var states []ChannelState
	for i := 0; i < dmx.DMXChannels; i++ {
		if channels[i] > 0 {
			states = append(states, ChannelState{Address: i + 1, Value: channels[i]})
		}
	}
	presetMu.RLock()
	ap := activePresetID
	presetMu.RUnlock()
	showMu.Lock()
	var as string
	var step int
	var loop bool
	if currentShow != nil {
		as = currentShow.id
		step = currentShow.currentStep
		loop = currentShow.loop
	}
	showMu.Unlock()
	state := DMXState{Channels: states, ActivePresetID: ap, ActiveShowID: as, ShowStep: step, ShowLoop: loop, Timestamp: time.Now().UnixMilli()}
	writeJSON(c, Message{Type: "dmx_state", Payload: mustMarshal(state)})
}

func handleGetProjectConfig(c *websocket.Conn) {
	if projectStore == nil {
		sendError(c, "config_error", "Project store not initialized", "")
		return
	}
	project := projectStore.Get()
	if project == nil {
		sendError(c, "config_error", "No project loaded", "")
		return
	}
	config := map[string]interface{}{
		"project_id":   project.ID,
		"project_name": project.Name,
		"fixtures":     project.Fixtures,
		"presets":      project.Presets,
		"shows":        project.Shows,
	}
	writeJSON(c, Message{Type: "project_config", Payload: mustMarshal(config)})
}

func handleStartMonitoring(c *websocket.Conn) {
	startMonitoring()
	writeJSON(c, Message{Type: "monitoring_started", Payload: json.RawMessage("{}")})
}

func handleStopMonitoring(c *websocket.Conn) {
	stopMonitoring()
	writeJSON(c, Message{Type: "monitoring_stopped", Payload: json.RawMessage("{}")})
}
