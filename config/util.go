package config

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/spf13/viper"
)

func getStringOrPanic(key string) string {
	value := viper.GetString(key)
	if value == "" {
		panic(fmt.Sprintf("Config value for '%s' is not set", key))
	}
	return value
}

func getIntOrPanic(key string) int {
	str := getStringOrPanic(key)
	value, err := strconv.Atoi(str)
	if err != nil {
		panic(fmt.Sprintf("Config value for '%s' is not a valid int: '%s'", key, str))
	}
	return value
}

func getInt64OrPanic(key string) int64 {
	str := getStringOrPanic(key)
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Config value for '%s' is not a valid int64: '%s'", key, str))
	}
	return value
}

func getBoolOrPanic(key string) bool {
	str := strings.ToLower(getStringOrPanic(key))
	if str != "true" && str != "false" {
		panic(fmt.Sprintf("Config value for '%s' is not a valid bool: '%s'", key, str))
	}
	return str == "true"
}
