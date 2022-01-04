package weixin

import (
	"micro-libs/utils/log"
	"micro-libs/utils/tool"
	"sync"
	"time"
)

type Weixin struct {
	sync.Mutex
	appId  string
	secret string
	token  string
	expire int64
}

func (w *Weixin) Init() error {
	w.Lock()
	defer w.Unlock()

	token, expire, err := AuthGetAccessToken(w.appId, w.secret)
	if err != nil {
		return err
	}

	expireTime := time.Now().Add(time.Duration(expire) * time.Second)

	w.token = token
	w.expire = expireTime.Unix()

	log.Debug("[Weixin] Refresh Weixin Access Token Success, Expire: %s ...", tool.TimeFormat(expireTime, "Y-m-d H:i:s"))

	return nil
}

func (w *Weixin) Set(appId, secret string) {
	w.Lock()
	defer w.Unlock()

	w.appId = appId
	w.secret = secret
}

func (w *Weixin) GetAccessToken() (string, error) {
	if w.token == "" || w.expire-60 < time.Now().Unix() {
		if err := w.Init(); err != nil {
			return "", err
		}
	}
	return w.token, nil
}

func (w *Weixin) SecurityMsgSecCheck(openid string, scene int, content string) (*MsgSecCheck, error) {
	token, err := w.GetAccessToken()
	if err != nil {
		return nil, err
	}
	return SecurityMsgSecCheck(token, openid, scene, content)
}

func NewWeixin(appId, secret string) *Weixin {
	wx := &Weixin{
		appId:  appId,
		secret: secret,
	}

	log.Debug("[Weixin] New Weixin SDK Success. AppId: %s, Secret: %s", appId, secret)

	return wx
}
