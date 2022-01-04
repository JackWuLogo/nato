package weixin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"micro-libs/utils/errors"
	"net/http"
	"net/url"
)

func Request(method, api string, body io.Reader) ([]byte, error) {
	// 创建请求对象
	rq, err := http.NewRequest(method, api, body)
	if err != nil {
		return nil, err
	}
	rq.Header.Set("Content-Type", "application/json;charset=utf-8")
	rq.Close = true // 不重用连接

	// 发起请求
	client := new(http.Client)
	rsp, err := client.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New(int32(rsp.StatusCode))
	}

	res, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var rerr = new(struct {
		Code int    `json:"errcode"`
		Msg  string `json:"errmsg"`
	})
	if err := json.Unmarshal(res, rerr); err != nil {
		return nil, err
	}

	if rerr.Code > 0 {
		return nil, errors.New(int32(rerr.Code), rerr.Msg)
	}

	return res, nil
}

func ApiPost(api string, token string, data interface{}) ([]byte, error) {
	var body io.Reader
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}
	var values = url.Values{}
	values.Add("access_token", token)
	return Request(http.MethodPost, fmt.Sprintf("%s?%s", api, values.Encode()), body)
}

func ApiGet(api string, token string, query map[string]string) ([]byte, error) {
	var values = url.Values{}
	values.Add("access_token", token)
	if len(query) > 0 {
		for k, v := range query {
			values.Add(k, v)
		}
	}
	return Request(http.MethodPost, fmt.Sprintf("%s?%s", api, values.Encode()), nil)
}
