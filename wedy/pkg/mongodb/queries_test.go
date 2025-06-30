package mongodb

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

type MongoDBTestSuit struct {
	suite.Suite
	col *mongo.Collection
}

func (d *MongoDBTestSuit) SetupSuite() {
	t := d.T()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
	}
	opts := options.Client().
		ApplyURI("mongodb://14.103.175.18:27017/").
		SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)
	col := client.Database("wedy").Collection("videos")
	d.col = col
	manyRes, err := d.col.InsertMany(ctx, []any{Video{
		Uid:  123456,
		VUid: 654321,
	}, Video{
		Uid:  1234567,
		VUid: 7654321,
	}, Video{
		Uid:  12345678,
		VUid: 87654321,
	}, Video{
		Uid:  123456,
		VUid: 76543210,
	}})
	assert.NoError(t, err)
	d.T().Log("Inserted IDs: ", manyRes.InsertedIDs)
}
func (d *MongoDBTestSuit) TestOR() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	filter := bson.A{
		bson.D{bson.E{Key: "uid", Value: 123456}},
		bson.D{bson.E{Key: "uid", Value: 654321}},
	}
	findRes, err := d.col.Find(ctx, bson.D{bson.E{Key: "$or", Value: filter}})
	assert.NoError(d.T(), err)
	var videos []Video
	err = findRes.All(ctx, &videos)
	assert.NoError(d.T(), err)
	d.T().Log("Inserted IDs: ", videos)
}
func (d *MongoDBTestSuit) TestAnd() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	filter := bson.A{
		bson.D{bson.E{Key: "uid", Value: 123456}},
		bson.D{bson.E{Key: "vuid", Value: 654321}},
	}
	findRes, err := d.col.Find(ctx, bson.D{bson.E{Key: "$and", Value: filter}})
	assert.NoError(d.T(), err)
	var videos []Video
	err = findRes.All(ctx, &videos)
	assert.NoError(d.T(), err)
	d.T().Log("Inserted IDs: ", videos)
}

func (d *MongoDBTestSuit) TestIn() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	filter := bson.D{bson.E{Key: "uid", Value: bson.D{bson.E{Key: "$in", Value: []int{123456, 1234567}}}}}
	findRes, err := d.col.Find(ctx, filter, options.Find().SetProjection(
		bson.D{bson.E{
			Key:   "uid",
			Value: 123456,
		}}))
	assert.NoError(d.T(), err)
	var videos []Video
	err = findRes.All(ctx, &videos)
	assert.NoError(d.T(), err)
	d.T().Log("Inserted IDs: ", videos)
}
func (d *MongoDBTestSuit) TestIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	rindex, err := d.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{bson.E{"vuid", 1}},
		Options: options.Index().SetUnique(true),
	})
	assert.NoError(d.T(), err)
	d.T().Log("Indexes created: ", rindex)
}
func (d *MongoDBTestSuit) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	filter := bson.A{
		bson.D{bson.E{Key: "uid", Value: 123456}},
		bson.D{bson.E{Key: "uid", Value: 1234567}},
	}
	res, err := d.col.DeleteMany(ctx, bson.D{bson.E{Key: "$or", Value: filter}})
	assert.NoError(d.T(), err)
	d.T().Log("Deleted ID: ", res.DeletedCount)
}
func TestMongoDBQueries(t *testing.T) {
	suite.Run(t, &MongoDBTestSuit{})
}
