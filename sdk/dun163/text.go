package dun163

import (
	"encoding/json"
)

type TextCheckResult struct {
	Antispam struct {
		TaskId       string `json:"taskId"`       // 检测任务ID
		DataId       string `json:"dataId"`       // 数据ID
		Suggestion   int    `json:"suggestion"`   // 建议动作，0：通过，1：嫌疑，2：不通过
		ResultType   int    `json:"resultType"`   // 结果类型，1：机器结果，2：人审结果
		IsRelatedHit bool   `json:"isRelatedHit"` // 是否关联检测命中，true：关联检测命中，false：原文命中
	} `json:"antispam"`
}

// TextCheck 关键词检测
func TextCheck(dun *Dun163, id string, context string) (*TextCheckResult, error) {
	const api = "http://as.dun.163.com/v5/text/check"
	const version = "v5"

	params := dun.CreateParams()
	params.Set("version", version)
	params.Set("dataId", id)
	params.Set("context", context)

	// 发起请求
	rsp, err := dun.Request(api, params)
	if err != nil {
		return nil, err
	}

	var res = new(TextCheckResult)
	if err := json.Unmarshal(rsp, res); err != nil {
		return nil, err
	}

	return res, nil
}
