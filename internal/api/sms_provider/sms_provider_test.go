package sms_provider

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/supabase/gotrue/internal/conf"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"net/url"
	"testing"
)

var handleApiRequest func(*http.Request) (*http.Response, error)

type SmsProviderTestSuite struct {
	suite.Suite
	Config *conf.GlobalConfiguration
}

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return handleApiRequest(req)
}

func TestSmsProvider(t *testing.T) {
	ts := &SmsProviderTestSuite{
		Config: &conf.GlobalConfiguration{
			Sms: conf.SmsProviderConfiguration{
				Twilio: conf.TwilioProviderConfiguration{
					AccountSid:        "test_account_sid",
					AuthToken:         "test_auth_token",
					MessageServiceSid: "test_message_service_id",
				},
				Messagebird: conf.MessagebirdProviderConfiguration{
					AccessKey:  "test_access_key",
					Originator: "test_originator",
				},
				Vonage: conf.VonageProviderConfiguration{
					ApiKey:    "test_api_key",
					ApiSecret: "test_api_secret",
					From:      "test_from",
				},
				Textlocal: conf.TextlocalProviderConfiguration{
					ApiKey: "test_api_key",
					Sender: "test_sender",
				},
				FlashMobile: conf.FlashMobileProviderConfiguration{
					User:    "test_user",
					Pass:    "test_pass",
					Masking: "test_masking",
				},
				FlashMobileV3: conf.FlashMobileV3ProviderConfiguration{
					ClientKey: "FM-0018-1dc3f842668e24872cddab15",
					ServerKey: "FMPA-edd01c47c0d48e1ee660b4a5d5f401112dbc",
					Masking:   "test_masking",
				},
			},
		},
	}
	suite.Run(t, ts)
}

func (ts *SmsProviderTestSuite) TestTwilioSendSms() {
	defer gock.Off()
	provider, err := NewTwilioProvider(ts.Config.Sms.Twilio)
	require.NoError(ts.T(), err)

	twilioProvider, ok := provider.(*TwilioProvider)
	require.Equal(ts.T(), true, ok)

	phone := "123456789"
	message := "This is the sms code: 123456"

	body := url.Values{
		"To":      {"+" + phone},
		"Channel": {"sms"},
		"From":    {twilioProvider.Config.MessageServiceSid},
		"Body":    {message},
	}

	cases := []struct {
		Desc           string
		TwilioResponse *gock.Response
		ExpectedError  error
	}{
		{
			Desc: "Successfully sent sms",
			TwilioResponse: gock.New(twilioProvider.APIPath).Post("").
				MatchHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(twilioProvider.Config.AccountSid+":"+twilioProvider.Config.AuthToken))).
				MatchType("url").BodyString(body.Encode()).
				Reply(200).JSON(SmsStatus{
				To:     "+" + phone,
				From:   twilioProvider.Config.MessageServiceSid,
				Status: "sent",
				Body:   message,
			}),
			ExpectedError: nil,
		},
		{
			Desc: "Sms status is failed / undelivered",
			TwilioResponse: gock.New(twilioProvider.APIPath).Post("").
				MatchHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(twilioProvider.Config.AccountSid+":"+twilioProvider.Config.AuthToken))).
				MatchType("url").BodyString(body.Encode()).
				Reply(200).JSON(SmsStatus{
				ErrorMessage: "failed to send sms",
				ErrorCode:    "401",
				Status:       "failed",
			}),
			ExpectedError: fmt.Errorf("twilio error: %v %v", "failed to send sms", "401"),
		},
		{
			Desc: "Non-2xx status code returned",
			TwilioResponse: gock.New(twilioProvider.APIPath).Post("").
				MatchHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(twilioProvider.Config.AccountSid+":"+twilioProvider.Config.AuthToken))).
				MatchType("url").BodyString(body.Encode()).
				Reply(500).JSON(twilioErrResponse{
				Code:     500,
				Message:  "Internal server error",
				MoreInfo: "error",
				Status:   500,
			}),
			ExpectedError: &twilioErrResponse{
				Code:     500,
				Message:  "Internal server error",
				MoreInfo: "error",
				Status:   500,
			},
		},
	}

	for _, c := range cases {
		ts.Run(c.Desc, func() {
			err = twilioProvider.SendSms(phone, message, SMSProvider)
			require.Equal(ts.T(), c.ExpectedError, err)
		})
	}
}

func (ts *SmsProviderTestSuite) TestMessagebirdSendSms() {
	defer gock.Off()
	provider, err := NewMessagebirdProvider(ts.Config.Sms.Messagebird)
	require.NoError(ts.T(), err)

	messagebirdProvider, ok := provider.(*MessagebirdProvider)
	require.Equal(ts.T(), true, ok)

	phone := "123456789"
	message := "This is the sms code: 123456"
	body := url.Values{
		"originator": {messagebirdProvider.Config.Originator},
		"body":       {message},
		"recipients": {phone},
		"type":       {"sms"},
		"datacoding": {"unicode"},
	}
	gock.New(messagebirdProvider.APIPath).Post("").MatchHeader("Authorization", "AccessKey "+messagebirdProvider.Config.AccessKey).MatchType("url").BodyString(body.Encode()).Reply(200).JSON(MessagebirdResponse{
		Recipients: MessagebirdResponseRecipients{
			TotalSentCount: 1,
		},
	})

	err = messagebirdProvider.SendSms(phone, message)
	require.NoError(ts.T(), err)
}

