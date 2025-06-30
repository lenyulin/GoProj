package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestMongodb(t *testing.T) {
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
	insertRes, err := col.InsertOne(ctx, Video{
		Title:   "Test Video",
		Content: "This is a test video",
		Uid:     123456,
		VUid:    1234566,
		Status:  0,
	})
	assert.NoError(t, err)
	oid := insertRes.InsertedID.(primitive.ObjectID)
	fmt.Printf("Inserted ID: %s\n", oid)
	filter := bson.M{"uid": 123456}
	findRes := col.FindOne(ctx, filter)
	if errors.Is(findRes.Err(), mongo.ErrNoDocuments) {
		fmt.Println("No documents found")
	} else {
		assert.NoError(t, findRes.Err())
		var video Video
		err := findRes.Decode(&video)
		fmt.Println(video)
		assert.NoError(t, err)
	}
	updateFilter := bson.D{bson.E{"uid", 123456}}
	set := bson.D{bson.E{Key: "$set", Value: bson.M{
		"title": "New title",
	}}}
	updateRes, err := col.UpdateOne(ctx, updateFilter, set)
	assert.NoError(t, err)
	fmt.Printf("update count%d", updateRes.MatchedCount)

	delFilter := bson.D{bson.E{"uid", 123456}}
	delRes, err := col.DeleteOne(ctx, delFilter)
	assert.NoError(t, err)
	fmt.Printf("delete count%d", delRes.DeletedCount)
}

type Video struct {
	Title   string `bson:"title,omitempty"`
	Content string `bson:"content,omitempty"`
	Uid     int64  `bson:"uid,omitempty"`
	VUid    int64  `bson:"vuid,omitempty"`
	Status  uint8  `bson:"status,omitempty"`
}
