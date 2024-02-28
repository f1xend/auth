package config

import (
	"github.com/joho/godotenv"
	"log"
)

type GRPCConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type HTTPConfig interface {
	Address() string
}

func Load(path string) error {
	if err := godotenv.Load(path); err != nil {
		log.Fatalf("failed to load config %v", err)
		return err
	}
	return nil
}
