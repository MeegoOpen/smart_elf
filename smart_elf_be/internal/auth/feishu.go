package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type FeishuAuth struct {
	apiHost      string
	pluginID     string
	pluginSecret string
	token        string
	tokenExpiry  time.Time
	mu           sync.Mutex
}

func NewFeishuAuth(apiHost, pluginID, pluginSecret string) *FeishuAuth {
	return &FeishuAuth{
		apiHost:      apiHost,
		pluginID:     pluginID,
		pluginSecret: pluginSecret,
	}
}

func (a *FeishuAuth) GetToken() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.token == "" || time.Now().After(a.tokenExpiry) {
		if err := a.refreshToken(); err != nil {
			return "", err
		}
	}

	return a.token, nil
}

func (a *FeishuAuth) refreshToken() error {
	url := fmt.Sprintf("%s/open_api/authen/plugin_token", a.apiHost)

	payload := map[string]interface{}{
		"plugin_id":     a.pluginID,
		"plugin_secret": a.pluginSecret,
		"type":          1, //todo
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get token, status code: %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Token      string `json:"token"`
			ExpireTime int64  `json:"expire_time"`
		} `json:"data"`
		Error struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Error.Code != 0 {
		return fmt.Errorf("failed to get token: %s", result.Error.Msg)
	}

	a.token = result.Data.Token
	a.tokenExpiry = time.Now().Add(time.Second * time.Duration(result.Data.ExpireTime))

	return nil
}
