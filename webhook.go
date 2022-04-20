package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func SendSuccess(vanity string, time string) {
	embed := []byte(fmt.Sprintf(`{
		"content": "@everyone",
		"embeds": [
		  {
			"title": "SNIPED!",
			"color": 65280,
			"fields": [
			  {
				"name": "Vanity",
				"value": "%s"
			  },
			  {
				"name": "Time Taken to attempt",
				"value": "%s"
			  }
			],
			"thumbnail": {
			  "url": "https://m.media-amazon.com/images/I/71+XQ4SFqnL._AC_SL1500_.jpg"
			}
		  }
		],
		"attachments": []
	  }`, vanity, time))

	resp, err := http.Post(Config.Main.Webhook, "application/json", bytes.NewBuffer(embed))

	if err != nil {
		fmt.Println(err)
		SendSuccess(vanity, time)
	}

	defer resp.Body.Close()
}

func SendRatelimit(vanity string, time string) {
	embed := []byte(fmt.Sprintf(`{
		"content": "@everyone",
		"embeds": [
		  {
			"title": "RATELIMITED",
			"color": 16763904,
			"fields": [
			  {
				"name": "Vanity",
				"value": "%s"
			  },
			  {
				"name": "Time Taken to attempt",
				"value": "%s"
			  }
			],
			"thumbnail": {
			  "url": "https://www.cambridge.org/elt/blog/wp-content/uploads/2019/07/Sad-Face-Emoji-480x480.png"
			}
		  }
		],
		"attachments": []
	  }`, vanity, time))

	resp, err := http.Post(Config.Main.Webhook, "application/json", bytes.NewBuffer(embed))

	if err != nil {
		fmt.Println(err)
		SendSuccess(vanity, time)
	}

	defer resp.Body.Close()
}

func SendFail(vanity string, time string, statuscode string) {
	embed := []byte(fmt.Sprintf(`{
		"content": "@everyone",
		"embeds": [
		  {
			"title": "FAILED TO SNIPE",
			"description": "Better luck next time",
			"color": 16763904,
			16711680
			"fields": [
			  {
				"name": "Vanity",
				"value": "%s"
			  },
			  {
				"name": "Time Taken to attempt",
				"value": "%s"
			  },
			  {
				"name": "Status Code",
				"value": "%s"
			  }
			],
			"thumbnail": {
			  "url": "https://www.cambridge.org/elt/blog/wp-content/uploads/2019/07/Sad-Face-Emoji-480x480.png"
			}
		  }
		],
		"attachments": []
	  }`, vanity, time, statuscode))

	resp, err := http.Post(Config.Main.Webhook, "application/json", bytes.NewBuffer(embed))

	if err != nil {
		fmt.Println(err)
		SendSuccess(vanity, time)
	}

	defer resp.Body.Close()
}
