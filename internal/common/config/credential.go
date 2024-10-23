package config

type Credential struct {
	URL      string
	Username string
	Password string
}

func (c *Config) NewCredential() *Credential {
	// can improve the credential logic within login_page
	return &Credential{}
}
