package mongodemo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
	Age  int                `bson:"age"`
}

func (u user) BSON() bson.M {
	return bson.M{
		"_id":  u.ID,
		"name": u.Name,
		"age":  u.Age,
	}
}

var (
	ctx = context.Background()
)

func TestGetOrCreate(t *testing.T) {
	u := &user{ID: primitive.NewObjectID(), Name: "john", Age: 18}
	t.Run("create and get", func(t *testing.T) {
		db, err := newTestDB(ctx)
		require.NoError(t, err)
		col := NewCollection(db.Collection(t.Name()))
		require.NoError(t, col.c.Drop(ctx))

		// get or create a new user
		createdUser := &user{}
		created, err := col.GetOrCreate(ctx, nil, u.BSON(), createdUser)
		require.NoError(t, err)
		require.True(t, created)
		require.NotEmpty(t, createdUser.ID, "should create ID")
		require.Equal(t, u.Name, createdUser.Name)
		require.Equal(t, u.Age, createdUser.Age)

		// get or create the same user
		gotUser := &user{}
		created, err = col.GetOrCreate(ctx, bson.M{"_id": u.ID}, bson.M{
			"name": u.Name,
			"age":  19,
		}, gotUser)
		require.NoError(t, err)
		require.False(t, created)
		require.Equal(t, createdUser, gotUser, "should be the same user and should not update the age since no insert")

		// malformed field
		updateRes, err := col.c.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{"$set": bson.M{"age": "an invalid field should not be string"}})
		require.Equal(t, int64(1), updateRes.ModifiedCount)
		require.NoError(t, err)
		created, err = col.GetOrCreate(ctx, nil, u.BSON(), createdUser)
		require.Error(t, err)
		require.Contains(t, err.Error(), "fail to unmarshal")
		require.False(t, created)
	})
}
