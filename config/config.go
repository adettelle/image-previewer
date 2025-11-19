package config

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

const (
	defaultHost          = "localhost"
	defaultPort          = "8080"
	defaultCacheCapacity = "5"
	defaultResize        = "scale"
)

type Config struct {
	Logger        *LoggerConf
	Context       *context.Context
	Host          string
	Port          string
	CacheCapacity string
	Resize        string
}

type LoggerConf struct {
	Level string
}

// check if there are environment vars and fill in the config structure from there;
// next, empty fields are filled in by default.
func New(ctx *context.Context) *Config {
	cfg := Config{
		Logger: &LoggerConf{
			Level: "INFO",
		},
		Context:       ctx,
		Host:          getEnvOrDefault("HOST", defaultHost),
		Port:          getEnvOrDefault("PORT", defaultPort),
		CacheCapacity: getEnvOrDefault("CAP", defaultCacheCapacity),
		Resize:        getEnvOrDefault("RESIZE", defaultResize),
	}

	ensureAddrIsCorrect(cfg.Host, cfg.Port)

	return &cfg
}

func getEnvOrDefault(envName string, defaultVal string) string {
	res := os.Getenv(envName)
	if res != "" {
		return res
	}
	return defaultVal
}

func ensureAddrIsCorrect(host, port string) {
	hostPort := fmt.Sprintf("%s:%s", host, port)
	_, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		log.Fatal(err)
	}
	_, err = strconv.Atoi(port)
	if err != nil {
		log.Fatal(fmt.Errorf("invalid port: '%s'", port))
	}
}
