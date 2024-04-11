package config

type Configuration struct {
	// General Settings
	Name                 string
	TelemetryUpdateDelay int `yaml:"telemetry-update-delay"`

	// MQTT Settings
	Broker   string
	Id       string
	User     string
	Password string

	// Sensors
	BootTimestamp bool `yaml:"boot-timestamp"`

	Cpu struct {
		Usage struct {
			Enabled bool
			Total bool
			Decimal int
		}
		Temperature struct {
			Enabled bool
			Decimal int
		}
	}

	Ram struct {
		Enabled bool
		Decimal int
	}

	Swap struct {
		Enabled bool
		Decimal int
	}

	Storage []Drive

	Network []Interface

	Rpi struct {
		PowerStatus bool `default:"false" yaml:"power-status"`
	}

	Advanced struct {
		StartDelay int    `default:"0" yaml:"start-delay"`
		DeviceID   string `default:"auto" yaml:"device-id"`
	}
}

type Drive struct {
	Drive   string
	Decimal int
}

type Interface struct {
	Interface string
	Bitrate   string
	Decimal   int
	Ingress   bool
	Egress    bool
}
