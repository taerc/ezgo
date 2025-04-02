package focalboard

import (
	"github.com/imroc/req/v3"
)

type TokenService struct {
	client *Client
}

type AccessToken struct {
	Token string `json:"token"`
}

func (s *TokenService) GetAccessToken() (*AccessToken , *req.Response, error) {

	var token AccessToken
	resp, err := s.client.client.R().
		SetBody(Login{
			Type:     "normal",
			UserName: s.client.username,
			Password: s.client.password,
		}).
		SetSuccessResult(&token).
		Post(s.client.RequestURL("/login"))
	return &token, resp, err 
}