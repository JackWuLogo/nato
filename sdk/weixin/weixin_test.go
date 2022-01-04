package weixin

import (
	"fmt"
	"testing"
)

func TestNewWeixin(t *testing.T) {
	appid := "wx13d73fc14bdf37b0"
	secret := "bc0a58ba15f126823d10fae6f88bf75f"

	wx, err := NewWeixin(appid, secret)
	if err != nil {
		fmt.Printf("NewWeixin Error: %s\n", err.Error())
		return
	}

	token, err := wx.GetAccessToken()
	if err != nil {
		fmt.Printf("GetAccessToken Error: %s\n", err.Error())
		return
	}

	fmt.Printf("Token: %s\n", token)
}
