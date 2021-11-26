package config

type Server struct {
	ListenAddr string `toml:"listen_addr"`
	CertFile   string `toml:"cert_file"`
	KeyFile    string `toml:"key_file"`
}
