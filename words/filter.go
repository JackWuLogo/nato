package words

import (
	"strings"
	"sync"
)

type Null struct{}

type Filter struct {
	sync.Mutex
	bad       map[string]Null // 屏蔽字
	invalid   map[string]Null // 无效字
	replace   rune            // 替换字符
	sensitive map[string]interface{}
}

func (f *Filter) SetInvalid(invalid string) {
	f.Lock()
	defer f.Unlock()

	f.invalid = make(map[string]Null)

	words := strings.Split(invalid, ",")
	for _, v := range words {
		f.invalid[v] = Null{}
	}
}

func (f *Filter) SetBad(bad []string) {
	f.Lock()
	defer f.Unlock()

	f.bad = make(map[string]Null)
	for _, v := range bad {
		f.bad[v] = Null{}
	}

	f.build()
}

func (f *Filter) SetReplace(replace string) {
	f.Lock()
	defer f.Unlock()

	runes := []rune(replace)
	if len(runes) > 0 {
		f.replace = runes[0]
	} else {
		f.replace = '*'
	}
}

func (f *Filter) build() {
	f.sensitive = make(map[string]interface{})
	for key := range f.bad {
		str := []rune(key)
		nowMap := f.sensitive
		for i := 0; i < len(str); i++ {
			if _, ok := nowMap[string(str[i])]; !ok {
				//如果该key不存在，
				thisMap := make(map[string]interface{})
				thisMap["isEnd"] = false
				nowMap[string(str[i])] = thisMap
				nowMap = thisMap
			} else {
				nowMap = nowMap[string(str[i])].(map[string]interface{})
			}
			if i == len(str)-1 {
				nowMap["isEnd"] = true
			}
		}
	}
}

func (f *Filter) Replace(txt string) (word string) {
	f.Lock()
	defer f.Unlock()

	str := []rune(txt)
	nowMap := f.sensitive
	start := -1
	tag := -1

	for i := 0; i < len(str); i++ {
		if _, ok := f.invalid[(string(str[i]))]; ok || string(str[i]) == "," {
			continue
		}
		if thisMap, ok := nowMap[string(str[i])].(map[string]interface{}); ok {
			tag++
			if tag == 0 {
				start = i
			}
			isEnd, _ := thisMap["isEnd"].(bool)
			if isEnd {
				for y := start; y < i+1; y++ {
					str[y] = f.replace
				}
				nowMap = f.sensitive
				start = -1
				tag = -1
			} else {
				nowMap = nowMap[string(str[i])].(map[string]interface{})
			}
		} else {
			if start != -1 {
				i = start + 1
			}
			nowMap = f.sensitive
			start = -1
			tag = -1
		}
	}

	return string(str)
}

func NewFilter() *Filter {
	return &Filter{
		bad:       make(map[string]Null),
		invalid:   make(map[string]Null),
		replace:   '*',
		sensitive: make(map[string]interface{}),
	}
}
