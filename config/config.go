package configs

type Config struct {
	DBName     string
	DBPassword string
	DBUser     string
	DBPort     string
	DBHost     string
}

func LoadConfig() *Config {
	viper := NewViper()

	return &Config{
		DBName:     viper.GetString("MYSQL_DATABASE"),
		DBPassword: viper.GetString("MYSQL_ROOT_PASSWORD"),
		DBUser:     viper.GetString("MYSQL_USER"),
		DBPort:     viper.GetString("MYSQL_PORT"),
		DBHost:     "127.0.0.1",
	}
}
