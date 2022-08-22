/*Package structs define models
 */
package types

// Token base Info
// swagger:model Token
type Token struct {
	Value string `json:"value"`
}

// Redirect  struct save request data
type Redirect struct {
	URL string
}

// AccessToken oauth2 access token
type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

// Oauth2ServerConfig oauth2 服务配置参数
type Oauth2ServerConfig struct {
	ClientID    string `json:"client_id" validate:"required"`    // appID
	Secret      string `json:"secret" validate:"required"`       // app Secret
	CallbackURL string `json:"callback_url" validate:"required"` // call url
	OauthServer string `json:"oauth_server" validate:"required"` // oauth server url
	OauthAPI    string `json:"oauth_api" validate:"required"`
	TokenAPI    string `json:"token_api" validate:"required"`
	UserAPI     string `json:"user_api" validate:"required"`
	ExternalID  string `json:"external_id" validate:"required"`
	UserName    string `json:"user_name" validate:"required"`
}
