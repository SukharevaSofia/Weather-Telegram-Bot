package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var telegramBaseUrl = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_API_KEY")
var WeatherApiKey = os.Getenv("OPENWEATHER_API_KEY")

const saintPetersburgCityId = "498817"
const telegramPollTimeout = 3600

func kelvinToCelsius(kelvin float64) float64 {
	return kelvin - 273.15
}

func main() {
	type TelegramWeatherRequest struct {
		Ok     bool `json:"ok"`
		Result []struct {
			UpdateId int `json:"update_id"`
			Message  struct {
				MessageId int `json:"message_id"`
				From      struct {
					Id           int    `json:"id"`
					IsBot        bool   `json:"is_bot"`
					FirstName    string `json:"first_name"`
					Username     string `json:"username"`
					LanguageCode string `json:"language_code"`
				} `json:"from"`
				Chat struct {
					Id        int    `json:"id"`
					FirstName string `json:"first_name"`
					Username  string `json:"username"`
					Type      string `json:"type"`
				} `json:"chat"`
				Date     int    `json:"date"`
				Text     string `json:"text"`
				Entities []struct {
					Offset int    `json:"offset"`
					Length int    `json:"length"`
					Type   string `json:"type"`
				} `json:"entities,omitempty"`
			} `json:"message"`
		} `json:"result"`
	}

	type WeatherResponse struct {
		Coord struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		} `json:"coord"`
		Weather []struct {
			Id          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Base string `json:"base"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		} `json:"main"`
		Visibility int `json:"visibility"`
		Wind       struct {
			Speed int `json:"speed"`
			Deg   int `json:"deg"`
		} `json:"wind"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Dt  int `json:"dt"`
		Sys struct {
			Type    int    `json:"type"`
			Id      int    `json:"id"`
			Country string `json:"country"`
			Sunrise int64  `json:"sunrise"`
			Sunset  int64  `json:"sunset"`
		} `json:"sys"`
		Timezone int    `json:"timezone"`
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Cod      int    `json:"cod"`
	}

	//OpenWeather API request
	weatherRes, err := http.Get(
		fmt.Sprintf(
			"https://api.openweathermap.org/data/2.5/weather?id=%s&lang=ru&appid=%s",
			saintPetersburgCityId,
			WeatherApiKey))
	if err != nil {
		fmt.Printf("OpenWeather request failed: %s\n", err)
		os.Exit(1)
	}
	weatherBody, err := io.ReadAll(weatherRes.Body)
	var fw WeatherResponse
	_ = json.Unmarshal(weatherBody, &fw)

	//TODO: forming a good-looking response with good-looking numbers

	weatherMessage := fmt.Sprintf(
		"Сегодня в Петербурге %s \n%.2f °C, ощущается как %.2f °C\nВлажность %d%%, скорость ветра %d м/с, облачность %d%%",
		fw.Weather[0].Description, kelvinToCelsius(fw.Main.Temp),
		kelvinToCelsius(fw.Main.FeelsLike), fw.Main.Humidity, fw.Wind.Speed, fw.Clouds.All)

	//Getting requests from users
	postUserUrl := fmt.Sprintf(
		"%s/getUpdates?timeout=%d",
		telegramBaseUrl,
		telegramPollTimeout)
	userRequestsRes, err := http.Get(postUserUrl)
	userRequestsBody, _ := io.ReadAll(userRequestsRes.Body)
	var requestList TelegramWeatherRequest
	_ = json.Unmarshal(userRequestsBody, &requestList)

	var offset = 0
	for _, value := range requestList.Result {
		if value.UpdateId > offset {
			offset = value.UpdateId
		}
		postSendWeatherUrl := fmt.Sprintf("%s/sendMessage", telegramBaseUrl)
		postSendWeatherData := url.Values{
			"chat_id": {fmt.Sprintf("%d", value.Message.From.Id)},
			"text":    {weatherMessage}}
		_, _ = http.PostForm(postSendWeatherUrl, postSendWeatherData)
	}
	// drop processed messages
	if offset > 0 {
		_, _ = http.Get(fmt.Sprintf("%s/getUpdates?offset=%d&limit=1", telegramBaseUrl, offset+1))
	}
}
