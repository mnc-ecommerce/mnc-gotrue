package mailer

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/supabase/gotrue/internal/conf"
	"github.com/supabase/gotrue/internal/models"
)

type MailClient interface {
	Mail(string, string, string, string, map[string]interface{}) error
}

// TemplateMailer will send mail and use templates from the site for easy mail styling
type TemplateMailer struct {
	SiteURL string
	Config  *conf.GlobalConfiguration
	Mailer  MailClient
}

func encodeRedirectParam(referrerURL string) string {
	redirectParam := ""
	if len(referrerURL) > 0 {
		if strings.ContainsAny(referrerURL, "&=#") {
			// if the string contains &, = or # it has not been URL
			// encoded by the caller, which means it should be URL
			// encoded by us otherwise, it should be taken as-is
			referrerURL = url.QueryEscape(referrerURL)
		}

		redirectParam = "&redirect_to=" + referrerURL
	}

	return redirectParam
}

func addLayout(content string, config *conf.GlobalConfiguration) string {
	return strings.ReplaceAll(BaseLayout(config), "{{content}}", content)
}

const defaultInviteMail = `
<div class="sm-w-280px" style="margin-bottom: 8px; height: auto; width: 620px; background-color: #F5F5F5">
<div class="sm-pl-5px" style="display: block; padding-left: 20px">
<p>You have been invited to create a user on {{ .SiteURL }}. Follow this link to accept the invite:</p>
<p><a href="{{ .ConfirmationURL }}">Accept the invite</a></p>
<p>Alternatively, enter the code: {{ .Token }}</p>
</div>
</div>
`

const defaultConfirmationMail = `
<p>Berikut adalah kode OTP Anda untuk mengakses akun Anda::</p>
<p style="border-width:3px; border-style:solid; border-color:#FF0000; padding: 1em;"><strong>{{ .Token }}</strong></p>
`

const defaultRecoveryMail = `
<p>Kami menerima permintaan reset password untuk akun Anda. Kami memahami betapa pentingnya keamanan akun Anda, dan kami siap membantu Anda mereset password dengan cepat dan mudah.</p>
<p>Silakan masukkan kode di bawah ini untuk mereset password Anda:</p>
<p style="border-width:3px; border-style:solid; border-color:#FF0000; padding: 1em;"><strong>{{ .Token }}</strong></p>
<p>Jika Anda tidak merasa melakukan permintaan reset password ini, silakan abaikan email ini.</p>
<p>Jika Anda mengalami kesulitan dalam mereset password Anda atau memiliki pertanyaan lainnya, jangan ragu untuk menghubungi tim dukungan kami di cs-aladinmall@misteraladin.com atau hubungi kami di nomor Whatsapp +62 811 113 8080.</p>
`

const defaultMagicLinkMail = `
<div class="sm-w-280px" style="margin-bottom: 8px; height: auto; width: 620px; background-color: #F5F5F5">
<div class="sm-pl-5px" style="display: block; padding-left: 20px">
<p>Follow this link to login:</p>
<p><a href="{{ .ConfirmationURL }}">Log In</a></p>
<p>Alternatively, enter the code: {{ .Token }}</p>
</div>
</div>
`

const defaultEmailChangeMail = `
<p>Confirm the update of your email address from {{ .Email }} to {{ .NewEmail }}:</p>
<p style="border-width:3px; border-style:solid; border-color:#FF0000; padding: 1em;"><strong>{{ .Token }}</strong></p>
`

const defaultReauthenticateMail = `
<p style="border-width:3px; border-style:solid; border-color:#FF0000; padding: 1em;"><strong>{{ .Token }}</strong></p>
`

const defaultSuccessRegisterMail = `
<p>Terima kasih telah bergabung dengan AladinMall! Kami senang sekali Anda menjadi pelanggan baru kami.</p>
<p>Kami ingin memberitahu Anda tentang AladinMall dan apa yang kami tawarkan. AladinMall adalah toko online yang menyediakan produk-produk berkualitas dan terpercaya dengan harga yang terjangkau. Kami selalu berusaha memberikan pengalaman belanja yang mudah, cepat, dan menyenangkan.</p>
<p>Kami memiliki berbagai macam produk, dari pakaian hingga aksesoris, dari kebutuhan rumah tangga hingga perlengkapan olahraga. Selain itu, kami juga menawarkan diskon menarik dan promo spesial untuk pelanggan setia kami.</p>
<p>Jangan lupa untuk mendaftar ke newsletter kami untuk mendapatkan informasi tentang produk terbaru dan promo eksklusif. Kami juga selalu siap membantu jika Anda memiliki pertanyaan atau masalah. Hubungi kami melalui email atau live chat.</p>
<p>Terima kasih atas kepercayaan Anda pada AladinMall. Kami harap Anda menikmati pengalaman belanja Anda di sini!</p>`

