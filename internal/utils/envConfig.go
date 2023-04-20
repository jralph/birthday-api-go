package utils

import (
	"log"
	"os"
	"strconv"
	"time"
)

func OverrideFromEnvStr(str *string, envName string) {
	if envValue := os.Getenv(envName); envValue != "" {
		*str = envValue
	}
}

func OverrideFromEnvInt(i *int, envName string) {
	if envValue := os.Getenv(envName); envValue != "" {
		envValueInt, err := strconv.Atoi(envValue)
		if err != nil {
			log.Fatalf("error parsing %s env variable: %s", envName, err)
		}
		*i = envValueInt
	}
}

func OverrideFromEnvDuration(duration *time.Duration, envName string) {
	if envValue := os.Getenv(envName); envValue != "" {
		parsedDuration, err := time.ParseDuration(envValue)
		if err != nil {
			log.Fatalf("error parsing %s env variable: %s", envName, err)
		}
		*duration = parsedDuration
	}
}
