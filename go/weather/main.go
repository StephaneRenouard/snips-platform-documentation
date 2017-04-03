package main

import (
	"os"
	"snips-platform-lambda-samples/go/common/mqtt"
	"snips-platform-lambda-samples/go/weather/conf"
	"encoding/json"
	"snips-platform-lambda-samples/go/weather/openweathermap"
	"strings"
	"snips-platform-lambda-samples/go/common/models"
	"log"
)

func main() {
	conf.Run(os.Args, func(sc conf.WeatherConf) {
		client := mqtt.NewClient(sc.PlatformHost, sc.PlatformPort)
		client.Connect()
		defer client.Disconnect()

		cityChannel := make(chan string)

		client.Subscribe("hermes/Intent/get_weather", func(topic string, payload []byte) {
			value := models.ParseIntent(payload).ParsedSlots[0].Value
			cityChannel <- strings.ToLower(value)
		})

		for {
			m := make(map[string]string)
			m["text"] = openweathermap.GetWeather(sc, <-cityChannel)
			message, err := json.Marshal(m)
			if err != nil {
				log.Fatalf("main: %v", err)
			}
			client.Publish("hermes/Tts", message)
		}
	})
}