package utils

import (
	"bytes"
	"html/template"
	"math"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
)

func Round(num float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Round(num*shift) / shift
}

func ParseForecast(forecast model.Forecast) (string, error) {
	htmlTemplate := `
		<b>Weather Forecast for {{.Name}}</b>
		Temperature: <b>{{.Main.Temp}}</b>
		Feels like: <b>{{.Main.FeelsLike}}</b>
		Min temp: <b>{{.Main.TempMin}}</b>
		Max temp: <b>{{.Main.TempMax}}</b>
		Pressure: <b>{{.Main.Pressure}}</b>
	`
	tmpl, err := template.New("weatherTemplate").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	var tplBuffer bytes.Buffer
	if err := tmpl.Execute(&tplBuffer, forecast); err != nil {
		return "", err
	}

	return tplBuffer.String(), nil
}
