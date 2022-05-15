package main

import (
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Updates struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

type Update struct {
	ID       int
	Username string
	Text     string
	Date     int
}

func main() {
	b, errr := os.ReadFile("./config/config.yaml")
	if errr != nil {
		log.Fatal(errr)
	}

	c, errc := config.ParseConfig(b)

	if errc != nil {
		log.Fatal(errc)
	}

	cl := http.Client{Timeout: 10 * time.Second}

	getUrl := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", c.ApiKeys.Telegram)

	mmap := make(map[int]Update)

	for {
		r, err := cl.Get(getUrl)
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(r.Body)
		updates := Updates{}
		errj := json.Unmarshal(b, &updates)
		if errj != nil {
			log.Fatal(errj)
		}
		if updates.Ok {
			for _, u := range updates.Result {
				_, ok := mmap[u.UpdateID]
				if !ok {
					upd := Update{}

					upd.ID = u.Message.From.ID
					upd.Text = u.Message.Text
					upd.Username = u.Message.From.Username
					upd.Date = u.Message.Date

					mmap[u.UpdateID] = upd

					fmt.Printf("Get %+v\n", upd)

					sendUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s",
						c.ApiKeys.Telegram, upd.ID, "I get <"+upd.Text+">")
					_, errp := cl.Post(sendUrl, "text/plain", nil)
					if errp != nil {
						fmt.Printf("Error = %v\n", errp)
					} else {
						fmt.Printf("I send smth\n")
					}
				}
			}
		} else {
			fmt.Printf("Get nothing")
		}
		_ = r.Body.Close()
		time.Sleep(1000 * time.Millisecond)
	}
}
