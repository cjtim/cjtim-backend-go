package models_test

import (
	"context"
	"testing"

	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Models(t *testing.T) {
	m, err := models.GetModels(nil)
	if err != nil {
		t.Fatal(err)
	}

	newData := bson.M{
		"test": "123567",
	}
	id, err := m.InsertOne("unit_test", newData)
	if err != nil {
		t.Fatal("Failed Insert DB!")
	}
	actual := &collections.UnitTestSchema{}
	m.FindOne("unit_test", &actual, newData)
	if actual.Test != newData["test"] {
		t.Fatal("Failed FindOne DB!")
	}
	if m.Destroy("unit_test", bson.M{"_id": id}) != nil {
		t.Fatal("Failed Destroy DB!")
	}
	if m.Client.Disconnect(context.TODO()) != nil {
		t.Fatal("Failed Disconect DB!")
	}

}
