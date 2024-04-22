package model

import (
  "errors"
  "fmt"
  "strings"
  "time"

  log "github.com/sirupsen/logrus"
  "gorm.io/gorm"

  "cems-dis/config"
)

type Device struct {
  Id          uint64                `json:"id"          gorm:"primaryKey"`
  UID         string                `json:"uid"         gorm:"index;size:30"`
  Name        string                `json:"name"        gorm:"size:30"`
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

type DeviceToken struct {
  Id                  uint64        `gorm:"primaryKey"`
  DeviceId            uint64        `gorm:"index"`
  LoginToken          string        `gorm:"index;size:64"`
  RefreshToken        string        `gorm:"index;size:64"`
  LoginExpiredAt      time.Time
  RefreshExpiredAt    time.Time
  CreatedAt           time.Time     `gorm:"autoCreateTime"`
}

type DeviceLogin struct {
  ApiKey          string            `json:"api_key"`
  Secret          string            `json:"secret"`
}


func (d *Device) Copy() *Device {
  return &Device{
    Id:         d.Id, 
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
    CreatedAt:  d.CreatedAt.Format(DEFAULT_DATE_TIME_FORMAT), 
    UpdatedAt:  d.UpdatedAt.Format(DEFAULT_DATE_TIME_FORMAT), 
  }
}


func (m *Model) GetDeviceById(id uint64) (*Device, error) {
  device := &Device{}
  err := m.DB.First(device, id).Error
  if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      return nil, errors.New(fmt.Sprintf("ID device '%d' tidak ada di database", id))
    } else {
      log.Warningf("DB error: %s", err.Error())
      return nil, errors.New("DB error")
    }
  }
  return device, nil
}

func (m *Model) GetDeviceByUid(uid string) (*Device, error) {
  device := &Device{}
  if err := m.DB.Where("uid = ?", uid).First(device).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      log.Warningf(fmt.Sprintf("DB error: %s", err.Error()))
      return nil, errors.New("DB error")
    } else {
      return nil, errors.New("UID tidak ada di database")
    }
  }
  return device, nil
}

func (m *Model) IsDeviceExist(id uint64) bool {
  device, _ := m.GetDeviceById(id)
  return device != nil
}

func (m *Model) IsDeviceUidTaken(uid string, idNot uint64) bool {
  device := &Device{}
  err := m.DB.Where("(uid = ?) AND (id <> ?)", uid, idNot).First(device).Error
  if (err != nil) && (errors.Is(err, gorm.ErrRecordNotFound)) {
    return false
  }
  return true
}

func (m *Model) IsDeviceApiKeyTaken(apiKey string, idNot uint64) bool {
  device := &Device{}
  err := m.DB.Where("(api_key = ?) AND (id <> ?)", apiKey, idNot).First(device).Error
  if (err != nil) && errors.Is(err, gorm.ErrRecordNotFound) {
    return false
  }
  return true
}

func (m *Model) CreateDeviceLoginToken(deviceId uint64) (*DeviceToken, error) {
  m.DB.Model(&DeviceToken{}).Where("(device_id = ?) AND ((login_expired_at = ?) OR (login_expired_at > ?))", deviceId, nil, time.Now()).
    Update("login_expired_at", time.Now())
  m.DB.Model(&DeviceToken{}).Where("(device_id = ?) AND ((refresh_expired_at = ?) OR (refresh_expired_at > ?))", deviceId, nil, time.Now()).
    Update("refresh_expired_at", time.Now())
  token := &DeviceToken{
    DeviceId:           deviceId, 
    LoginToken:         GenerateDeviceLoginToken(deviceId), 
    RefreshToken:       GenerateDeviceLoginToken(deviceId), 
    LoginExpiredAt:     time.Now().Add(time.Duration(config.DeviceLoginTokenAge()) * time.Second), 
    RefreshExpiredAt:   time.Now().Add(time.Duration(config.DeviceRefreshTokenAge()) * time.Second), 
  }
  err := m.DB.Create(token).Error
  if err != nil {
    log.Warningf("DB error: %s", err.Error())
    return nil, errors.New("DB error")
  }
  return token, nil
}

func (m *Model) DeleteDeviceById(id uint64) error {
  err := m.DB.Delete(&Device{}, id).Error
  if err != nil {
    log.Warningf("DB error: %s", err.Error())
    return errors.New("DB error")
  }
  return nil
}
