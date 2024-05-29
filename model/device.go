package model

import (
  "errors"
  "fmt"
  "strings"
  "time"
  log "github.com/sirupsen/logrus"
  "gorm.io/gorm"
  "cems-dis/utils"
)

const MAX_DEVICE_UID_LENGTH = 32

type Device struct {
  UID         string                `json:"uid"         gorm:"size:32;primaryKey"`
  Name        string                `json:"name"        gorm:"size:50"`
  Latitude    *float64              `json:"latitude"`
  Longitude   *float64              `json:"longitude"`
  ApiKey      string                `json:"apikey"      gorm:"index;size:30"`
  Secret      string                `json:"secret"      gorm:"size:64"`
  Enabled     bool                  `json:"enabled"     gorm:"index"`
  CreatedAt   time.Time             `json:"created_at"  gorm:"autoCreateTime"`
  UpdatedAt   time.Time             `json:"updated_at"  gorm:"autoUpdateTime"`
}

type DeviceOut struct {
  Device
  CreatedAt   string                `json:"created_at"`
  UpdatedAt   string                `json:"updated_at"`
}


func (d *Device) Copy() *Device {
  return &Device{
    UID:        d.UID, 
    Name:       d.Name, 
    Latitude:   d.Latitude, 
    Longitude:  d.Longitude, 
    ApiKey:     d.ApiKey, 
    Secret:     d.Secret, 
    Enabled:    d.Enabled, 
    CreatedAt:  d.CreatedAt, 
    UpdatedAt:  d.UpdatedAt, 
  }
}

func (d *Device) Trim() {
  d.UID = strings.Trim(d.UID, " ")
  d.Name = strings.Trim(d.Name, " ")
  d.ApiKey = strings.Trim(d.ApiKey, " ")
  d.Secret = strings.Trim(d.Secret, " ")
}

func (o *Device) Update(n *Device) {
  o.UID         = n.UID
  o.Name        = n.Name
  o.Latitude    = n.Latitude
  o.Longitude   = n.Longitude
  o.ApiKey      = n.ApiKey
  o.Secret      = n.Secret
  o.Enabled     = n.Enabled
}

func (d *Device) Out() *DeviceOut {
  return &DeviceOut{
    Device:     *d, 
    CreatedAt:  utils.TimeToString(d.CreatedAt), 
    UpdatedAt:  utils.TimeToString(d.UpdatedAt), 
  }
}


func (m *Model) GetDeviceByUid(uid string) (*Device, error) {
  device := &Device{}
  if err := m.DB.Where("uid = ?", uid).First(device).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      return nil, nil
    }
    log.Warningf(fmt.Sprintf("DB error: %s", err.Error()))
    return nil, errors.New("DB error")
  }
  return device, nil
}

func (m *Model) IsDeviceExist(uid string) bool {
  device, _ := m.GetDeviceByUid(uid)
  return device != nil
}

func (m *Model) IsDeviceUidTaken(uid string, uidNot string) bool {
  device := &Device{}
  err := m.DB.Where("(uid = ?) AND (uid <> ?)", uid, uidNot).First(device).Error
  if (err != nil) && (errors.Is(err, gorm.ErrRecordNotFound)) {
    return false
  }
  return true
}

func (m *Model) IsDeviceApiKeyTaken(apiKey string, uidNot string) bool {
  device := &Device{}
  err := m.DB.Where("(api_key = ?) AND (uid <> ?)", apiKey, uidNot).First(device).Error
  if (err != nil) && errors.Is(err, gorm.ErrRecordNotFound) {
    return false
  }
  return true
}

func (m *Model) DeleteDeviceByUid(uid string) error {
  err := m.DB.Delete(&Device{UID: uid}).Error
  if err != nil {
    log.Warningf("DB error: %s", err.Error())
    return errors.New("DB error")
  }
  return nil
}
