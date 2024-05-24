package model

import (
	"errors"
	"time"
	log "github.com/sirupsen/logrus"
  "gorm.io/gorm"
  "cems-dis/config"
)


type DeviceToken struct {
  Id                  uint64        `gorm:"primaryKey"`
  DEV                 string        `gorm:"column:uid;size:32;index"`
  Device              Device        `gorm:"foreignKey:DEV;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
  LoginToken          string        `gorm:"index;size:64"`
  RefreshToken        string        `gorm:"index;size:64"`
  LoginExpiredAt      time.Time     `gorm"index"`
  RefreshExpiredAt    time.Time     `gorm"index"`
  CreatedAt           time.Time     `gorm:"autoCreateTime"`
}

type DeviceLogin struct {
  ApiKey          string            `json:"api_key"`
  Secret          string            `json:"secret"`
}


func (m *Model) GetDeviceToken(token string) (*DeviceToken, error) {
	deviceToken := &DeviceToken{}
	err := m.DB.Where("login_token = ?", token).Order("id DESC").First(deviceToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Warningf("DB error: %s", err.Error())
		return nil, errors.New("DB error")
	}
	return deviceToken, nil
}

func (m *Model) CreateDeviceLoginToken(uid string) (*DeviceToken, error) {
  m.DB.Where("(login_expired_at + interval '1 minute') < current_timestamp").Delete(&DeviceToken{})
  m.DB.Model(&DeviceToken{}).Where("(uid = ?) AND ((refresh_expired_at = ?) OR (refresh_expired_at > ?))", uid, nil, time.Now()).
    Update("refresh_expired_at", time.Now())
  token := &DeviceToken{
    DEV:                uid, 
    LoginToken:         GenerateDeviceLoginToken(uid), 
    RefreshToken:       GenerateDeviceLoginToken(uid), 
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
