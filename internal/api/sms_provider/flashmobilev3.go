package sms_provider

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/supabase/gotrue/internal/conf"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	defaultFlashMobileV3ApiBase = "https://app.flashmobile.co.id"
	defaultFlashMobileAuthPath  = "/auth/v2/access-token"
	defaultFlashMobileSMSPath   = "/sms/v1/single"
)

type FlashMobileAuthRequest struct {
	ClientKey string `json:"client_key"`
	ServerKey string `json:"server_key"`
}

type FlashMobileDispatchRequest struct {
	ExternalId  string `json:"external_id"`
	Text        string `json:"text"`
	Phone       string `json:"phone"`
	Sender      string `json:"sender"`
	CallbackUrl string `json:"callback_url"`
}

type GeneralFlashMobileV3Response struct {
	Status      any    `json:"status"`
	Message     string `json:"message"`
	Description string `json:"description"`
	Data        any    `json:"data"`
	Meta        struct {
	} `json:"meta"`
}

type FlashMobileToken struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"-"`
}

func (t *FlashMobileToken) isExpired() bool {
	return t.Expiry.Before(time.Now())
}

type FlashMobileV3Provider struct {
	APIPath string
	Config  conf.FlashMobileV3ProviderConfiguration
	token   *FlashMobileToken
	lock    *sync.Mutex
}

func NewFlashMobileV3Provider(config conf.FlashMobileV3ProviderConfiguration) (SmsProvider, error) {
	if err := config.Validate(); nil != err {
		return nil, err
	}

	return &FlashMobileV3Provider{
		APIPath: defaultFlashMobileV3ApiBase,
		Config:  config,
		token:   new(FlashMobileToken),
		lock:    new(sync.Mutex),
	}, nil
}

func (p *FlashMobileV3Provider) SendMessage(phone, message, channel string) error {
	switch channel {
	case SMSProvider:
		return p.dispatch(phone, message)
	default:
		return fmt.Errorf("channel type %q is not supported for FlashMobile", channel)
	}
}

func (p *FlashMobileV3Provider) dispatch(phone string, message string) error {
	if p.token.isExpired() {
		if err := p.auth(); nil != err {
			return err
		}
	}

	var (
		response GeneralFlashMobileV3Response
		url      = defaultFlashMobileV3ApiBase + defaultFlashMobileSMSPath
		headers  = map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + p.token.Token,
		}
		request = FlashMobileDispatchRequest{
			ExternalId:  fmt.Sprintf("%d", time.Now().Unix()),
			Text:        message,
			Phone:       phone,
			Sender:      p.Config.Masking,
			CallbackUrl: "https://aladinmall.id",
		}
	)

	requestBody, _ := json.Marshal(request)
	headers["X-Signature"] = p.sign(request.ExternalId, defaultFlashMobileSMSPath, requestBody)

	if err := do(http.MethodPost, url, headers, requestBody, &response); nil != err {
		return err
	}

	return nil
}

func (p *FlashMobileV3Provider) auth() error {
	var (
		response GeneralFlashMobileV3Response
		url      = defaultFlashMobileV3ApiBase + defaultFlashMobileAuthPath
		headers  = map[string]string{
			"Content-Type": "application/json",
		}
		request = FlashMobileAuthRequest{
			ClientKey: p.Config.ClientKey,
			ServerKey: p.Config.ServerKey,
		}
	)

	requestBody, _ := json.Marshal(request)
	if err := do(http.MethodPost, url, headers, requestBody, &response); nil != err {
		return err
	}

	var t FlashMobileToken
	if err := mapstructure.Decode(response.Data, &t); nil != err {
		return err
	}

	p.lock.Lock()
	p.token = &t
	p.token.Expiry = time.Now().Add(time.Hour * 24)
	p.lock.Unlock()

	return nil
}

func (p *FlashMobileV3Provider) sign(extId, path string, payload []byte) string {
	var payloadMinify = new(bytes.Buffer)
	if err := json.Compact(payloadMinify, payload); nil != err {
		return ""
	}

	var (
		str2Sign = fmt.Sprintf(
			"%s|%s|%s/%s",
			p.token.Token,
			extId,
			path,
			payloadMinify.String(),
		)
		hasher = sha256.New()
	)

	hasher.Write([]byte(str2Sign))

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func do(method, url string, headers map[string]string, body []byte, val any) error {
	var bodyReader io.Reader
	if nil != body {
		bodyReader = bytes.NewReader(body)
	}

	request, err := http.NewRequest(method, url, bodyReader)
	if nil != err {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}

	client := &http.Client{Timeout: defaultTimeout}
	response, err := client.Do(request)
	if nil != err {
		return err
	}

	var responseBody []byte
	if nil != response.Body {
		responseBody, err = io.ReadAll(response.Body)
		if nil != err {
			return err
		}
	}

	if response.StatusCode >= 300 {
		return fmt.Errorf("status_code: %d body: %s", response.StatusCode, string(responseBody))
	}

	return json.Unmarshal(responseBody, val)
}
