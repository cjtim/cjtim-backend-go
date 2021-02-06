package line

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

var _ = godotenv.Load()
var restyClient = resty.New()
var LineBot, err = linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))

func LineIsTokenValid(accToken string) error {
	resp, err := restyClient.R().SetQueryParam("access_token", accToken).Get("https://api.line.me/oauth2/v2.1/verify")
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New(string(resp.Body()))
	}
	return nil
}

func LineGetProfile(accToken string) (*linebot.UserProfileResponse, error) {
	resp, err := restyClient.R().SetHeader("Authorization", "Bearer "+accToken).Get("https://api.line.me/v2/profile")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New(string(resp.Body()))
	}
	profile := &linebot.UserProfileResponse{}
	if resp.StatusCode() == 200 {
		body := resp.Body()
		err := json.Unmarshal(body, &profile)
		if err != nil {
			return nil, err
		}
		return profile, nil
	}
	return nil, errors.New(string(resp.Body()))
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
