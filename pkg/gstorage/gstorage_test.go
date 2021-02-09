package gstorage_test

import (
	"context"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/cjtim/cjtim-backend-go/pkg/gstorage"
	"google.golang.org/api/option"
)

func Test_GstorageList(t *testing.T) {
	ctx := context.Background()
	gClient, err := storage.NewClient(ctx, option.WithCredentialsFile("../../serviceAcc.json"))
	if err != nil {
		t.Fatal(err)
	}
	client := gstorage.Bucket{Client: gClient}
	_, err = client.List("")
	if err != nil {
		t.Fatal(err)
	}
	err = client.Client.Close()
	if err != nil {
		t.Fatal(err)
	}

}
