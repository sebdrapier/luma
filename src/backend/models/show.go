package models

type Show struct {
	ID    string     `yaml:"id" json:"id"`
	Name  string     `yaml:"name" json:"name"`
	Steps []ShowStep `yaml:"steps" json:"steps"`
}

type ShowStep struct {
	PresetID string `yaml:"preset_id" json:"preset_id"`
	Duration int    `yaml:"duration" json:"duration"`
	FadeMS   int    `yaml:"fade_ms" json:"fade_ms"`
}
