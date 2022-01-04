package dun163

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"micro-libs/utils/log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Code   int             `json:"code"`
	Msg    string          `json:"msg"`
	Result json.RawMessage `json:"result"`
}

type Dun163 struct {
	sync.Mutex
	SecretId   string
	SecretKey  string
	BusinessId string
}

func (d *Dun163) Set(secretId, secretKey, businessId string) {
	d.Lock()
	defer d.Unlock()

	d.SecretId = secretId
	d.SecretKey = secretKey
	d.BusinessId = businessId

	log.Debug("[Dun163] Reset Dun163 SDK Success. SecretId: %s, SecretKey: %s, BusinessId: %s", secretId, secretKey, businessId)
}

func (d *Dun163) CreateParams() url.Values {
	return url.Values{
		"secretId":   []string{d.SecretId},
		"businessId": []string{d.BusinessId},
		"timestamp":  []string{strconv.FormatInt(time.Now().UnixNano()/1000000, 10)},
		"nonce":      []string{strconv.FormatInt(rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000), 10)},
	}
}

func (d *Dun163) Request(api string, params url.Values) ([]byte, error) {
	params["signature"] = []string{Sign(params, d.SecretKey)}

	resp, err := http.Post(api, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res = new(Result)
	if err := json.Unmarshal(body, res); err != nil {
		return nil, err
	}

	if res.Code != 200 {
		return nil, fmt.Errorf("错误码: %d", res.Code)
	}

	return res.Result, nil
}

func Sign(params url.Values, secret string) string {
	var paramStr string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		paramStr += key + params[key][0]
	}

	paramStr += secret

	md5Reader := md5.New()
	md5Reader.Write([]byte(paramStr))
	return hex.EncodeToString(md5Reader.Sum(nil))
}

func NewDun163(secretId, secretKey, businessId string) *Dun163 {
	d := &Dun163{
		SecretId:   secretId,
		SecretKey:  secretKey,
		BusinessId: businessId,
	}

	log.Debug("[Dun163] New Dun163 SDK Success. SecretId: %s, SecretKey: %s, BusinessId: %s", secretId, secretKey, businessId)

	return d
}
