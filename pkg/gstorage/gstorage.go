package gstorage

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/url"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
	"github.com/cjtim/cjtim-backend-go/config"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

var projectID = config.Config.GProjectID
var bucketName = config.Config.GBucketName

type Bucket struct {
	Client *storage.Client
}

func GetClient() (Bucket, error) {
	ctx := context.Background()
	var client, err = storage.NewClient(ctx, option.WithCredentialsFile("./config/serviceAcc.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return Bucket{}, err
	}
	return Bucket{Client: client}, nil
}

func (s *Bucket) Upload(path string, byteData []byte) (string, error) {
	downloadToken := uuid.New().String()
	bucket := s.Client.Bucket(bucketName)
	ctx := context.Background()
	wc := bucket.Object(path).NewWriter(ctx)
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

func (s *Bucket) Delete(path string) error {
	return s.Client.Bucket(bucketName).Object(path).Delete(context.TODO())
}

func (s *Bucket) List(filename string) ([]string, error) {
	query := &storage.Query{StartOffset: filename}
	var names []string
	it := s.Client.Bucket(bucketName).Objects(context.TODO(), query)
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
