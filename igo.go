package igo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/vivym/igo/internal/common/cookiejar"
	"github.com/vivym/igo/internal/service"
	"github.com/vivym/igo/internal/service/drive"
	"github.com/vivym/igo/internal/session"
	"golang.org/x/net/publicsuffix"
)

// Client provides the iCloud instance
type Client struct {
	http      *resty.Client
	session   *session.Session
	cookieJar *cookiejar.Jar
	services  map[string]service.Service
}

// New creates a new iCloud instance
func New() *Client {
	session := session.NewSession()
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	http := resty.New().
		SetCookieJar(cookieJar).
		SetRetryCount(3).
		SetRetryWaitTime(3 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second).
		SetHeaders(session.ClientSetting.DefaultHeaders)
	client := Client{
		http:      http,
		session:   session,
		cookieJar: cookieJar,
		services:  make(map[string]service.Service),
	}
	return &client
}

// Close ensures everything used is released
func (c *Client) Close() {
	// TODO: autoSave (session)
	for _, service := range c.services {
		service.Stop()
	}
}

// SaveSession dumps session to a writer
func (c *Client) SaveSession(writer io.Writer) error {
	c.session.Auth.Cookies, _ = c.cookieJar.Dumps()
	data, err := json.Marshal(c.session)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// LoadSession loads session from a reader
func (c *Client) LoadSession(reader io.Reader) error {
	err := json.NewDecoder(reader).Decode(c.session)
	if err != nil {
		return err
	}
	c.http.SetHeaders(c.session.ClientSetting.DefaultHeaders)
	cookies := c.session.Auth.Cookies
	if cookies != "" {
		return c.cookieJar.Loads(cookies)
	}
	return nil
}

// Login starts to login
func (c *Client) Login(username, password string) error {
	c.session.Username = username
	c.session.Password = password

	if err := c.signin(); err != nil {
		return err
	}
	if err := c.accountLogin(); err != nil {
		return err
	}
	return nil
}

// TwoFactorAuthenticationIsRequired checks whether the two factor authentication is required.
func (c *Client) TwoFactorAuthenticationIsRequired() bool {
	return c.session.TwoFactorAuthentication && c.session.Auth.XAppleTwoSVTrustToken == ""
}

// SetSecurityCode sets the SecurityCode and trusts the Client
func (c *Client) SetSecurityCode(securityCode string) error {
	c.session.SecurityCode = securityCode
	if err := c.securitycode(); err != nil {
		return err
	}
	if err := c.trust(); err != nil {
		return err
	}
	if err := c.accountLogin(); err != nil {
		return err
	}
	return nil
}

func (c *Client) signin() error {
	signinBody := struct {
		AccountName string        `json:"accountName"`
		Password    string        `json:"password"`
		RememberMe  bool          `json:"rememberMe"`
		TrustTokens []interface{} `json:"trustTokens"`
	}{
		AccountName: c.session.Username,
		Password:    c.session.Password,
		RememberMe:  true,
	}

	rsp, err := c.http.R().
		SetHeaders(map[string]string{
			"Referer":                  "https://idmsa.apple.com/appleauth/auth/signin",
			"Accept":                   "application/json, text/javascript, */*; q=0.01",
			"Origin":                   "https://idmsa.apple.com",
			"X-Apple-Widget-Key":       c.session.ClientSetting.XAppleWidgetKey,
			"X-Apple-I-FD-Client-Info": c.session.ClientSetting.XAppleIFDClientInfoStr(),
		}).
		SetBody(signinBody).
		Post("https://idmsa.apple.com/appleauth/auth/signin")
	if err != nil {
		return err
	}

	res := &struct {
		Reason   string `json:"reason,omitempty"`
		AuthType string `json:"authType"`
	}{}
	if err := json.Unmarshal(rsp.Body(), res); err != nil {
		return err
	}
	if res.Reason != "" {
		return errors.New(res.Reason)
	}

	c.session.Auth.AuthType = res.AuthType
	c.session.Auth.SessionToken = rsp.Header().Get("X-Apple-Session-Token")
	c.session.Auth.SessionID = rsp.Header().Get("X-Apple-ID-Session-Id")
	c.session.Auth.Scnt = rsp.Header().Get("scnt")

	if res.AuthType == "hsa2" {
		c.session.TwoFactorAuthentication = true
	}

	if c.session.Auth.SessionToken == "" {
		return errors.New("no session token")
	}
	return nil
}

func (c *Client) accountLogin() error {
	reqBody := struct {
		DsWebAuthToken string `json:"dsWebAuthToken"`
		ExtendedLogin  bool   `json:"extended_login"`
		TrustToken     string `json:"trustToken,omitempty"`
	}{
		DsWebAuthToken: c.session.Auth.SessionToken,
		ExtendedLogin:  true,
		TrustToken:     c.session.Auth.XAppleTwoSVTrustToken,
	}

	rsp, err := c.http.R().
		SetQueryParams(map[string]string{
			"clientBuildNumber":     c.session.ClientSetting.ClientBuildNumber,
			"clientMasteringNumber": c.session.ClientSetting.ClientMasteringNumber,
			"clientId":              c.session.ClientSetting.ClientID,
		}).
		SetHeader("Origin", "https://www.icloud.com").
		SetHeader("Referer", "https://www.icloud.com/").
		SetBody(reqBody).
		Post(accountLoginURL)
	if err != nil {
		return err
	}

	accountInfo := session.AccountInfo{}
	if err := json.Unmarshal(rsp.Body(), &accountInfo); err != nil {
		return err
	}
	if accountInfo.Reason != "" {
		return errors.New(accountInfo.Reason)
	}
	c.session.AccountInfo = &accountInfo

	return nil
}

func (c *Client) securitycode() error {
	referer := "https://idmsa.apple.com/appleauth/auth/authorize/signin?client_id=" +
		c.session.ClientSetting.ClientID +
		"&response_mode=web_message&response_type=code"
	reqBody := struct {
		SecurityCode map[string]string `json:"securityCode"`
	}{
		SecurityCode: map[string]string{
			"code": c.session.SecurityCode,
		},
	}

	rsp, err := c.http.R().
		SetHeaders(map[string]string{
			// "Origin":                   "https://www.icloud.com",
			"Referer":                  referer,
			"Host":                     "idmsa.apple.com",
			"scnt":                     c.session.Auth.Scnt,
			"X-Apple-Widget-Key":       c.session.ClientSetting.XAppleWidgetKey,
			"X-Apple-I-FD-Client-Info": c.session.ClientSetting.XAppleIFDClientInfoStr(),
			"X-Apple-ID-Session-Id":    c.session.Auth.SessionID,
			"Content-Type":             "application/json",
		}).
		// "{\"securityCode\": {\"code\": \"" + c.session.SecurityCode + "\"}}"
		SetBody(reqBody).
		Post("https://idmsa.apple.com/appleauth/auth/verify/trusteddevice/securitycode")
	if err != nil {
		return err
	}

	if rsp.StatusCode() != 204 {
		return fmt.Errorf("status code: %d", rsp.StatusCode())
	}

	return nil
}

func (c *Client) trust() error {
	referer := "https://idmsa.apple.com/appleauth/auth/authorize/signin?client_id=" +
		c.session.ClientSetting.ClientID +
		"&response_mode=web_message&response_type=code"

	rsp, err := c.http.R().
		SetHeaders(map[string]string{
			"Referer":                  referer,
			"scnt":                     c.session.Auth.Scnt,
			"X-Apple-Widget-Key":       c.session.ClientSetting.XAppleWidgetKey,
			"X-Apple-I-FD-Client-Info": c.session.ClientSetting.XAppleIFDClientInfoStr(),
			"X-Apple-ID-Session-Id":    c.session.Auth.SessionID,
		}).
		Get("https://idmsa.apple.com/appleauth/auth/2sv/trust")
	if err != nil {
		return err
	}

	c.session.Auth.SessionToken = rsp.Header().Get("X-Apple-Session-Token")
	c.session.Auth.XAppleTwoSVTrustToken = rsp.Header().Get("X-Apple-TwoSV-Trust-Token")

	return nil
}

// EnableDrive enables iCloud Drive service
func (c *Client) EnableDrive() *drive.Drive {
	if _, ok := c.services["drive"]; !ok {
		service := drive.New(c.http, c.session)
		c.services["drive"] = service
	}
	drive, _ := c.services["drive"].(*drive.Drive)
	return drive
}
