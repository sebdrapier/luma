package models

type Project struct {
	ID           string    `yaml:"id" json:"id"`
	Name         string    `yaml:"name" json:"name"`
	USBInterface string    `yaml:"usb_interface" json:"usb_interface"`
	Fixtures     []Fixture `yaml:"fixtures" json:"fixtures"`
	Presets      []Preset  `yaml:"presets" json:"presets"`
	Shows        []Show    `yaml:"shows" json:"shows"`
}
