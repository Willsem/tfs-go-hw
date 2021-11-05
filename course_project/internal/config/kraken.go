package config

type Kraken struct {
	ApiKey    string `toml:"api_key"`
	RestUrl   string `toml:"rest_url"`
	SocketUrl string `toml:"socket_url"`
}
