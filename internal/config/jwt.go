package config

type JWTConfig struct {
	Issuer     string `yaml:"issuer" env:"ISSUER" required:"true"`
	SignIssuer string `yaml:"sign_issuer" env:"SIGN_ISSUER" required:"true"`
}
