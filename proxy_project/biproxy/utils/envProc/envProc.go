package envProc

import (
        "os"
        log "github.com/sirupsen/logrus"
        "strconv"
)

func GetEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}

func GetEnvAsInt(key string, defaultValue int) int {
    valueStr, exists := os.LookupEnv(key)
    if !exists {
        return defaultValue
    }

    value, err := strconv.Atoi(valueStr)
    if err != nil {
        log.Warnf("Error parsing %s as int, using default value: %d. Error: %v", key, defaultValue, err)
        return defaultValue
    }

    return value
}
