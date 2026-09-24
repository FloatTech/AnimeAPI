// Package pixiv (copied from plugin/pixiv/api/token.go)
package pixiv

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// TokenStore ...
type TokenStore struct {
	client       *Client
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"-"`
	mu           sync.Mutex
}

// NewTokenStore ...
func NewTokenStore(refreshToken string, c *Client) *TokenStore {
	return &TokenStore{
		RefreshToken: refreshToken,
		client:       c,
	}
}

// GetAccessToken ...
func (t *TokenStore) GetAccessToken() (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if time.Now().Before(t.ExpiresAt) && t.AccessToken != "" {
		log.Println("access_token is valid")
		return t.AccessToken, nil
	}

	if err := t.refreshPixivAccessToken(); err != nil {
		return "", err
	}

	t.ExpiresAt = time.Now().Add(time.Duration(t.ExpiresIn/2) * time.Second)
	return t.AccessToken, nil
}

// refreshPixivAccessToken 用 refresh_token 刷新 access_token
func (t *TokenStore) refreshPixivAccessToken() error {
	endpoint := "https://oauth.secure.pixiv.net/auth/token"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", "MOBrBDS8blbauoSck0ZfDbtuzpyT")
	data.Set("client_secret", "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj")
	data.Set("refresh_token", t.RefreshToken)

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234 (Android 11; Pixel 5)")

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New("刷新失败: " + resp.Status + "\nbody: " + string(body))
	}

	var tokenRes TokenStore
	err = json.Unmarshal(body, &tokenRes)

	t.AccessToken = tokenRes.AccessToken
	t.ExpiresIn = tokenRes.ExpiresIn

	return err
}
