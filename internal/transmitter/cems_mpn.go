package transmitter

import (
  "bytes"
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  tokens "cems-dis/internal/token_store"
  "cems-dis/model"
)

type LoginData struct {
  Token             string      `json:"token"`
  TokenTTL          int64       `json:"token_ttl"`
  UserId            int         `json:"user_id"`
  UserName          string      `json:"user_name"`
  IsAdmin           bool        `json:"is_admin"`
  IdInstansi        int         `json:"id_instansi"`
  AppTitle          string      `json:"app_title"`
  LogoLeft          string      `json:"logo_left"`
  HeadBgColor       string      `json:"head_bgcolor"`
  HasAccountMenu    bool        `json:"has_account_menu"`
}

type LoginResponse struct {
  Success           bool        `json:"success"`
  Error             *string     `json:"error"`
  Login             *LoginData  `json:"data"`
}

type CemsMpnProtocol struct {
	model *model.Model
}

func (p *CemsMpnProtocol) cemsMpnLogin(station model.RelayStation) (*tokens.LoginToken, error) {
	url := fmt.Sprintf("%s/api/v1/users/login", station.BaseURL)
	var body []byte
	payload := map[string]string{
		"user_id": 	station.Username, 
		"password":	station.Password, 
	}
	body, _ = json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	res_body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		msg := fmt.Sprintf("Error %d", res.StatusCode)
		log.Warningf("cems_mpn.Login => %s. Body: %s", msg, res_body)
		return nil, errors.New(msg)
	}
	loginResp := &LoginResponse{}
	err = json.Unmarshal(res_body, loginResp)
	if err != nil {
		log.Warningf("cems_mpn.Login => Error unmarshalling response body: %s", err.Error())
		return nil, err
	}
	if loginResp.Error != nil {
		log.Warningf("cems_mpn.Login => Error: %s", loginResp.Error)
		return nil, errors.New(*loginResp.Error)
	}
	token := tokens.RegisterToken(station.Protocol, station.BaseURL, loginResp.Login.Token, "", loginResp.Login.TokenTTL - 5)
	return token, nil
}

func (p *CemsMpnProtocol) cemsMpnGetToken(task model.Transmission, station model.RelayStation) (*tokens.LoginToken, error) {
	token := tokens.FindToken(station.Protocol, station.BaseURL)
	if token == nil {
		return p.cemsMpnLogin(station)
	}
	if token.IsExpired() {
		return p.cemsMpnLogin(station)
	}
	return token, nil
}

func (p *CemsMpnProtocol) Send(task model.Transmission, station model.RelayStation) {
	record, _ := p.model.GetRawDataById(task.RawDataId)
	if record == nil {
		p.model.SetTransmissionError(task, 0, "Invalid raw data ID")
		return
	}

	p.model.SetTransmissionStarted(task)
	token, err := p.cemsMpnGetToken(task, station)
	if err != nil {
		msg := fmt.Sprintf("Get token failed: %s", err.Error())
		log.Warningf("cems_mpn.Send => %s", msg)
		p.model.SetTransmissionError(task, 0, msg)
		return
	}

	url := fmt.Sprintf("%s/api/v1/cems/push", station.BaseURL)
	var body []byte
	body, _ = json.Marshal(record.CemsPayload())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Warningf("cems_mpn.Send => Error: %s", err.Error())
		p.model.SetTransmissionError(task, 0, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-CSRF-TOKEN", token.GetAccessToken())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Warningf("cems_mpn.Send => Error: %s", err.Error())
		p.model.SetTransmissionError(task, 0, err.Error())
		return
	}
	defer res.Body.Close()
	res_body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		log.Warningf("cems_mpn.Send => Error %d", res.StatusCode)
		p.model.SetTransmissionError(task, res.StatusCode, "Unknown error")
		return
	}
	p.model.SetTransmissionSuccess(task, string(res_body))
}

func NewCemsMpnProtocol(model *model.Model) *CemsMpnProtocol {
	return &CemsMpnProtocol{
		model: model, 
	}
}
