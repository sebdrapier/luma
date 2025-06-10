package models

type Fixture struct {
	ID          string           `yaml:"id" json:"id"`
	Name        string           `yaml:"name" json:"name"`
	Description string           `yaml:"description" json:"description"`
	Type        string           `yaml:"type" json:"type"`
	Channels    []FixtureChannel `yaml:"channels" json:"channels"`
}

type FixtureChannel struct {
	Name           string `yaml:"name" json:"name"`
	Description    string `yaml:"description" json:"description"`
	Min            int    `yaml:"min" json:"min"`
	Max            int    `yaml:"max" json:"max"`
	ChannelAddress int    `yaml:"channel_address" json:"channel_address"`
}
