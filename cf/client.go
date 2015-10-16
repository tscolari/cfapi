package cf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tscolari/cfapi/uaa"
)

type Client struct {
	accessToken string
	endpoint    string
	client      *http.Client
}

func NewClient(endpoint, accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		endpoint:    endpoint,
		client:      &http.Client{},
	}
}

func (c *Client) Get(path string, response interface{}) error {
	return c.fetch("GET", path, nil, response)
}

func (c *Client) Put(path string, options map[string]string, response interface{}) error {
	return c.fetch("PUT", path, options, response)
}

func (c *Client) Post(path string, options map[string]string, response interface{}) error {
	return c.fetch("POST", path, options, response)
}

func (c *Client) Delete(path string, options map[string]string) error {
	return c.fetch("DELETE", path, options, nil)
}

func (c *Client) CurrentTokens() uaa.Tokens {
	return uaa.Tokens{
		AccessToken: c.accessToken,
	}
}

func (c *Client) fetch(method, path string, options map[string]string, response interface{}) error {
	req, err := c.createRequest(method, path, options)
	if err != nil {
		return err
	}
	resp, err := c.executeRequest(req)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, response)
}

func (c *Client) createRequest(method, path string, body map[string]string) (*http.Request, error) {
	var requestBody io.Reader

	if body != nil {
		json, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("Invalid options format: %s", err.Error())
		}
		requestBody = strings.NewReader(string(json))
	}

	req, err := http.NewRequest(method, c.endpoint+path, requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")
	return req, err
}

func (c *Client) executeRequest(request *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect: %s", err.Error())
	}

	return resp, nil
}

func (c *Client) parseResponse(resp *http.Response, returnObj interface{}) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 500 {
		var errResp Error
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return fmt.Errorf("%s: %s", http.StatusText(resp.StatusCode), body)
		}

		return errors.New(strings.TrimSpace(errResp.Description))
	}

	if resp.StatusCode >= 400 {
		return errors.New(http.StatusText(resp.StatusCode))
	}

	if returnObj == nil {
		return nil
	}

	err = json.Unmarshal(body, returnObj)
	if err != nil {
		return fmt.Errorf("Failed to parse response: %s", err.Error())
	}

	return nil
}
