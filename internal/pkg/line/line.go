package line

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/cjtim/cjtim-backend-go/configs"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/utils"
	"github.com/line/line-bot-sdk-go/linebot"
)

var LineBot, _ = linebot.New(configs.Config.LineChannelSecret, configs.Config.LineChannelAccessToken)

func LineIsTokenValid(accToken string) error {

	resp, body, err := utils.Http(&utils.HttpReq{
		Method: http.MethodGet,
		URL:    "https://api.line.me/oauth2/v2.1/verify",
		Querys: map[string]string{
			"access_token": accToken,
		},
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(string(body))
	}
	return nil
}

func LineGetProfile(accToken string) (*linebot.UserProfileResponse, error) {
	resp, body, err := utils.Http(&utils.HttpReq{
		Method: http.MethodGet,
		URL:    "https://api.line.me/v2/profile",
		Headers: map[string]string{
			"Authorization": "Bearer " + accToken,
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(string(body))
	}
	profile := &linebot.UserProfileResponse{}
	if resp.StatusCode == 200 {
		body := body
		err := json.Unmarshal(body, &profile)
		if err != nil {
			return nil, err
		}
		return profile, nil
	}
	return nil, errors.New(string(body))
}

func GetContent(messageID string) ([]byte, string, error) {
	content, err := LineBot.GetMessageContent(messageID).Do()
	if err != nil {
		return nil, "", err
	}
	fileByte, err := ioutil.ReadAll(content.Content)
	if err != nil {
		return nil, "", err
	}
	return fileByte, content.ContentType, nil
}

func Reply(replayToken string, msgs map[string]interface{}) error {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + configs.Config.LineChannelAccessToken,
	}
	resp, body, err := utils.Http(&utils.HttpReq{
		Method:  http.MethodPost,
		URL:     "https://api.line.me/v2/bot/message/reply",
		Headers: headers,
		Body: map[string]interface{}{
			"replyToken": replayToken,
			"messages":   msgs,
		},
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(string(body))
	}
	return nil
}

func Broadcast(msgs map[string]interface{}) error {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + configs.Config.LineChannelAccessToken,
	}
	resp, body, err := utils.Http(&utils.HttpReq{
		Method:  http.MethodPost,
		URL:     "https://api.line.me/v2/bot/message/broadcast",
		Headers: headers,
		Body: map[string]interface{}{
			"messages": msgs,
		},
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(string(body))
	}
	return nil
}
