package files

import (
	"context"

	"github.com/cjtim/cjtim-backend-go/pkg/gstorage"
	"github.com/cjtim/cjtim-backend-go/pkg/line"
	"github.com/cjtim/cjtim-backend-go/pkg/rebrandly"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"github.com/cjtim/cjtim-backend-go/repository"
	"go.mongodb.org/mongo-driver/bson"
)

// Add - Upload file and save to DB
func Add(fullFileName string, byteData []byte, lineUID string) (*repository.FileScheama, error) {
	client, err := gstorage.GetClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	objPath := ("users/" + lineUID + "/files/" + fullFileName)
	url, err := client.Upload(objPath, byteData)
	if err != nil {
		return nil, err
	}
	shortURL, err := rebrandly.Add(url)
	if err != nil {
		return nil, err
	}
	data := &repository.FileScheama{
		FileName: fullFileName,
		URL:      *shortURL,
		LineUID:  lineUID,
	}
	_, err = repository.FileRepo.InsertOne(context.TODO(), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AddFromLine - Contents being upload via line chat will parse here
func AddFromLine(messageID string, lineUID string) (*repository.FileScheama, error) {
	fileByte, fileType, err := line.GetContent(messageID)
	if err != nil {
		return nil, err
	}
	ext, err := utils.ContentTypeToExtension(fileType)
	if ext == nil {
		ext = []string{""}
	}
	if err != nil {
		return nil, err
	}
	return Add(messageID+ext[0], fileByte, lineUID)
}

// Delete - Remove file from storage and rebrandly
func Delete(fullFileName string, lineUID string) error {
	file := repository.FileScheama{}
	err := repository.FileRepo.FindOne(&file, bson.M{"lineUid": lineUID, "fileName": fullFileName})
	if err != nil {
		return err
	}
	// Delete from storage
	gClient, err := gstorage.GetClient()
	if err != nil {
		return err
	}
	defer gClient.Close()

	err = gClient.Delete("users/" + lineUID + "/files/" + fullFileName)
	if err != nil {
		return err
	}
	// Delete Rebrandly URL
	err = rebrandly.Delete(file.URL.RebrandlyID)
	if err != nil {
		return err
	}
	_, err = repository.FileRepo.DeleteOne(context.TODO(), bson.M{"fileName": fullFileName, "lineUid": lineUID})
	if err != nil {
		return err
	}
	return nil
}
