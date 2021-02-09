package line

import (
	"fmt"

	"github.com/cjtim/cjtim-backend-go/pkg/airvisual"
)

func WeatherFlexMessage(data *airvisual.AirVisualResponse) map[string]interface{} {
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
	return map[string]interface{}{
		"type":    "flex",
		"altText": "รายงานอากาศ",
		"contents": map[string]interface{}{
			"type": "bubble",
			"size": "kilo",
			"header": map[string]interface{}{
				"type":   "box",
				"layout": "vertical",
				"contents": []map[string]interface{}{
					{
						"type":    "text",
						"text":    headerMessage,
						"color":   "#ffffff",
						"align":   "start",
						"size":    "xl",
						"gravity": "center",
					},
					{
						"type":    "text",
						"text":    AQIValue,
						"color":   "#ffffff",
						"align":   "start",
						"size":    "xxl",
						"gravity": "center",
						"margin":  "lg",
					},
				},
				"backgroundColor": bgColor,
				"paddingTop":      "19px",
				"paddingAll":      "12px",
				"paddingBottom":   "16px",
				"action": map[string]interface{}{
					"type":  "uri",
					"label": "action",
					"uri":   "https://www.iqair.com/th-en/thailand/bangkok/phaya-thai",
					"altUri": map[string]interface{}{
						"desktop": "https://www.iqair.com/th-en/thailand/bangkok/phaya-thai",
					},
				},
			},
			"body": map[string]interface{}{
				"type":   "box",
				"layout": "vertical",
				"contents": []map[string]interface{}{
					{
						"type":   "box",
						"layout": "horizontal",
						"contents": []map[string]interface{}{
							{
								"type":  "text",
								"text":  innerMessage,
								"color": "#000000",
								"size":  "lg",
								"wrap":  true,
							},
						},
						"flex": 1,
					},
				},
				"spacing":    "lg",
				"paddingAll": "12px",
			},
			"styles": map[string]interface{}{
				"footer": map[string]interface{}{
					"separator": false,
				},
			},
		},
	}
}
