package main

import (
  "fmt"
  "os"
)


const (
  DEFAULT_PORT = "8080"
  DEFAULT_ADDRESS = ""
)


func getenvOrElse(key, def string) string {
  v, ok := os.LookupEnv(key)
  if !ok {
    return def
  }
  return v
}

type Config struct {
  Port string
  Address string
}

func (c Config) AddressString() string {
  return fmt.Sprintf("%s:%s", c.Address, c.Port)
}

func NewEnvConfig() Config {
  port := getenvOrElse("PORT", DEFAULT_PORT)
  address := getenvOrElse("ADDRESS", DEFAULT_ADDRESS)

  return Config{
    Port: port,
    Address: address,
  }
}
