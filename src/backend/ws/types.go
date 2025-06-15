package ws

import (
	"context"
	"encoding/json"

	"elano.fr/src/backend/models"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ChannelUpdatePayload struct {
	DMXAddress int `json:"dmx_address"`
	Value      int `json:"value"`
}

type PresetPayload map[string]int

type ShowStep struct {
	Preset   PresetPayload `json:"preset"`
	Duration int           `json:"duration"`
	FadeMs   int           `json:"fade_ms"`
}

type ShowPayload struct {
	Steps []ShowStep `json:"steps"`
	Loop  bool       `json:"loop,omitempty"`
}

type ShowController struct {
	cancel      context.CancelFunc
	id          string
	currentStep int
	showData    *models.Show
	loop        bool
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type ChannelState struct {
	Address int  `json:"address"`
	Value   byte `json:"value"`
}

type DMXState struct {
	Channels       []ChannelState `json:"channels"`
	ActivePresetID string         `json:"active_preset_id,omitempty"`
	ActiveShowID   string         `json:"active_show_id,omitempty"`
	ShowStep       int            `json:"show_step,omitempty"`
	ShowLoop       bool           `json:"show_loop,omitempty"`
	Timestamp      int64          `json:"timestamp"`
}
