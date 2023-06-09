package main

import (
	"fmt"
	"time"
)

const saintPetersburgCityId = "498817"

func kelvinToCelsius(kelvin float64) float64 {
	return kelvin - 273.15
}

func main() {
	fw := getWeatherFromApi()

	weatherTgMessage := fmt.Sprintf(
		"Сегодня в Петербурге %s \n%.2f °C, ощущается как %.2f °C\nВлажность %d%%, скорость ветра %d м/с, облачность %d%%",
		fw.Weather[0].Description, kelvinToCelsius(fw.Main.Temp),
		kelvinToCelsius(fw.Main.FeelsLike), fw.Main.Humidity, fw.Wind.Speed, fw.Clouds.All)

	sendResponses(weatherTgMessage)
	db := createDb("Weather", "WeatherTable")

	t := time.Now()
	if (t.Hour() == 9 || t.Hour() == 21) && (t.Minute() == 0) {
		fmt.Println("Adding to db")
		addToDb(db, "WeatherTable", fw)
	}
}
