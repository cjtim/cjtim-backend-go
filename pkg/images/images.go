package images

import (
	"mime"

	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/cjtim/cjtim-backend-go/pkg/gstorage"
	"github.com/cjtim/cjtim-backend-go/pkg/line"
	"github.com/line/line-bot-sdk-go/linebot"
)

// Upload Image LineUser send from chat
func Add(e *linebot.Event, m *models.Models) (*collections.FileScheama, error) {
	message := e.Message.(*linebot.ImageMessage)
	fileByte, fileType, err := line.GetContent(message.ID)
	if err != nil {
		return nil, err
	}
	ext, err := mime.ExtensionsByType(fileType)
	if ext == nil {
		ext = []string{""}
	}
	if err != nil {
		return nil, err
	}
	Gclient, err := gstorage.GetClient()
	defer Gclient.Client.Close()
	if err != nil {
		return nil, err
	}
	objPath := ("users/" + e.Source.UserID + "/files/" + message.ID + ext[0])
	url, err := Gclient.Upload(objPath, fileByte)
	if err != nil {
		return nil, err
	}
	data := &collections.FileScheama{
		FileName: message.ID + ext[0],
		URL:      url,
		LineUID:  e.Source.UserID,
	}
	_, err = m.InsertOne("files", data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
