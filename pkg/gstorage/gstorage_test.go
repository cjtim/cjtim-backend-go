package gstorage_test

import (
	"context"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/cjtim/cjtim-backend-go/pkg/gstorage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

func Test_GstorageList(t *testing.T) {
	ctx := context.Background()
	gClient, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		t.Fatal(err)
	}
	client := gstorage.Bucket{Client: gClient}
	_, err = client.List("")
	assert.Error(t, storage.ErrBucketNotExist)
}