func (ts *SmsProviderTestSuite) TestVonageSendSms() {
	defer gock.Off()
	provider, err := NewVonageProvider(ts.Config.Sms.Vonage)
	require.NoError(ts.T(), err)

	vonageProvider, ok := provider.(*VonageProvider)
	require.Equal(ts.T(), true, ok)

	phone := "123456789"
	message := "This is the sms code: 123456"

	body := url.Values{
		"from":       {vonageProvider.Config.From},
		"to":         {phone},
		"text":       {message},
		"api_key":    {vonageProvider.Config.ApiKey},
		"api_secret": {vonageProvider.Config.ApiSecret},
	}

	gock.New(vonageProvider.APIPath).Post("").MatchType("url").BodyString(body.Encode()).Reply(200).JSON(VonageResponse{
		Messages: []VonageResponseMessage{
			{Status: "0"},
		},
	})

	err = vonageProvider.SendSms(phone, message)
	require.NoError(ts.T(), err)
}

func (ts *SmsProviderTestSuite) TestTextLocalSendSms() {
	defer gock.Off()
	provider, err := NewTextlocalProvider(ts.Config.Sms.Textlocal)
	require.NoError(ts.T(), err)

	textlocalProvider, ok := provider.(*TextlocalProvider)
	require.Equal(ts.T(), true, ok)

	phone := "123456789"
	message := "This is the sms code: 123456"
	body := url.Values{
		"sender":  {textlocalProvider.Config.Sender},
		"apikey":  {textlocalProvider.Config.ApiKey},
		"message": {message},
		"numbers": {phone},
	}

	gock.New(textlocalProvider.APIPath).Post("").MatchType("url").BodyString(body.Encode()).Reply(200).JSON(TextlocalResponse{
		Status: "success",
		Errors: []TextlocalError{},
	})

	err = textlocalProvider.SendSms(phone, message)
	require.NoError(ts.T(), err)
}

func (ts *SmsProviderTestSuite) TestFlashMobileSendSms() {
	defer gock.Off()
	provider, err := NewFlashMobileProvider(ts.Config.Sms.FlashMobile)
	require.NoError(ts.T(), err)

	flashMobileProvider, ok := provider.(*FlashMobileProvider)
	require.Equal(ts.T(), true, ok)

	phone := "123456789"
	message := "This is the sms code: 123456"

	requestURL, _ := url.Parse(flashMobileProvider.APIPath)
	urlQuery := requestURL.Query()
	urlQuery.Set("uid", flashMobileProvider.Config.User)
	urlQuery.Set("password", flashMobileProvider.Config.Pass)
	urlQuery.Set("sender", flashMobileProvider.Config.Masking)
	urlQuery.Set("phone", phone)
	urlQuery.Set("text", message)
	requestURL.RawQuery = urlQuery.Encode()

	gock.New(flashMobileProvider.APIPath).Get(requestURL.RawPath).Reply(200).JSON(FlashMobileResponse{
		Status:  1,
		Message: "Success",
		MsgID:   "1234",
	})

	err = flashMobileProvider.SendSms(phone, message)
	require.NoError(ts.T(), err)
}

func (ts *SmsProviderTestSuite) TestFlashMobileV3SendSms() {
	defer gock.Off()
	provider, err := NewFlashMobileV3Provider(ts.Config.Sms.FlashMobileV3)
	require.NoError(ts.T(), err)

	flashMobileProviderV3, ok := provider.(*FlashMobileV3Provider)
	require.Equal(ts.T(), true, ok)

	phone := "082213770600"
	message := "This is the sms code: 123456"

	gock.New(flashMobileProviderV3.APIPath).Post(defaultFlashMobileAuthPath).Reply(200).JSON(GeneralFlashMobileV3Response{
		Status:      200,
		Message:     "success",
		Description: "success",
		Data: FlashMobileToken{
			Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDA2MjU1NjksIm1lcmNoYW50X2lkIjoiMTgiLCJtZXJjaGFudF9uYW1lIjoiYWxhZGlubWFsbCIsIm1lcmNoYW50X3V1aWQiOiIwZTcwMDQ3MTAyYzg0Nzg1YTUzY2I5MTMxYmZlZjEwZSJ9.6S0JjrAvopRejJwiiDGZO6KdikQBAoTON8qan--iWp0",
		},
		Meta: struct{}{},
	})

	gock.New(flashMobileProviderV3.APIPath).Post("/sms/v1/single").Reply(200).JSON(GeneralFlashMobileV3Response{
		Status:      200,
		Message:     "success",
		Description: "success",
		Data: map[string]any{
			"external_id":    "EXT-1000004",
			"transaction_id": "FM-ad74e85e661a2e60933e402a5cc",
			"phone":          "6282213770600",
		},
		Meta: struct{}{},
	})

	err = flashMobileProviderV3.dispatch(phone, message)
	require.NoError(ts.T(), err)
}
