package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/cjtim/cjtim-backend-go/datasource"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Models -- center for invoke any operation on Model
type Models struct {
	Client *mongo.Client
}

// GetModels by pass Mongo Client pointer
func GetModels(c *mongo.Client) *Models {
	models := &Models{}
	if c == nil {
		DBchannel := make(chan *mongo.Client)
		go datasource.MongoClient(DBchannel) // GoRoutine connectDB
		models.Client = <-DBchannel
		return models
	}
	models.Client = c
	return models
}

func (s *Models) FindAll(collectionName string, results interface{}, filter bson.M) error {
	var tmp []bson.M
	collection := s.Client.Database(os.Getenv("MONGO_DB")).Collection(collectionName)
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer cursor.Close(context.TODO())
	if cursor.All(context.TODO(), &tmp) != nil {
		log.Fatal(err)
	}
	jbyte, _ := json.Marshal(tmp)
	json.Unmarshal(jbyte, &results)
	return nil
}

func (s *Models) FindOne(collectionName string, results interface{}, filter bson.M) error {
	var tmp bson.M
	collection := s.Client.Database(os.Getenv("MONGO_DB")).Collection(collectionName)
	resp := collection.FindOne(context.TODO(), filter)
	if err := resp.Decode(&tmp); err != nil {
		return err
	}
	jbyte, err := json.Marshal(&tmp)
	if err != nil {
		return nil
	}
	json.Unmarshal(jbyte, &results)
	return nil
}

func (s *Models) InsertOne(collectionName string, data interface{}) (interface{}, error) {
	collection := s.Client.Database(os.Getenv("MONGO_DB")).Collection(collectionName)
	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return insertResult.InsertedID, nil
}

func (s *Models) Update(collectionName string, data interface{}, filter interface{}) (*mongo.UpdateResult, error) {
	collection := s.Client.Database(os.Getenv("MONGO_DB")).Collection(collectionName)
	return collection.UpdateMany(context.TODO(), filter, bson.M{
		"$set": data,
	})
}

func (s *Models) Destroy(collectionName string, filter bson.M) error {
	collection := s.Client.Database(os.Getenv("MONGO_DB")).Collection(collectionName)
	if filter == nil {
		return errors.New("no filter apply!")
	}
	collection.DeleteMany(context.TODO(), filter)
	return nil
}
