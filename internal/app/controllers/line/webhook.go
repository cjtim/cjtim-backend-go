package line_controllers

import (
	"bytes"
	"net/http"

	"github.com/cjtim/cjtim-backend-go/configs"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/airvisual"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/files"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/line"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

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
	events, err := linebot.ParseRequest(configs.Config.LineChannelSecret, httpReq)
	if err != nil {
		return err
	}
	// Webhook type
	for _, event := range events {
		switch EventMessageType(event) {
		case "location":
			location := event.Message.(*linebot.LocationMessage)
			weatherData, err := airvisual.GetByLocation(location.Latitude, location.Longitude)
			if err != nil {
				zap.L().Error("Line webhook error", zap.String("event", "location"), zap.Error(err))
				return err
			}
			err = line.Reply(event.ReplyToken, []interface{}{line.WeatherFlexMessage(weatherData)})
			if err != nil {
				zap.L().Error("Line webhook error - reply", zap.String("event", "location"), zap.Error(err))
			}
			return err
		case "file":
			message := event.Message.(*linebot.FileMessage)
			return addFile(message.ID, event.Source.UserID, event.ReplyToken)

		case "image":
			message := event.Message.(*linebot.ImageMessage)
			return addFile(message.ID, event.Source.UserID, event.ReplyToken)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func addFile(msgId, userId, replyToken string) error {
	newfile, err := files.Client.AddFromLine(msgId, userId)
	if err != nil {
		zap.L().Error("Line webhook error",
			zap.String("event", "file"),
			zap.String("msgId", msgId),
			zap.String("userId", userId),
			zap.Error(err),
		)
	}
	msgs := []interface{}{
		map[string]interface{}{
			"type": "text",
			"text": newfile.URL.ShortURL,
		}}
	err = line.Reply(replyToken, msgs)
	if err != nil {
		zap.L().Error("error reply msg", zap.Error(err))
	}
	return err
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
