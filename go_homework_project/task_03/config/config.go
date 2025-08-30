package config

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SeverPort  string
}

func LoadConfig() (Config, error) {
	var cfg Config

	cfg.DBHost = "127.0.0.1"
	cfg.DBPort = "3306"
	cfg.DBName = "gorm_test"
	cfg.DBUser = "root"
	cfg.DBPassword = "root"

	return cfg, nil
}
