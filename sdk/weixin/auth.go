package weixin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// AuthGetAccessToken 获取小程序全局唯一后台接口调用凭据（access_token）
func AuthGetAccessToken(appId, secret string) (string, int64, error) {
	values := url.Values{}
	values.Add("grant_type", "client_credential")
	values.Add("appid", appId)
	values.Add("secret", secret)

	api := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?%s", values.Encode())

	rsp, err := Request(http.MethodGet, api, nil)
	if err != nil {
		return "", 0, err
	}

	res := new(struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	})
	if err := json.Unmarshal(rsp, res); err != nil {
		return "", 0, err
	}

	return res.AccessToken, res.ExpiresIn, nil
}
