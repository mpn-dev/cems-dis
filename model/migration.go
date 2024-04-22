package model

import (
  "io/ioutil"

  "gorm.io/gorm"

  log "github.com/sirupsen/logrus"
)


type migration struct {
  db    *gorm.DB
}


func (m migration) createTable(tableName string, blueprint interface{}) {
  if !m.db.Migrator().HasTable(tableName) {
    log.Printf("[MIGRATION] Creating table '%s'\n", tableName)
    m.db.AutoMigrate(blueprint)
  }
}

func (m migration) execSqlFile(file string) {
  b, e := ioutil.ReadFile(file)
  if(e != nil) {
    log.Println(e.Error())
    panic(e.Error())
  }

  m.db.Exec(string(b))
}

func (m migration) initDevices() {
  m.createTable("devices", &Device{})
}

func (m migration) initDeviceTokens() {
  m.createTable("device_tokens", &DeviceToken{})
}

func (m migration) initRawData() {
  m.createTable("cems_records", &RawData{})
}

func (m migration) Run() {
  m.initDevices()
  m.initDeviceTokens()
  m.initRawData()
}

func newMigration(db *gorm.DB) migration {
  return migration{db: db}
}
