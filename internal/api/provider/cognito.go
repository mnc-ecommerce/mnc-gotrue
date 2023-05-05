package provider

import (
	"context"
	"errors"
	"strings"

	"github.com/supabase/gotrue/internal/conf"
	"golang.org/x/oauth2"
)

// Cognito

const defaultcognitoAuthBase = "https://aladinmall.auth.ap-southeast-1.amazoncognito.com"

type cognitoProvider struct {
	*oauth2.Config
	Host string
}

type cognitoUser struct {
	Sub               string `json:"sub"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
}

// NewCognitoProvider creates a Cognito account provider.
func NewCognitoProvider(ext conf.OAuthProviderConfiguration, scopes string) (OAuthProvider, error) {
	if err := ext.Validate(); err != nil {
		return nil, err
	}

	oauthScopes := []string{
		"openid",
	}

	if scopes != "" {
		oauthScopes = append(oauthScopes, strings.Split(scopes, ",")...)
	}

	host := chooseHost(ext.URL, defaultcognitoAuthBase)
	return &cognitoProvider{
		Config: &oauth2.Config{
			ClientID:     ext.ClientID,
			ClientSecret: ext.Secret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  host + "/oauth2/authorize",
				TokenURL: host + "/oauth2/token",
			},
			RedirectURL: ext.RedirectURI,
			Scopes:      oauthScopes,
		},
		Host: host,
	}, nil
}

func (g cognitoProvider) GetOAuthToken(code string) (*oauth2.Token, error) {
	return g.Exchange(context.Background(), code)
}

func (g cognitoProvider) GetUserData(ctx context.Context, tok *oauth2.Token) (*UserProvidedData, error) {
	var u cognitoUser

	if err := makeRequest(ctx, tok, g.Config, g.Host+"/oauth2/userInfo", &u); err != nil {
		return nil, err
	}

	if u.Email == "" {
		return nil, errors.New("unable to find email with Cognito provider")
	}

	return &UserProvidedData{
		Metadata: &Claims{
			Subject:           u.Sub,
			Name:              u.Name,
			Email:             u.Email,
			EmailVerified:     true,
			FullName:          u.Name,
			ProviderId:        u.Sub,
			GivenName:         u.GivenName,
			FamilyName:        u.FamilyName,
			PreferredUsername: u.PreferredUsername,
		},
		Emails: []Email{{
			Email:    u.Email,
			Verified: true,
			Primary:  true,
		}},
	}, nil
}
