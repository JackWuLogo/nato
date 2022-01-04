package dtype

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// 角色基本数据
type RoleBase struct {
	RoleId     int64             `bson:"role_id" json:"role_id" redis:"role_id"`             // 角色ID
	AccountId  int64             `bson:"account_id" redis:"account_id" json:"account_id"`    // 账号ID
	ServerId   int32             `bson:"server_id" redis:"server_id" json:"server_id"`       // 服务器ID
	Online     bool              `bson:"online" json:"online" redis:"online"`                // 在线状态
	Name       string            `bson:"name" redis:"name" json:"name"`                      // 角色名称
	Qualifier  []string          `bson:"qualifier" redis:"qualifier" json:"qualifier"`       // 穿戴称号的修饰词
	NakedAttrs map[int32]float64 `bson:"naked_attrs" json:"naked_attrs" redis:"naked_attrs"` // 扩展潜力点
	CurZlUse   map[int32]int32   `bson:"cur_zl_use" json:"cur_zl_use" redis:"cur_zl_use"`    // 佩戴注灵信息
	Forbid     *Forbid           `bson:"forbid" json:"forbid" redis:"forbid"`                // 角色封禁
}

type Forbid struct {
	Status bool   `bson:"status" json:"status" redis:"status"` // 禁止状态
	Reason string `bson:"reason" json:"reason" redis:"reason"` // 禁止原因
	Expire int64  `bson:"expire" json:"expire" redis:"expire"` // 到期时间
}

// 附魔信息
type EnchantInfo struct {
	Id        int32 `bson:"id" json:"id" redis:"id"`                      // 属性id
	Value     int64 `bson:"value" json:"value" redis:"value"`             // 数值
	Level     int32 `bson:"level" json:"level" redis:"level"`             // 附魔符等阶
	Timestamp int64 `bson:"timestamp" json:"timestamp" redis:"timestamp"` // 到期的时间戳
	ItemId    int32 `bson:"item_id" json:"item_id" redis:"item_id"`       // 道具id(用于发给前端读属性区间)
}

var typ = reflect.TypeOf(RoleBase{})
var testData = &RoleBase{
	RoleId:    1000000000000000,
	AccountId: 123123,
	ServerId:  1,
	Online:    true,
	Name:      "测试",
	Qualifier: nil, // 称号修饰词
	NakedAttrs: map[int32]float64{
		1: 0.1,
		2: 0.2,
	},
	CurZlUse: make(map[int32]int32), // 当前幻化信息
}

func BenchmarkUnsafe(b *testing.B) {
	f := UnsafeStructToByte(testData)
	_ = UnsafeByteToStruct(typ, f)
}

func BenchmarkJson(b *testing.B) {
	bt, _ := json.Marshal(testData)
	c := reflect.New(typ)
	_ = json.Unmarshal(bt, c.Interface())
}

func BenchmarkGob(b *testing.B) {
	var buf = new(bytes.Buffer)
	_ = gob.NewEncoder(buf).Encode(testData)
	c := reflect.New(typ)
	_ = gob.NewDecoder(buf).Decode(c.Interface())
}

func TestUnsafeToByte(t *testing.T) {
	f := UnsafeStructToByte(testData)
	b := UnsafeByteToStruct(typ, f)
	if res, ok := b.Interface().(*RoleBase); ok {
		fmt.Printf("Ok! RoleId: %+v\n", res.NakedAttrs)
	}
}

func TestUnsafeCopyStruct(t *testing.T) {
	b := UnsafeCopyStruct(testData)
	testData.RoleId = 20000000000001
	if res, ok := b.(*RoleBase); ok {
		fmt.Printf("Ok! RoleId: %d, Old: %d\n", res.RoleId, testData.RoleId)
	}
}
