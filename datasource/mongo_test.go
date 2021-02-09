package datasource_test

import (
	"context"
	"testing"

	"github.com/cjtim/cjtim-backend-go/models"
)

func Test_Database(t *testing.T) {
	m, err := models.GetModels(nil)
	if err != nil {
		t.Fatal(err)
	}
	if m.Client.Disconnect(context.TODO()) != nil {
		t.Fatal("Failed Disconect DB!")
	}
}
