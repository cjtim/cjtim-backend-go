package files

import (
	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/cjtim/cjtim-backend-go/pkg/gstorage"
)

// Add - Upload file from LineUser chat, from web
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
	data := &collections.FileScheama{
		FileName: fullFileName,
		URL:      url,
		LineUID:  lineUID,
	}
	_, err = m.InsertOne("files", data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
