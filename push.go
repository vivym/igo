package igo

import (
	"encoding/json"
	"errors"
	"strings"
)

// GetPushToken gets pushToken
func (c *Client) GetPushToken() error {
	url := c.session.AccountInfo.Webservices["push"].URL
	url = strings.Replace(url, ":443", "", 1) + "/getToken"

	reqBody := struct {
		PushTopics   []string `json:"pushTopics"`
		PushTokenTTL int      `json:"pushTokenTTL"`
	}{
		PushTopics:   c.session.PushInfo.Topics,
		PushTokenTTL: 43200,
	}

	rsp, err := c.http.R().
		EnableTrace().
		SetPathParams(map[string]string{
			"attempt":               "1",
			"clientBuildNumber":     c.session.ClientSetting.ClientBuildNumber,
			"clientMasteringNumber": c.session.ClientSetting.ClientMasteringNumber,
			"clientId":              c.session.ClientSetting.ClientID,
			"dsid":                  c.session.AccountInfo.DsInfo.DsID,
		}).
		SetHeader("Origin", "https://www.icloud.com").
		SetHeader("Referer", "https://www.icloud.com/").
		SetBody(reqBody).
		Post(url)
	if err != nil {
		return err
	}

	res := struct {
		Reason           string   `json:"reason,omitempty"`
		PushTokenTTL     int      `json:"pushTokenTTL"`
		WebCourierURL    string   `json:"webCourierURL"`
		RegisteredTopics []string `json:"registeredTopics"`
		PushToken        string   `json:"pushToken"`
	}{}
	if err := json.Unmarshal(rsp.Body(), &res); err != nil {
		return err
	}
	if res.Reason != "" {
		return errors.New(res.Reason)
	}
	c.session.PushInfo.TTL = res.PushTokenTTL
	c.session.PushInfo.Topics = res.RegisteredTopics
	c.session.PushInfo.Token = res.PushToken
	c.session.PushInfo.WebCourierURL = res.WebCourierURL

	return nil
}
