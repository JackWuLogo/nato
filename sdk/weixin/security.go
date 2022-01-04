package weixin

import "encoding/json"

type MsgSecCheck struct {
	TraceId string               `json:"trace_id"`
	Result  MsgSecCheck_Result   `json:"result"`
	Detail  []MsgSecCheck_Detail `json:"detail"`
}

type MsgSecCheck_Result struct {
	Suggust string `json:"suggust"` // 建议，有risky、pass、review三种值
	Label   int    `json:"label"`   // 命中标签枚举值，100 正常；10001 广告；20001 时政；20002 色情；20003 辱骂；20006 违法犯罪；20008 欺诈；20012 低俗；20013 版权；21000 其他
}

type MsgSecCheck_Detail struct {
	Strategy string `json:"strategy"` // 策略类型
	Errcode  int    `json:"errcode"`  // 错误码，仅当该值为0时，该项结果有效
	Suggest  string `json:"suggest"`  // 建议，有risky、pass、review三种值
	Label    int    `json:"label"`    // 命中标签枚举值，100 正常；10001 广告；20001 时政；20002 色情；20003 辱骂；20006 违法犯罪；20008 欺诈；20012 低俗；20013 版权；21000 其他
	Prob     int    `json:"prob"`     // 0-100，代表置信度，越高代表越有可能属于当前返回的标签（label）
	Level    int    `json:"level"`    // 0-100，代表置信度，越高代表越有可能属于当前返回的标签（label）
	Keyword  string `json:"keyword"`  // 命中的自定义关键词
}

func SecurityMsgSecCheck(accessToken string, openid string, scene int, content string) (*MsgSecCheck, error) {
	api := "https://api.weixin.qq.com/wxa/msg_sec_check"
	params := map[string]interface{}{
		"version": 2,
		"openid":  openid,
		"scene":   scene,
		"content": content,
	}
	rsp, err := ApiPost(api, accessToken, params)
	if err != nil {
		return nil, err
	}

	res := new(MsgSecCheck)
	if err := json.Unmarshal(rsp, res); err != nil {
		return nil, err
	}

	return res, nil
}
