package utils

import (
	"log"
	"os"
	"strconv"
	"time"
)

// OverrideFromEnvStr will replace the given string pointer with the given env value
func OverrideFromEnvStr(str *string, envName string) {
	if envValue := os.Getenv(envName); envValue != "" {
		*str = envValue
	}
}

// OverrideFromEnvInt will replace the given int pointer with the given env value, converted to an int
func OverrideFromEnvInt(i *int, envName string) {
	if envValue := os.Getenv(envName); envValue != "" {
		envValueInt, err := strconv.Atoi(envValue)
		if err != nil {
			log.Fatalf("error parsing %s env variable: %s", envName, err)
		}
		*i = envValueInt
	}
}

// OverrideFromEnvDuration will replace the given duration pointer with the given env value, parsed as a duration
func OverrideFromEnvDuration(duration *time.Duration, envName string) {
	if envValue := os.Getenv(envName); envValue != "" {
		parsedDuration, err := time.ParseDuration(envValue)
		if err != nil {
			log.Fatalf("error parsing %s env variable: %s", envName, err)
		}
		*duration = parsedDuration
	}
}
