package token_store

import (
	"time"
)

type LoginToken struct {
	protocol			string
	baseUrl 			string
	accessToken		string
	refreshToken	string
	lastLogin			time.Time
	expiresSec		int64
}

var tokens []*LoginToken

func findToken(protocol string, baseUrl string) (int, *LoginToken) {
	for i, t := range tokens {
		if (t.protocol == protocol) && (t.baseUrl == baseUrl) {
			return i, t
		}
	}
	return -1, nil
}

func (t *LoginToken) SetAccessToken(token string) {
	t.accessToken = token
}

func (t *LoginToken) GetAccessToken() string {
	return t.accessToken
}

func (t *LoginToken) SetRefreshToken(token string) {
	t.refreshToken = token
}

func (t *LoginToken) GetRefreshToken() string {
	return t.refreshToken
}

func (t *LoginToken) RenewLoginTime() {
	t.lastLogin = time.Now()
}

func (t *LoginToken) IsExpired() bool {
	return t.lastLogin.Add(time.Duration(t.expiresSec) * time.Second).Before(time.Now())
}

func RegisterToken(protocol string, baseUrl string, accessToken string, refreshToken string, expiresSec int64) *LoginToken {
	token_ttl := expiresSec
	if token_ttl < 0 {
		token_ttl = 0
	}

	_, tok := findToken(protocol, baseUrl)
	if tok != nil {
		tok.accessToken 	= accessToken
		tok.refreshToken	= refreshToken
		tok.expiresSec 		= token_ttl
		tok.lastLogin 		= time.Now()
		return tok
	}

	token := &LoginToken{
		protocol:				protocol, 
		baseUrl: 				baseUrl, 
		accessToken:		accessToken, 
		refreshToken:		refreshToken, 
		lastLogin:			time.Now(), 
		expiresSec:			token_ttl, 
	}

	tokens = append(tokens, token)
	return token
}

func RemoveToken(protocol string, baseUrl string) {
	idx, tok := findToken(protocol, baseUrl)
	if tok != nil {
		newTokens := []*LoginToken{}
		for i, t := range(tokens) {
			if i != idx {
				newTokens = append(newTokens, t)
			}
		}
		tokens = newTokens
	}
}

func FindToken(protocol string, baseUrl string) *LoginToken {
	_, t := findToken(protocol, baseUrl)
	return t
}

func Init() {
	if tokens == nil {
		tokens = []*LoginToken{}
	}
}
