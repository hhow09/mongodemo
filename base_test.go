package mongodemo

import (
	"context"
	"errors"
	"runtime"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newTestDB(ctx context.Context) (*mongo.Database, error) {
	client, err := newClient(ctx)
	if err != nil {
		return nil, errors.Join(err, errors.New("fail to create client"))
	}
	pc, _, _, _ := runtime.Caller(1)
	path := runtime.FuncForPC(pc).Name()
	db := path[strings.LastIndex(path, "/")+1:]
	db = strings.Replace(db, ".", "_", -1)
	return client.Database(db), nil
}

func newClient(ctx context.Context) (*mongo.Client, error) {
	opt := options.Client().ApplyURI("mongodb://" + "localhost:27017")
	client, err := mongo.Connect(context.Background(), opt)
	if err != nil {
		return nil, errors.Join(err, errors.New("fail to connect"))
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.Join(err, errors.New("fail to ping"))
	}
	return client, nil
}