// ValidateEmail returns nil if the email is valid,
// otherwise an error indicating the reason it is invalid
func (m TemplateMailer) ValidateEmail(email string) error {
	return checkmail.ValidateFormat(email)
}

// InviteMail sends a invite mail to a new user
func (m *TemplateMailer) InviteMail(user *models.User, otp, referrerURL string) error {
	redirectParam := encodeRedirectParam(referrerURL)

	url, err := getSiteURL(referrerURL, m.Config.API.ExternalURL, m.Config.Mailer.URLPaths.Invite, "token="+user.ConfirmationToken+"&type=invite"+redirectParam)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"SiteURL":         m.Config.SiteURL,
		"ConfirmationURL": url,
		"Email":           user.Email,
		"Token":           otp,
		"TokenHash":       user.ConfirmationToken,
		"Data":            user.UserMetaData,
	}

	return m.Mailer.Mail(
		user.GetEmail(),
		string(withDefault(m.Config.Mailer.Subjects.Invite, "You have been invited")),
		m.Config.Mailer.Templates.Invite,
		addLayout(defaultInviteMail, m.Config),
		data,
	)
}

// ConfirmationMail sends a signup confirmation mail to a new user
func (m *TemplateMailer) ConfirmationMail(user *models.User, otp, referrerURL string) error {
	redirectParam := encodeRedirectParam(referrerURL)
	fragment := "token=" + user.ConfirmationToken + "&type=signup" + redirectParam
	url, err := getSiteURL(referrerURL, m.Config.API.ExternalURL, m.Config.Mailer.URLPaths.Confirmation, fragment)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"SiteURL":         m.Config.SiteURL,
		"ConfirmationURL": url,
		"Email":           user.Email,
		"Token":           otp,
		"TokenHash":       user.ConfirmationToken,
		"Data":            user.UserMetaData,
	}

	return m.Mailer.Mail(
		user.GetEmail(),
		string(withDefault(m.Config.Mailer.Subjects.Confirmation, "Confirm Your Email")),
		m.Config.Mailer.Templates.Confirmation,
		addLayout(defaultConfirmationMail, m.Config),
		data,
	)
}

// ReauthenticateMail sends a reauthentication mail to an authenticated user
func (m *TemplateMailer) ReauthenticateMail(user *models.User, otp string) error {
	data := map[string]interface{}{
		"SiteURL": m.Config.SiteURL,
		"Email":   user.Email,
		"Token":   otp,
		"Data":    user.UserMetaData,
	}

	return m.Mailer.Mail(
		user.GetEmail(),
		string(withDefault(m.Config.Mailer.Subjects.Reauthentication, "Confirm reauthentication")),
		m.Config.Mailer.Templates.Reauthentication,
		addLayout(defaultReauthenticateMail, m.Config),
		data,
	)
}

// EmailChangeMail sends an email change confirmation mail to a user
func (m *TemplateMailer) EmailChangeMail(user *models.User, otpNew, otpCurrent, referrerURL string) error {
	type Email struct {
		Address   string
		Otp       string
		TokenHash string
		Subject   string
		Template  string
	}
	emails := []Email{
		{
			Address:   user.EmailChange,
			Otp:       otpNew,
			TokenHash: user.EmailChangeTokenNew,
			Subject:   string(withDefault(m.Config.Mailer.Subjects.EmailChange, "Confirm Email Change")),
			Template:  m.Config.Mailer.Templates.EmailChange,
		},
	}

	currentEmail := user.GetEmail()
	if m.Config.Mailer.SecureEmailChangeEnabled && currentEmail != "" {
		emails = append(emails, Email{
			Address:   currentEmail,
			Otp:       otpCurrent,
			TokenHash: user.EmailChangeTokenCurrent,
			Subject:   string(withDefault(m.Config.Mailer.Subjects.Confirmation, "Confirm Email Address")),
			Template:  m.Config.Mailer.Templates.EmailChange,
		})
	}

	redirectParam := encodeRedirectParam(referrerURL)

	errors := make(chan error)
	for _, email := range emails {
		url, err := getSiteURL(
			referrerURL,
			m.Config.API.ExternalURL,
			m.Config.Mailer.URLPaths.EmailChange,
			"token="+email.TokenHash+"&type=email_change"+redirectParam,
		)
		if err != nil {
			return err
		}
		go func(address, token, tokenHash, template string) {
			data := map[string]interface{}{
				"SiteURL":         m.Config.SiteURL,
				"ConfirmationURL": url,
				"Email":           user.GetEmail(),
				"NewEmail":        user.EmailChange,
				"Token":           token,
				"TokenHash":       tokenHash,
				"Data":            user.UserMetaData,
			}
			errors <- m.Mailer.Mail(
				address,
				string(withDefault(m.Config.Mailer.Subjects.EmailChange, "Confirm Email Change")),
				template,
				addLayout(defaultEmailChangeMail, m.Config),
				data,
			)
		}(email.Address, email.Otp, email.TokenHash, email.Template)
	}

	for i := 0; i < len(emails); i++ {
		e := <-errors
		if e != nil {
			return e
		}
	}

	return nil
}

