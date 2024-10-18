package config

type Config struct {
	Hour       string
	Category   string
	Filter     bool
	Max        int
	Adjustment bool
}

func Set(hour, category string, filter bool, max int, adjustment bool) *Config {
	return &Config{
		Hour:       hour,
		Category:   category,
		Filter:     filter,
		Max:        max,
		Adjustment: adjustment,
	}
}
