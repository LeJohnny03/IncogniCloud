package config

type Config struct {
	CookieExpiration int `env:"COOKIE_EXPIRATION" envDefault:"7200"`
}
