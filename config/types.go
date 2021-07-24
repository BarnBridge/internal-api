package config

type database struct {
	Host     string
	Port     string
	SSLMode  string
	DBName   string
	User     string
	Password string

	ConnectionString string `mapstructure:"connection-string"`
}

type metrics struct {
	Port int64
}

type api struct {
	Port        string
	DevCors     bool   `mapstructure:"dev-cors"`
	DevCorsHost string `mapstructure:"dev-cors-host"`
}

type addresses struct {
	Bond             string
	ExcludeTransfers []string `mapstructure:"exclude-transfers"`
}
