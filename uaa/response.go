package uaa

type authenticationResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	RefreshToken     string `json:"refresh_token"`
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
