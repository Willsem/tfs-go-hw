package config

type Kraken struct {
	ApiKey    string `toml:"api_key"`
	ApiSecret string `toml:"api_secret"`
	RestUrl   string `toml:"rest_url"`
	SocketUrl string `toml:"socket_url"`
}