// RecoveryMail sends a password recovery mail
func (m *TemplateMailer) RecoveryMail(user *models.User, otp, referrerURL string) error {
	redirectParam := encodeRedirectParam(referrerURL)
	url, err := getSiteURL(referrerURL, m.Config.API.ExternalURL, m.Config.Mailer.URLPaths.Recovery, "token="+user.RecoveryToken+"&type=recovery"+redirectParam)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"SiteURL":         m.Config.SiteURL,
		"ConfirmationURL": url,
		"Email":           user.Email,
		"Token":           otp,
		"TokenHash":       user.RecoveryToken,
		"Data":            user.UserMetaData,
	}

	return m.Mailer.Mail(
		user.GetEmail(),
		string(withDefault(m.Config.Mailer.Subjects.Recovery, "Reset Your Password")),
		m.Config.Mailer.Templates.Recovery,
		addLayout(defaultRecoveryMail, m.Config),
		data,
	)
}

// MagicLinkMail sends a login link mail
func (m *TemplateMailer) MagicLinkMail(user *models.User, otp, referrerURL string) error {
	redirectParam := encodeRedirectParam(referrerURL)
	fragment := "token=" + user.RecoveryToken + "&type=magiclink" + redirectParam
	url, err := getSiteURL(referrerURL, m.Config.API.ExternalURL, m.Config.Mailer.URLPaths.Recovery, fragment)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"SiteURL":         m.Config.SiteURL,
		"ConfirmationURL": url,
		"Email":           user.Email,
		"Token":           otp,
		"TokenHash":       user.RecoveryToken,
		"Data":            user.UserMetaData,
	}

	return m.Mailer.Mail(
		user.GetEmail(),
		string(withDefault(m.Config.Mailer.Subjects.MagicLink, "Your Magic Link")),
		m.Config.Mailer.Templates.MagicLink,
		addLayout(defaultMagicLinkMail, m.Config),
		data,
	)
}

// Send can be used to send one-off emails to users
func (m TemplateMailer) Send(user *models.User, subject, body string, data map[string]interface{}) error {
	return m.Mailer.Mail(
		user.GetEmail(),
		subject,
		"",
		body,
		data,
	)
}

// GetEmailActionLink returns a magiclink, recovery or invite link based on the actionType passed.
func (m TemplateMailer) GetEmailActionLink(user *models.User, actionType, referrerURL string) (string, error) {
	var err error

	redirectParam := encodeRedirectParam(referrerURL)

	var url string
	switch actionType {
	case "magiclink":
		url, err = getSiteURL(referrerURL, m.Config.API.ExternalURL, m.Config.Mailer.URLPaths.Recovery, "token="+user.RecoveryToken+"&type=magiclink"+redirectParam)
	case "recovery":
		url, err = getSiteURL(referrerURL, m.Config.API.ExternalURL, m.Config.Mailer.URLPaths.Recovery, "token="+user.RecoveryToken+"&type=recovery"+redirectParam)
	case "invite":
		url, err = getSiteURL(referrerURL, m.Config.API.ExternalURL, m.Config.Mailer.URLPaths.Invite, "token="+user.ConfirmationToken+"&type=invite"+redirectParam)
	case "signup":
		url, err = getSiteURL(referrerURL, m.Config.API.ExternalURL, m.Config.Mailer.URLPaths.Confirmation, "token="+user.ConfirmationToken+"&type=signup"+redirectParam)
	case "email_change_current":
		url, err = getSiteURL(referrerURL, m.Config.API.ExternalURL, m.Config.Mailer.URLPaths.EmailChange, "token="+user.EmailChangeTokenCurrent+"&type=email_change"+redirectParam)
	case "email_change_new":
		url, err = getSiteURL(referrerURL, m.Config.API.ExternalURL, m.Config.Mailer.URLPaths.EmailChange, "token="+user.EmailChangeTokenNew+"&type=email_change"+redirectParam)
	default:
		return "", fmt.Errorf("invalid email action link type: %s", actionType)
	}
	if err != nil {
		return "", err
	}
	return url, nil
}

func (m *TemplateMailer) SuccessSignupMail(user *models.User) error {
	return m.Mailer.Mail(
		user.GetEmail(),
		"Selamat datang di AladinMall",
		"",
		addLayout(defaultSuccessRegisterMail, m.Config),
		map[string]interface{}{},
	)
}
