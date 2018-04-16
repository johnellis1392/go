package main

import (
	"fmt"
	"os"
)

const (
	DEFAULT_PORT = "8080"
	DEFAULT_ADDR = ""
)

type config struct {
	port string
	addr string
}

func (c config) AddressString() string {
	return fmt.Sprintf("%s:%s", c.addr, c.port)
}

func getenvOrDefault(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func envConfig() config {
	c := config{
		port: getenvOrDefault("PORT", DEFAULT_PORT),
		addr: getenvOrDefault("ADDR", DEFAULT_ADDR),
	}
	return c
}
