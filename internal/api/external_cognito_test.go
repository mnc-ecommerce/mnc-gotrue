package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	jwt "github.com/golang-jwt/jwt"
)

const (
	cognitoUser        string = `{"name":"Cognito Test","email":"cognito@example.com","sub":"cognitotestid"}`
	cognitoUserNoEmail string = `{"name":"Cognito Test","sub":"cognitotestid"}`
)

func (ts *ExternalTestSuite) TestSignupExternalCognito() {
	req := httptest.NewRequest(http.MethodGet, "http://localhost/authorize?provider=cognito", nil)
	w := httptest.NewRecorder()
	ts.API.handler.ServeHTTP(w, req)
	ts.Require().Equal(http.StatusFound, w.Code)
	u, err := url.Parse(w.Header().Get("Location"))
	ts.Require().NoError(err, "redirect url parse failed")
	q := u.Query()
	ts.Equal(ts.Config.External.Cognito.RedirectURI, q.Get("redirect_uri"))
	ts.Equal(ts.Config.External.Cognito.ClientID, q.Get("client_id"))
	ts.Equal("code", q.Get("response_type"))
	ts.Equal("openid", q.Get("scope"))

	claims := ExternalProviderClaims{}
	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}
	_, err = p.ParseWithClaims(q.Get("state"), &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(ts.Config.JWT.Secret), nil
	})
	ts.Require().NoError(err)

	ts.Equal("cognito", claims.Provider)
	ts.Equal(ts.Config.SiteURL, claims.SiteURL)
}

func CognitoTestSignupSetup(ts *ExternalTestSuite, tokenCount *int, userCount *int, code string, user string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/token":
			*tokenCount++
			ts.Equal(code, r.FormValue("code"))
			ts.Equal("authorization_code", r.FormValue("grant_type"))
			ts.Equal(ts.Config.External.Cognito.RedirectURI, r.FormValue("redirect_uri"))

			w.Header().Add("Content-Type", "application/json")
			fmt.Fprint(w, `{"access_token":"cognito_token","expires_in":100000}`)
		case "/oauth2/userinfo":
			*userCount++
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprint(w, user)
		default:
			w.WriteHeader(500)
			ts.Fail("unknown cognito oauth call %s", r.URL.Path)
		}
	}))

	ts.Config.External.Cognito.URL = server.URL
	ts.Config.External.Cognito.ApiURL = server.URL

	return server
}

func (ts *ExternalTestSuite) TestSignupExternalCognito_AuthorizationCode() {
	ts.Config.DisableSignup = false
	tokenCount, userCount := 0, 0
	code := "authcode"
	server := CognitoTestSignupSetup(ts, &tokenCount, &userCount, code, cognitoUser)
	defer server.Close()

	u := performAuthorization(ts, "cognito", code, "")

	assertAuthorizationSuccess(ts, u, tokenCount, userCount, "cognito@example.com", "Cognito Test", "cognitotestid", "")
}

func (ts *ExternalTestSuite) TestSignupExternalCognitoDisableSignupErrorWhenNoUser() {
	ts.Config.DisableSignup = true
	tokenCount, userCount := 0, 0
	code := "authcode"
	server := CognitoTestSignupSetup(ts, &tokenCount, &userCount, code, cognitoUser)
	defer server.Close()

	u := performAuthorization(ts, "cognito", code, "")

	assertAuthorizationFailure(ts, u, "Signups not allowed for this instance", "access_denied", "cognito@example.com")
}

func (ts *ExternalTestSuite) TestSignupExternalCognitoDisableSignupErrorWhenNoEmail() {
	ts.Config.DisableSignup = true
	tokenCount, userCount := 0, 0
	code := "authcode"
	server := CognitoTestSignupSetup(ts, &tokenCount, &userCount, code, cognitoUserNoEmail)
	defer server.Close()

	u := performAuthorization(ts, "cognito", code, "")

	assertAuthorizationFailure(ts, u, "Error getting user email from external provider", "server_error", "cognito@example.com")

}

func (ts *ExternalTestSuite) TestSignupExternalCognitoDisableSignupSuccessWithPrimaryEmail() {
	ts.Config.DisableSignup = true

	ts.createUser("cognitotestid", "cognito@example.com", "Cognito Test", "http://example.com/avatar", "")

	tokenCount, userCount := 0, 0
	code := "authcode"
	server := CognitoTestSignupSetup(ts, &tokenCount, &userCount, code, cognitoUser)
	defer server.Close()

	u := performAuthorization(ts, "cognito", code, "")

	assertAuthorizationSuccess(ts, u, tokenCount, userCount, "cognito@example.com", "Cognito Test", "cognitotestid", "http://example.com/avatar")
}

func (ts *ExternalTestSuite) TestInviteTokenExternalCognitoSuccessWhenMatchingToken() {
	// name should be populated from Cognito API
	ts.createUser("cognitotestid", "cognito@example.com", "", "http://example.com/avatar", "invite_token")

	tokenCount, userCount := 0, 0
	code := "authcode"
	server := CognitoTestSignupSetup(ts, &tokenCount, &userCount, code, cognitoUser)
	defer server.Close()

	u := performAuthorization(ts, "cognito", code, "invite_token")

	assertAuthorizationSuccess(ts, u, tokenCount, userCount, "cognito@example.com", "Cognito Test", "cognitotestid", "http://example.com/avatar")
}

func (ts *ExternalTestSuite) TestInviteTokenExternalCognitoErrorWhenNoMatchingToken() {
	tokenCount, userCount := 0, 0
	code := "authcode"
	cognitoUser := `{"name":"Cognito Test","avatar":{"href":"http://example.com/avatar"}}`
	server := CognitoTestSignupSetup(ts, &tokenCount, &userCount, code, cognitoUser)
	defer server.Close()

	w := performAuthorizationRequest(ts, "cognito", "invite_token")
	ts.Require().Equal(http.StatusNotFound, w.Code)
}

func (ts *ExternalTestSuite) TestInviteTokenExternalCognitoErrorWhenWrongToken() {
	ts.createUser("cognitotestid", "cognito@example.com", "", "", "invite_token")

	tokenCount, userCount := 0, 0
	code := "authcode"
	cognitoUser := `{"name":"Cognito Test","avatar":{"href":"http://example.com/avatar"}}`
	server := CognitoTestSignupSetup(ts, &tokenCount, &userCount, code, cognitoUser)
	defer server.Close()

	w := performAuthorizationRequest(ts, "cognito", "wrong_token")
	ts.Require().Equal(http.StatusNotFound, w.Code)
}

func (ts *ExternalTestSuite) TestInviteTokenExternalCognitoErrorWhenEmailDoesntMatch() {
	ts.createUser("cognitotestid", "cognito@example.com", "", "", "invite_token")

	tokenCount, userCount := 0, 0
	code := "authcode"
	cognitoUser := `{"name":"Cognito Test", "email":"other@example.com", "avatar":{"href":"http://example.com/avatar"}}`
	server := CognitoTestSignupSetup(ts, &tokenCount, &userCount, code, cognitoUser)
	defer server.Close()

	u := performAuthorization(ts, "cognito", code, "invite_token")

	assertAuthorizationFailure(ts, u, "Invited email does not match emails from external provider", "invalid_request", "")
}
