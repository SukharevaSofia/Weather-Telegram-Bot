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

const telegramPollTimeout = 3600

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

func sendResponses(weatherMessage string) {
	postUserUrl := fmt.Sprintf(
		"%s/getUpdates?timeout=%d",
		telegramBaseUrl,
		telegramPollTimeout)
	userRequestsRes, err := http.Get(postUserUrl)
	if err != nil {
		fmt.Println("Failed to get users' requests")
		os.Exit(1)
	}
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
