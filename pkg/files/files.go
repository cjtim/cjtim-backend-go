package files

import (
	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/cjtim/cjtim-backend-go/pkg/gstorage"
	"github.com/cjtim/cjtim-backend-go/pkg/rebrandly"
	"go.mongodb.org/mongo-driver/bson"
)

// Add - Upload file and save to DB
func Add(fullFileName string, byteData []byte, lineUID string, m *models.Models) (*collections.FileScheama, error) {
	Gclient, err := gstorage.GetClient()
	defer Gclient.Client.Close()
	if err != nil {
		return nil, err
	}
	objPath := ("users/" + lineUID + "/files/" + fullFileName)
	url, err := Gclient.Upload(objPath, byteData)
	if err != nil {
		return nil, err
	}
	shortURL, err := rebrandly.Add(url)
	if err != nil {
		return nil, err
	}
	data := &collections.FileScheama{
		FileName: fullFileName,
		URL:      *shortURL,
		LineUID:  lineUID,
	}
	_, err = m.InsertOne("files", data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Delete - Remove file from storage and rebrandly
func Delete(fullFileName string, lineUID string, m *models.Models) error {
	file := &collections.FileScheama{}
	err := m.FindOne("files", file, bson.M{"lineUid": lineUID, "fileName": fullFileName})
	if err != nil {
		return err
	}
	// Delete from storage
	gClient, err := gstorage.GetClient()
	defer gClient.Client.Close()
	if err != nil {
		return err
	}
	err = gClient.Delete("users/" + lineUID + "/files/" + fullFileName)
	if err != nil {
		return err
	}
	// Delete Rebrandly URL
	err = rebrandly.Delete(file.URL.RebrandlyID)
	if err != nil {
		return err
	}
	err = m.Destroy("files", bson.M{"fileName": fullFileName, "lineUid": lineUID})
	if err != nil {
		return err
	}
	return nil
}
