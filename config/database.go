package config

import "fmt"

type dbConfig struct {
  host        string
  port        int
  user        string
  password    string
  name        string
}

func (c dbConfig) String() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s", 
		c.user, c.password, c.host, c.port, c.name)
}

func newDbConfig() *dbConfig {
  return &dbConfig{
    host:  		getStringOrPanic("DB_HOST"), 
    port:  		getIntOrPanic("DB_PORT"), 
    user:  		getStringOrPanic("DB_USER"), 
    password:	getStringOrPanic("DB_PASS"), 
    name:  		getStringOrPanic("DB_NAME"), 
  }
}
