package models

type Preset struct {
	ID          string         `yaml:"id" json:"id"`
	Name        string         `yaml:"name" json:"name"`
	Description string         `yaml:"description" json:"description"`
	Channels    []ChannelValue `yaml:"channels" json:"channels"`
}

type ChannelValue struct {
	DMXAddress int  `yaml:"dmx_address" json:"dmx_address"`
	Value      byte `yaml:"value" json:"value"`
}
