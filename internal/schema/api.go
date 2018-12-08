package schema

// APIConfig defines external API settings.
type APIConfig struct {
	URL       string
	Name      string
	Secret    string
	Client    interface{}
	Endpoints APIEndpoints
}

// APIEndpoints defined extrnal API endpoints.
type APIEndpoints struct {
	Authenticate string
	GetDevices   string `toml:"get_devices"`
	UpdateDevice string `toml:"update_device"`
}
