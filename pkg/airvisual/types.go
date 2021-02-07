package airvisual

type location struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
type forcasts struct {
	Timestamp       string  `json:"ts"`
	AQIUS           int     `json:"aqius"`
	AQICN           int     `json:"aqicn"`
	Temperature     float64 `json:"tp"`
	TemperatureMin  float64 `json:"tp_min"`
	PressureHPa     float64 `json:"pr"`
	Humidity        float64 `json:"hu"`
	WindSpeed       float64 `json:"ws"`
	WindDirection   float64 `json:"wd"`
	WeatherIconCode string  `json:"ic"`
}
type currentWeather struct {
	Timestamp       string  `json:"ts"`
	Temperature     float64 `json:"tp"`
	PressureHPa     float64 `json:"pr"`
	Humidity        float64 `json:"hu"`
	WindSpeed       float64 `json:"ws"`
	WindDirection   float64 `json:"wd"`
	WeatherIconCode string  `json:"ic"`
}
type currentPollution struct {
	Timestamp   string            `json:"ts"`
	AQIUS       int               `json:"aqius"`
	MainUS      string            `json:"mainus"`
	AQICN       int               `json:"aqicn"`
	MainCN      string            `json:"maincn"`
	PM25        pollutionUnitData `json:"p2"`
	PM10        pollutionUnitData `json:"p1"`
	OzoneO3     pollutionUnitData `json:"o3"`
	NitrogenNO2 pollutionUnitData `json:"n2"`
	Sulfur      pollutionUnitData `json:"s2"`
	Carbon      pollutionUnitData `json:"co"`
}
type pollutionUnitData struct {
	Concentration float64 `json:"conc"`
	AQIUS         float64 `json:"aqius"`
	AQICN         float64 `json:"aqicn"`
}
type current struct {
	Weather   currentWeather   `json:"weather"`
	Pollution currentPollution `json:"pollution"`
}
type history struct {
	Weather   []currentWeather   `json:"weather"`
	Pollution []currentPollution `json:"pollution"`
}
type units struct {
	PM25        float64 `json:"p2"`
	PM10        float64 `json:"p1"`
	OzoneO3     float64 `json:"o3"`
	NitrogenNO2 float64 `json:"n2"`
	Sulfur      float64 `json:"s2"`
	Carbon      float64 `json:"co"`
}

type data struct {
	Name      string     `json:"name"`
	City      string     `json:"city"`
	State     string     `json:"state"`
	Country   string     `json:"country"`
	Location  location   `json:"location"`
	Forecasts []forcasts `json:"forecasts"`
	Current   current    `json:"current"`
	History   history    `json:"history"`
	Units     units      `json:"units"`
}

// AirVisualResponse - Response from https://api.airvisual.com/v2/city
type AirVisualResponse struct {
	Status string `json:"status"`
	Data   data   `json:"data"`
}
