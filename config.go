package main

type Config struct {
	sortField string
}

var DefaultConfig = NewDefaultConfig()

func NewDefaultConfig() Config {
	return Config{
		sortField: "id",
	}
}
