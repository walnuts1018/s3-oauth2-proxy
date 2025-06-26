package config

type AppConfig struct {
	Port          string `env:"PORT" envDefault:"8080"`
	SessionSecret string `env:"SESSION_SECRET,required"`
}
