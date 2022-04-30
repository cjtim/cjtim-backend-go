package gstorage

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/url"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
	"github.com/cjtim/cjtim-backend-go/configs"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

var bucketName = configs.Config.GBucketName

type c struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

func GetClient() (*c, error) {
	var client, err = storage.NewClient(
		context.TODO(),
		option.WithCredentialsFile(configs.Config.GCLOUD_CREDENTIAL),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}
	return &c{
		client: client,
		bucket: client.Bucket(bucketName),
	}, nil
}

func (c *c) Upload(path string, byteData []byte) (string, error) {
	downloadToken := uuid.New().String()
	wc := c.bucket.Object(path).NewWriter(context.TODO())
	wc.Metadata = map[string]string{
		"firebaseStorageDownloadTokens": downloadToken,
	}
	data := bytes.NewReader(byteData)
	if _, err := io.Copy(wc, data); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	downloadURL := ("https://firebasestorage.googleapis.com/v0/b/" + bucketName + "/o/" +
		url.QueryEscape(wc.Name) + "?alt=media&token=" + downloadToken)
	return downloadURL, nil
}

func (c *c) Delete(path string) error {
	return c.client.Bucket(bucketName).Object(path).Delete(context.TODO())
}

func (c *c) List(filename string) ([]string, error) {
	query := &storage.Query{StartOffset: filename}
	var names []string
	it := c.bucket.Objects(context.TODO(), query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		names = append(names, attrs.Name)
	}
	return names, nil
}

func (c *c) Close() error {
	return c.client.Close()
}
