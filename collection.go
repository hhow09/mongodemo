package mongodemo

import (
	"context"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

func NewCollection(c *mongo.Collection) *Collection {
	return &Collection{c: c}
}

type Collection struct {
	c *mongo.Collection
}

// GetOrCreate is a helper function to find a document by query, if not found, insert a new document with the patch.
// noted that the field in patch will be overwritten by the field of same name in query.
func (col *Collection) GetOrCreate(ctx context.Context, query, patch bson.M, out interface{}) (bool, error) {
	if patch == nil {
		patch = query
	} else {
		for k, v := range query {
			patch[k] = v
		}
	}
	update := bson.M{"$setOnInsert": patch}
	result := &operation.FindAndModifyResult{}
	res := col.c.Database().RunCommand(ctx, bson.D{
		{Key: "findAndModify", Value: col.c.Name()},
		{Key: "query", Value: query},
		{Key: "update", Value: update},
		{Key: "new", Value: out != nil},
		{Key: "upsert", Value: true},
	})
	if err := res.Decode(result); err != nil {
		return false, errors.Join(err, errors.New("fail to decode result"))
	}
	created := result.LastErrorObject.Upserted != nil
	if out != nil {
		if err := bson.Unmarshal(result.Value, out); err != nil {
			return created, errors.Join(err, errors.New("fail to unmarshal result"))
		}
	}
	return created, nil
}
