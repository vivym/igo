package igo

import (
	"encoding/json"
	"math/rand"
	"strings"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36"

// Session holds the login state
type Session struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Auth     struct {
		AuthType              string `json:"auth_type"`
		SessionToken          string `json:"session_token"`
		SessionID             string `json:"session_id"`
		Scnt                  string `json:"scnt"`
		XAppleTwoSVTrustToken string `json:"X-Apple-TwoSV-Trust-Token"`
		Cookies               string `json:"cookies"`
	} `json:"auth"`
	TwoFactorAuthentication bool          `json:"two_factor_authentication"`
	SecurityCode            string        `json:"security_code"`
	ClientSetting           ClientSetting `json:"client_setting"`
	AccountInfo             *AccountInfo  `json:"account_info"`
	Apps                    []AppInfo     `json:"apps"`
	PushInfo                PushInfo      `json:"push_info"`
}

// ClientSetting holds settings of client
type ClientSetting struct {
	ClientID              string            `json:"client_id"`
	Lang                  string            `json:"lang"`
	Locale                string            `json:"locale"`
	XAppleWidgetKey       string            `json:"X-Apple-Widget-Key"`
	XAppleIFDClientInfo   map[string]string `josn:"X-Apple-I-FD-Client-Info"`
	ClientBuildNumber     string            `json:"client_build_number"`
	ClientMasteringNumber string            `json:"client_mastering_number"`
	DefaultHeaders        map[string]string `json:"default_headers"`
}

// PushInfo holds settings of push service
type PushInfo struct {
	Topics        []string `json:"topics"`
	Token         string   `json:"token"`
	TTL           int      `json:"ttl"`
	WebCourierURL string   `json:"web_courier_url"`
}

// XAppleIFDClientInfoStr returns XAppleIFDClientInfo in string format
func (c *ClientSetting) XAppleIFDClientInfoStr() string {
	data, err := json.Marshal(c.XAppleIFDClientInfo)
	if err != nil {
		return ""
	}
	return string(data)
}

// NewSession creates a new default session object
func NewSession() *Session {
	session := Session{
		ClientSetting: ClientSetting{
			ClientID:        NewClientID(),
			Lang:            "zh-cn",
			Locale:          "zh-cn_CN",
			XAppleWidgetKey: "d39ba9916b7251055b22c7f910e2ea796ee65e98b2ddecea8f5dde8d9d1a815d",
			XAppleIFDClientInfo: map[string]string{
				"U": userAgent,
				"L": "zh-CN",
				"Z": "GMT+08:00",
				"V": "1.1",
				"F": "",
			},
			ClientBuildNumber:     "2002Hotfix2",
			ClientMasteringNumber: "2002Hotfix2",
			DefaultHeaders: map[string]string{
				"User-Agent":       userAgent,
				"X-Requested-With": "XMLHttpRequest",
			},
		},
		Apps: []AppInfo{
			AppInfo{
				AppName:   "contacts",
				PushTopic: "73f7bfc9253abaaa423eba9a48e9f187994b7bd9",
			},
			AppInfo{
				AppName:   "calendar",
				PushTopic: "dce593a0ac013016a778712b850dc2cf21af8266",
			},
			AppInfo{
				AppName:   "mail",
				PushTopic: "e850b097b840ef10ce5a7ed95b171058c42cc435",
			},
		},
		PushInfo: PushInfo{
			TTL: 43200,
		},
	}

	topics := make([]string, 0)
	for _, app := range session.Apps {
		if app.PushTopic != "" {
			topics = append(topics, app.PushTopic)
		}
	}
	session.PushInfo.Topics = topics
	return &session
}

// NewClientID creates a new client id
func NewClientID() string {
	parts := make([]string, 0, 5)
	charset := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
	for _, length := range []int{8, 4, 4, 4, 12} {
		part := ""
		for i := 0; i < length; i++ {
			part += charset[rand.Intn(len(charset))]
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, "-")
}
