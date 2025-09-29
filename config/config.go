package config

type Config struct {
	PgHost     string `env:"DB_HOST" default:"localhost"`
	PgUser     string `env:"DB_USER" default:"salam"`
	PgPassword string `env:"DB_PASSWORD" default:"salam"`
	PgPort     int    `env:"DB_PORT" default:"5432"`
	Db         string `env:"DB_NAME" default:"salam"`
}
