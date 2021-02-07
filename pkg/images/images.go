package images

import (
	"mime"

	"github.com/cjtim/cjtim-backend-go/datasource/collections"
	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/cjtim/cjtim-backend-go/pkg/files"
	"github.com/cjtim/cjtim-backend-go/pkg/line"
	"github.com/line/line-bot-sdk-go/linebot"
)

// Add - Upload Image and save to DB
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
	return files.Add(message.ID+ext[0], fileByte, e.Source.UserID, m)
}
