package cf

import "github.com/tscolari/cfapi/uaa"

type RefresherClient struct {
	Client
	cfEndpoint     string
	tokens         uaa.Tokens
	uaaClient      uaa.UAAClient
	OnTokenRefresh func(newTokens uaa.Tokens)
}

func NewRefresherClient(cfEndpoint string, tokens uaa.Tokens, uaaClient uaa.UAAClient) *RefresherClient {
	cfClient := *NewClient(cfEndpoint, tokens.AccessToken)
	return &RefresherClient{
		Client:     cfClient,
		cfEndpoint: cfEndpoint,
		tokens:     tokens,
		uaaClient:  uaaClient,
	}
}

func (c *RefresherClient) Get(path string, response interface{}) error {
	return c.fetch("GET", path, nil, response)
}

func (c *RefresherClient) Put(path string, options map[string]string, response interface{}) error {
	return c.fetch("PUT", path, options, response)
}

func (c *RefresherClient) Post(path string, options map[string]string, response interface{}) error {
	return c.fetch("POST", path, options, response)
}

func (c *RefresherClient) Delete(path string, options map[string]string) error {
	return c.fetch("DELETE", path, options, nil)
}

func (c *RefresherClient) fetch(method, path string, options map[string]string, response interface{}) error {
	err := c.Client.fetch(method, path, options, response)
	if err != nil && err.Error() == "Unauthorized" {
		err = c.refreshTokens()
		if err != nil {
			return err
		}
		return c.Client.fetch(method, path, options, response)
	}
	return err
}

func (c *RefresherClient) refreshTokens() error {
	tokens, err := c.uaaClient.RefreshToken(c.tokens.RefreshToken)
	if err != nil {
		return err
	}

	c.tokens = *tokens
	c.Client = *NewClient(c.cfEndpoint, tokens.AccessToken)

	if c.OnTokenRefresh != nil {
		c.OnTokenRefresh(c.tokens)
	}

	return nil
}
