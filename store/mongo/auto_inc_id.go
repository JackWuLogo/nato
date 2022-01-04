package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"micro-libs/utils/errors"
	"strings"
)

const AutoIncNamePrefix = "auto_inc"

// 自增ID数据表
type AutoIncId struct {
	Id  string `bson:"_id" json:"id"`
	Num int64  `bson:"n" json:"n"`
}

// 生成自增ID数据表名
func GetAutoIncName(id string) string {
	return fmt.Sprintf("%s_%s", AutoIncNamePrefix, id)
}

// 生成自增ID数据表名
func CheckTableIsAutoInc(name string) bool {
	return strings.HasPrefix(name, AutoIncNamePrefix)
}

// GetAutoIncId 生成自增ID
func GetAutoIncId(ctx context.Context, db *mongo.Database, id string) (int64, error) {
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	res := db.Collection(GetAutoIncName(id)).FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"n": int64(1)}}, opts)
	if res.Err() != nil {
		return 0, res.Err()
	}

	var incId = new(AutoIncId)
	if err := res.Decode(&incId); err != nil {
		return 0, err
	} else if incId.Num == 0 {
		return 0, errors.Invalid("invalid inc id")
	}

	return incId.Num, nil
}
