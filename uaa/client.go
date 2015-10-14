package uaa

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	endpoint string
}

func NewClient(endpoint string) Client {
	client := Client{
		endpoint: endpoint,
	}
	return client
}

func (c *Client) Authenticate(username, password string) (*Tokens, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("scope", "")
	data.Set("username", username)
	data.Set("password", password)

	return c.fetchToken(data)
}

func (c *Client) RefreshToken(refreshToken string) (*Tokens, error) {
	data := url.Values{
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
		"scope":         {""},
	}

	return c.fetchToken(data)
}

func (c *Client) fetchToken(data url.Values) (*Tokens, error) {
	path := fmt.Sprintf("%s/oauth/token", c.endpoint)
	request, err := http.NewRequest("POST", path, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	respBytes, err := c.runRequest(request)
	if err != nil {
		return nil, err
	}

	uaaResp := new(authenticationResponse)
	err = json.Unmarshal(respBytes, &uaaResp)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse response (%s)", err.Error())
	}

	if uaaResp.ErrorCode != "" {
		return nil, fmt.Errorf("UAA Error: %s (%s)", uaaResp.ErrorDescription, uaaResp.ErrorCode)
	}

	return &Tokens{
		AccessToken:  uaaResp.AccessToken,
		RefreshToken: uaaResp.RefreshToken,
		TokenType:    uaaResp.TokenType,
	}, nil
}

func (c *Client) runRequest(request *http.Request) ([]byte, error) {
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("cf:")))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}
