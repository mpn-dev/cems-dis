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

func (p *CemsMpnProtocol) cemsMpnLogin(task *model.Transmission) (*tokens.LoginToken, error) {
	url := fmt.Sprintf("%s/api/v1/users/login", task.BaseURL)
	var body []byte
	payload := map[string]string{
		"user_id": 	task.Username, 
		"password":	task.Password, 
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

	raw_body, _ := io.ReadAll(res.Body)
	res_body := string(raw_body)
	loginResp := &LoginResponse{}
	if err = json.Unmarshal(raw_body, loginResp); err != nil {
		log.Warningf("cems_mpn.Login => Error unmarshalling response body: %s", err.Error())
	}
	if loginResp.Error != nil {
		log.Warningf("cems_mpn.Login => Error, code: %d, body: %s", res.StatusCode, *loginResp.Error)
		return nil, errors.New(*loginResp.Error)
	}
	if res.StatusCode != 200 {
		log.Warningf("cems_mpn.Login => Error, code: %d, body: %s", res.StatusCode, res_body)
		return nil, errors.New(fmt.Sprintf("Error %d: %s", res.StatusCode, res_body))
	}
	if loginResp.Login == nil {
		return nil, errors.New("Failed parsing login information")
	}
	token := tokens.RegisterToken(task.Protocol, task.BaseURL, loginResp.Login.Token, "", loginResp.Login.TokenTTL - 5)
	return token, nil
}

func (p *CemsMpnProtocol) cemsMpnGetToken(task *model.Transmission) (*tokens.LoginToken, error) {
	token := tokens.FindToken(task.Protocol, task.BaseURL)
	if token == nil {
		return p.cemsMpnLogin(task)
	}
	if token.IsExpired() {
		return p.cemsMpnLogin(task)
	}
	return token, nil
}

func (p *CemsMpnProtocol) Send(task *model.Transmission) Result {
	record, _ := p.model.GetRawDataById(task.RawDataId)
	if record == nil {
		return Error(task, 0, "Invalid raw data ID")
	}

	token, err := p.cemsMpnGetToken(task)
	if err != nil {
		msg := fmt.Sprintf("Get token failed: %s", err.Error())
		log.Warningf("cems_mpn.Send => %s", msg)
		return Error(task, 0, msg)
	}

	url := fmt.Sprintf("%s/api/v1/cems/push", task.BaseURL)
	var body []byte
	body, _ = json.Marshal(record.CemsPayload())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Warningf("cems_mpn.Send => Error: %s", err.Error())
		return Error(task, 0, err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-CSRF-TOKEN", token.GetAccessToken())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Warningf("cems_mpn.Send => Error: %s", err.Error())
		return Error(task, 0, err.Error())
	}
	defer res.Body.Close()
	raw_body, _ := io.ReadAll(res.Body)
	res_body := string(raw_body)
	if res.StatusCode != 200 {
		log.Warningf("cems_mpn.Send => Error %d, body: %s", res.StatusCode, res_body)
		return Error(task, res.StatusCode, fmt.Sprintf("Error: `%s`", res_body))
	}
	return Success(task, 200, res_body)
}

func NewCemsMpnProtocol(model *model.Model) *CemsMpnProtocol {
	return &CemsMpnProtocol{
		model: model, 
	}
}
