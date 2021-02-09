package line_controllers

import (
	"bytes"
	"net/http"
	"os"

	"github.com/cjtim/cjtim-backend-go/models"
	"github.com/cjtim/cjtim-backend-go/pkg/airvisual"
	"github.com/cjtim/cjtim-backend-go/pkg/files"
	"github.com/cjtim/cjtim-backend-go/pkg/images"
	"github.com/cjtim/cjtim-backend-go/pkg/line"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

var _ = godotenv.Load()

// Webhook - for line webhook
func Webhook(c *fiber.Ctx) error {
	// Convert fiber.Ctx to http.Request
	// for linebot.ParseRequest
	httpReq, err := http.NewRequest(c.Method(), c.OriginalURL(), bytes.NewReader(c.Body()))
	if err != nil {
		return err
	}
	for k, v := range utils.HeadersToMapStr(c) {
		httpReq.Header.Set(k, v)
	}
	event, err := linebot.ParseRequest(os.Getenv("LINE_CHANNEL_SECRET"), httpReq)
	if err != nil {
		return err
	}
	// Webhook type
	if len(event) > 0 {
		switch EventMessageType(event[0]) {
		case "location":
			location := event[0].Message.(*linebot.LocationMessage)
			weatherData, err := airvisual.GetByLocation(location.Latitude, location.Longitude)
			if err != nil {
				return err
			}
			err = line.Reply(event[0].ReplyToken,
				[]interface{}{line.WeatherFlexMessage(weatherData)})
			if err != nil {
				return err
			}
			return nil
		case "file":
			message := event[0].Message.(*linebot.FileMessage)
			fileByte, fileType, err := line.GetContent(message.ID)
			if err != nil {
				return err
			}
			ext, err := utils.ContentTypeToExtension(fileType)
			if err != nil {
				return nil
			}
			_, err = files.Add(message.ID+ext[0], fileByte, event[0].Source.UserID,
				c.Locals("db").(*models.Models))
			return err
		case "image":
			_, err = images.Add(event[0], c.Locals("db").(*models.Models))
			return err
		}
	}
	return c.SendStatus(fiber.StatusOK)
}

// EventMessageType - Check event message type
func EventMessageType(e *linebot.Event) linebot.MessageType {
	switch e.Message.(type) {
	case *linebot.TextMessage:
		return linebot.MessageTypeText
	case *linebot.ImageMessage:
		return linebot.MessageTypeImage
	case *linebot.VideoMessage:
		return linebot.MessageTypeVideo
	case *linebot.AudioMessage:
		return linebot.MessageTypeAudio
	case *linebot.FileMessage:
		return linebot.MessageTypeFile
	case *linebot.LocationMessage:
		return linebot.MessageTypeLocation
	case *linebot.StickerMessage:
		return linebot.MessageTypeSticker
	}
	return ""
}
