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
    log.Printf("[MIGRATION] Creating table '%s'", tableName)
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

func (m migration) initRelayStations() {
  m.createTable("relay_stations", &RelayStation{})
}

func (m migration) initDeviceTokens() {
  m.createTable("device_tokens", &DeviceToken{})
}

func (m migration) initTransmissions() {
  m.createTable("transmissions", &TransmissionTable{})
}

func (m migration) initRawData() {
  m.createTable("raw_data", &RawData{})
}

func (m migration) initPushRequest() {
  m.createTable("push_requests", &PushRequest{})
}

func (m migration) Run() {
  m.initDevices()
  m.initRelayStations()
  m.initDeviceTokens()
  m.initTransmissions()
  m.initRawData()
  m.initPushRequest()
}

func newMigration(db *gorm.DB) migration {
  return migration{db: db}
}
