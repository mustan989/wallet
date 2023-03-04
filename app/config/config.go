package config

type Config struct {
	Database *Database `json:"database" yaml:"database"`
	Server   *Server   `json:"server" yaml:"server"`
}

type Database struct {
	Host string `json:"host" yaml:"host" env:"DATABASE_HOST"`
	Port int    `json:"port" yaml:"port" env:"DATABASE_PORT"`
	User string `json:"user" yaml:"user" env:"DATABASE_USER"`
	Pass string `json:"pass" yaml:"pass" env:"DATABASE_PASS"`
	Name string `json:"name" yaml:"name" env:"DATABASE_NAME"`
}

type Server struct {
	Port int `json:"port" yaml:"port" env:"SERVER_PORT"`
}
