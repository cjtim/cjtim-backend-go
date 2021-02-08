package line

import (
	"fmt"

	"github.com/cjtim/cjtim-backend-go/pkg/airvisual"
	"github.com/line/line-bot-sdk-go/linebot"
)

func WeatherFlexMessage(data *airvisual.AirVisualResponse) *linebot.FlexMessage {
	var headerMessage, bgColor string
	AQIValue := fmt.Sprintf("%v", data.Data.Current.Pollution.AQIUS) + " AQIUS"
	innerMessage := fmt.Sprintf(
		`City: %v
AQI us: %v
AQI cn: %v
Temperature:%vº
Humidity: %v
Wind Speed: %v(m/s)`,
		data.Data.City,
		data.Data.Current.Pollution.AQIUS,
		data.Data.Current.Pollution.AQICN,
		data.Data.Current.Weather.Temperature,
		data.Data.Current.Weather.Humidity,
		data.Data.Current.Weather.WindSpeed,
	)
	if data.Data.Current.Pollution.AQIUS <= 50 {
		headerMessage = "อากาศดีจัง"
		bgColor = `#27ACB2`
	} else if data.Data.Current.Pollution.AQIUS > 50 {
		headerMessage = "อากาศดีพอใช้"
		bgColor = `#E9AF29`
	} else if data.Data.Current.Pollution.AQIUS > 100 {
		headerMessage = "อากาศไม่ค่อยดีนะ"
		bgColor = `#E9632D`
	} else {
		headerMessage = "อากาศแย่มากเลย🤢"
		bgColor = `#FF6B6E`
	}
	return &linebot.FlexMessage{
		AltText: "รายงานอากาศ",
		Contents: &linebot.BubbleContainer{
			Type: "bubble",
			Size: "kilo",
			Header: &linebot.BoxComponent{
				Type:   "box",
				Layout: "vertical",
				Contents: []linebot.FlexComponent{
					&linebot.TextComponent{
						Type:    "text",
						Text:    headerMessage,
						Color:   "#ffffff",
						Align:   "start",
						Size:    "xl",
						Gravity: "center",
					},
					&linebot.TextComponent{
						Type:    "text",
						Text:    AQIValue,
						Color:   "#ffffff",
						Align:   "start",
						Size:    "xxl",
						Gravity: "center",
						Margin:  "lg",
					},
				},
				BackgroundColor: bgColor,
			},
			Body: &linebot.BoxComponent{
				Type:   "box",
				Layout: "vertical",
				Contents: []linebot.FlexComponent{
					&linebot.BoxComponent{
						Type:   "box",
						Layout: "horizontal",
						Contents: []linebot.FlexComponent{
							&linebot.TextComponent{
								Type:  "text",
								Text:  innerMessage,
								Color: "#000000",
								Size:  "lg",
								Wrap:  true,
							},
						},
						Flex: linebot.IntPtr(1),
					},
				},
				Spacing: linebot.FlexComponentSpacingTypeLg,
			},
			Styles: &linebot.BubbleStyle{
				Footer: &linebot.BlockStyle{
					Separator: false,
				},
			},
		},
	}
}
