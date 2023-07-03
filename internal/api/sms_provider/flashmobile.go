package sms_provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/supabase/gotrue/internal/conf"
	"github.com/supabase/gotrue/internal/utilities"
)

const (
	defaultFlashMobileApiBase = "https://sms.flashmobile.co.id:881"
)

var urlFlashMobileApiBase = []string{
	"https://sms.flashmobile.co.id:881",
	"https://sms.flashmobile.co.id:882",
}

type FlashMobileProvider struct {
	Config  *conf.FlashMobileProviderConfiguration
	APIPath string
}

type FlashMobileResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	MsgID   string `json:"msg_id"`
}

// Creates a SmsProvider with the FlashMobile Config
func NewFlashMobileProvider(config conf.FlashMobileProviderConfiguration) (SmsProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	apiPath := defaultFlashMobileApiBase
	return &FlashMobileProvider{
		Config:  &config,
		APIPath: apiPath,
	}, nil
}

func (t *FlashMobileProvider) SendMessage(phone string, message string, channel string) error {
	switch channel {
	case SMSProvider:
		return t.SendSms(phone, message)
	default:
		return fmt.Errorf("channel type %q is not supported for FlashMobile", channel)
	}
}

// Send an SMS containing the OTP with FlashMobile's API
func (t *FlashMobileProvider) SendSms(phone string, message string) error {

	resp := &FlashMobileResponse{}
	var err error
	for _, vBaseUrl := range urlFlashMobileApiBase {
		requestURL, _ := url.Parse(vBaseUrl + "/v1/send")
		urlQuery := requestURL.Query()
		urlQuery.Set("uid", t.Config.User)
		urlQuery.Set("password", t.Config.Pass)
		urlQuery.Set("sender", t.Config.Masking)
		urlQuery.Set("phone", phone)
		urlQuery.Set("text", message)
		requestURL.RawQuery = urlQuery.Encode()

		client := &http.Client{Timeout: defaultTimeout}
		r, err := http.NewRequest("GET", requestURL.String(), nil)
		if err != nil {
			continue
		}

		res, err := client.Do(r)
		fmt.Println("Response FlashMobile SendSms", res)
		if err != nil {
			continue
		}
		defer utilities.SafeClose(res.Body)

		derr := json.NewDecoder(res.Body).Decode(resp)
		if derr != nil {
			err = derr
			continue
		}

		if resp.Status == 0 {
			err = fmt.Errorf("textlocal error: %v", resp.Message)
			continue
		}
	}

	fmt.Println("SendSms", err)
	return err
}
