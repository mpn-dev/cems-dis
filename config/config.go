package config

import (
  "path/filepath"

  "github.com/spf13/viper"
)

type config struct {
  serverPort              int
  database                *dbConfig
  userLoginTokenAgeShort  int
  userLoginTokenAgeLong   int
  deviceLoginTokenAge     int
  deviceRefreshTokenAge   int
  isDBLoggerEnabled       bool
}

var appConfig *config

func Load(file string) {
  filename := filepath.Base(file)
  ext := filepath.Ext(file)
  filenameWithoutExt := filename[0:len(filename) - len(ext)]
  viper.SetConfigName(filenameWithoutExt)
  viper.SetConfigType("yaml")
  viper.AddConfigPath(".")
  viper.SetEnvPrefix("CEMS")
  viper.AutomaticEnv()
  _ = viper.ReadInConfig()

  appConfig = &config{
    serverPort:               viper.GetInt("SERVER_PORT"), 
    database: 	              newDbConfig(), 
    userLoginTokenAgeShort:   getIntOrPanic("USER_LOGIN_TOKEN_AGE_SHORT"), 
    userLoginTokenAgeLong:    getIntOrPanic("USER_LOGIN_TOKEN_AGE_LONG"), 
    deviceLoginTokenAge:      getIntOrPanic("DEVICE_LOGIN_TOKEN_AGE"), 
    deviceRefreshTokenAge:    getIntOrPanic("DEVICE_REFRESH_TOKEN_AGE"), 
    isDBLoggerEnabled:        getBoolOrPanic("ENABLE_DB_LOGGER"), 
  }
}

func ServerPort() int {
  return appConfig.serverPort
}

func DbConfig() *dbConfig {
  return appConfig.database;
}

func UserLoginTokenAgeShort() int {
  return appConfig.userLoginTokenAgeShort
}

func UserLoginTokenAgeLong() int {
  return appConfig.userLoginTokenAgeLong
}

func DeviceLoginTokenAge() int {
  return appConfig.deviceLoginTokenAge
}

func DeviceRefreshTokenAge() int {
  return appConfig.deviceRefreshTokenAge
}

func IsDBLoggerEnabled() bool {
  return appConfig.isDBLoggerEnabled
}
