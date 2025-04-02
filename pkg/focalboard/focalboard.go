package focalboard

import (
	"net/url"
	"strings"

	"github.com/imroc/req/v3"
)

type Login struct {
	Type     string `json:"type"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

const (
	defaultBaseUrl = "http://47.122.5.188/"
	apiVersionPath = "api/v2"
	userAgent = "focalboard-go/0.1.0"
)

type Client struct {
	client *req.Client
	baseURL *url.URL
	username string
	password string
	token string

	Token *TokenService
}

func NewClient(username, password string) (*Client, error) {

	c:= &Client{
		client: req.C(),	
	}

	c.username = username
	c.password = password

	c.Token = &TokenService{
		client: c,
	}
	c.setBaseURL(defaultBaseUrl)

	c.Token = &TokenService{
		client: c,
	}

	// headers 
	c.client.SetCommonHeader("User-Agent", userAgent)
	c.client.SetCommonHeader("Content-Type", "application/json")
	c.client.SetCommonHeader("X-Requested-With", "XMLHttpRequest")
	c.setDumpAll()
	acessToken, _, err := c.Token.GetAccessToken()

	if err!= nil {
		return nil, err
	}
	c.token = acessToken.Token
	c.client.SetCommonHeader("Authorization", "Bearer " + c.token)
	return c, nil
}

func (c *Client) setBaseURL(urlStr string) error {
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(baseURL.Path, apiVersionPath) {
		baseURL.Path += apiVersionPath
	}

	// Update the base URL of the client.
	c.baseURL = baseURL

	return nil
}

func (c *Client) RequestURL(path string) string{
	u := *c.baseURL
	u.Path = c.baseURL.Path + path
	return u.String()
}

func (c *Client) setDumpAll() {
	c.client.EnableDumpAll()		
}
