package configuration

type Logger struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type Config struct {
	DatabaseDsn string  `json:"database_dsn"`
	Logger      *Logger `json:"logger"`
}
