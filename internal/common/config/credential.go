package config

type Credential struct {
	URL      string
	Username string
	Password string
}

func (c *Config) NewCredential() *Credential {
	return &Credential{}
}
