package line_webhook

import (
	"bytes"
	"net/http"
	"os"

	"github.com/cjtim/cjtim-backend-go/pkg/files"
	"github.com/cjtim/cjtim-backend-go/pkg/images"
	"github.com/cjtim/cjtim-backend-go/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

var _ = godotenv.Load()

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
		case "file":
			return files.Add(c, event[0])
		case "image":
			return images.Add(c, event[0])
		}
	}
	return c.SendStatus(fiber.StatusOK)
}

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
