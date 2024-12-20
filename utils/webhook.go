package utils

import (
	"bombelaio-keydrop-golang/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func sendDiscordWebhook(user Users, prize models.PrizesWinner, giveawayId string) {
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title": "We have a winner!",
				"color": 0x00FF00, // green color
				"thumbnail": map[string]interface{}{
					"url": fmt.Sprintf("https://cdn.key-drop.com//%v", prize.ItemImg),
				},
				"fields": []map[string]interface{}{
					{
						"name":  "**Our winner is!**",
						"value": fmt.Sprintf("%v, has won a giveaway!", user.Name),
					},
					{
						"name":  "Prize:",
						"value": fmt.Sprintf("||%v - %s %v, that is worth: %g %s||", prize.Title, prize.Subtitle, prize.Condition, prize.Price, prize.Currency),
					},
					{
						"name":  "Giveaway url:",
						"value": fmt.Sprintf("[Key-Drop](https://key-drop.com/giveaways/keydrop/%s)", giveawayId),
					},
				},

				"footer": map[string]interface{}{
					"text":     "Sent from my Golang application",
					"icon_url": fmt.Sprintf("https://cdn.key-drop.com//%v", prize.ItemImg),
				},
				"timestamp": time.Now().UTC().Format("2006-01-02T15:04:05.999Z"),
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", DiscordWebhook, bytes.NewBuffer(jsonPayload))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode > 400 && resp.StatusCode < 500 {
		Sleep(10000)
		Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error sending webhook: %v", resp.StatusCode))
		sendDiscordWebhook(user, prize, giveawayId)
	}

}
