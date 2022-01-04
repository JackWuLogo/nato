package scheme

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// VersionModel 配置表版本
type VersionModel struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	Key        string             `bson:"key" json:"key"`                 // 标识, 即该配置在数据库中的表名
	Name       string             `bson:"name" json:"name"`               // 名称
	Total      int                `bson:"total" json:"total"`             // 配置项数量
	Version    int64              `bson:"version" json:"version"`         // 版本号
	UpdateUser int64              `bson:"update_user" json:"update_user"` // 更新者ID
	UpdateTime int64              `bson:"update_time" json:"update_time"` // 更新时间
	CreateUser int64              `bson:"create_user" json:"create_user"` // 创建者ID
	CreateTime int64              `bson:"create_time" json:"create_time"` // 创建时间
}

var (
	VersionModelIndex = []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "key", Value: bsonx.Int32(1)},
			},
			Options: options.Index().SetUnique(true),
		},
	}
)

// Client 客户端导出数据
type Client struct {
	Table   string `json:"table"`   // 数据表key
	Version int64  `json:"version"` // 数据表版本
	Attrs   []byte `json:"attrs"`   // 数据表配置信息
}
