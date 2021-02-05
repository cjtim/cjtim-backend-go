package line

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/go-resty/resty/v2"
	"github.com/line/line-bot-sdk-go/linebot"
)

var restyClient = resty.New()
var LineBot, err = linebot.New("<channel secret>", "<channel access token>")

func LineIsTokenValid(accToken string) error {
	resp, err := restyClient.R().SetQueryParam("access_token", accToken).Get("https://api.line.me/oauth2/v2.1/verify")
	fmt.Println(resp.StatusCode())
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New(string(resp.Body()))
	}
	return nil
}

func LineGetProfile(accToken string) {
	resp, err := restyClient.R().SetHeader("Authorization", accToken).Get("https://api.line.me/v2/profile")
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode() == 200 {
		fmt.Println(string(resp.Body()))
	}
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
