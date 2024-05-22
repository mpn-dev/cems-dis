package model

import (
  "crypto/sha256"
  "encoding/hex"
  "fmt"
  "strings"
  "time"

  log "github.com/sirupsen/logrus"
  "golang.org/x/crypto/bcrypt"
  "gorm.io/gorm"
  "gorm.io/gorm/logger"
  "gorm.io/driver/postgres"

  "cems-dis/config"
  "cems-dis/utils"
)


type Model struct {
  DB *gorm.DB
}

var model *Model


func GeneratePasswordHash(password string) (string, error) {
  hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
  if err != nil {
    return "", err
  }
  return string(hash), nil
}

func IsUserPasswordMatch(hash string, password string) bool {
  return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func GenerateUserLoginToken(userId int64, userAgent string) string {
  bytes := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%s|%d", utils.RandomString32(), userId, userAgent, time.Now().UnixMicro())))
  return hex.EncodeToString(bytes[:])
}

func GenerateDeviceLoginToken(uid string) string {
  bytes := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%d", utils.RandomString32(), uid, time.Now().UnixMicro())))
  return hex.EncodeToString(bytes[:])
}

func SetSearchKeywords(sql *gorm.DB, fields []string, keywords string) *gorm.DB {
  q := strings.Trim(keywords, " ")
  if (len(fields) == 0) || (len(q) == 0) {
    return sql
  }
  ff := strings.Join(fields, ", ")
  qq := strings.Split(q, " ")
  for _, qx := range qq {
    sql = sql.Where(fmt.Sprintf("CONCAT('|', CONCAT_WS('|', %s), '|') ILIKE ?", ff), fmt.Sprintf("%%%s%%", qx))
  }
  return sql
}


func New() (*Model, error) {
  if model == nil {
    url := config.DbConfig().String()
    logLevel := logger.Silent
    if config.IsDBLoggerEnabled() {
      logLevel = logger.Error
    }
    db, err := gorm.Open(postgres.Open(url), &gorm.Config{
      Logger: logger.Default.LogMode(logLevel), 
    })

    if err != nil {
      log.Warningf("Error connecting to database: %s", err)
      return nil, err
    }
    newMigration(db).Run()
    model = &Model{DB: db}
  }
  return model, nil
}
