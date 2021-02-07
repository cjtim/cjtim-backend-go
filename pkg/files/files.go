package files

import (
	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/cjtim/cjtim-backend-go/pkg/gstorage"
	"github.com/cjtim/cjtim-backend-go/pkg/rebrandly"
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
