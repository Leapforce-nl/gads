package gads

import (
	"encoding/json"
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type AuthConfig struct {
	file         string             `json:"-"`
	OAuth2Config *oauth2.Config     `json:"oauth2.Config"`
	OAuth2Token  *oauth2.Token      `json:"oauth2.Token"`
	tokenSource  oauth2.TokenSource `json:"-"`
	Auth         Auth               `json:"gads.Auth"`
}

func NewCredentialsFromCode(ctx context.Context, customerId string, developerToken string, userAgent string, token *oauth2.Token) (ac AuthConfig, err error) {
	ac.Auth.CustomerId = customerId
	ac.Auth.DeveloperToken = developerToken
	ac.Auth.UserAgent = userAgent

	ac.Auth.Client = oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	return ac, err
}

func NewCredentialsFromFile(pathToFile string, ctx context.Context) (ac AuthConfig, err error) {
	data, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return ac, err
	}
	if err := json.Unmarshal(data, &ac); err != nil {
		return ac, err
	}
	ac.file = pathToFile
	ac.tokenSource = ac.OAuth2Config.TokenSource(ctx, ac.OAuth2Token)
	ac.Auth.Client = ac.OAuth2Config.Client(ctx, ac.OAuth2Token)
	return ac, err
}

func NewCredentials(ctx context.Context) (ac AuthConfig, err error) {
	return NewCredentialsFromFile(*configJson, ctx)
}

// Save writes the contents of AuthConfig back to the JSON file it was
// loaded from.
func (c AuthConfig) Save() error {
	if c.file == "" {
		return nil
	}

	configData, err := json.MarshalIndent(&c, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.file, configData, 0600)
}

// Token implements oauth2.TokenSource interface and store updates to
// config file.
func (c AuthConfig) Token() (token *oauth2.Token, err error) {
	// use cached token
	if c.OAuth2Token.Valid() {
		return c.OAuth2Token, err
	}

	// get new token from tokens source and store
	c.OAuth2Token, err = c.tokenSource.Token()
	if err != nil {
		return nil, err
	}
	return token, c.Save()
}
