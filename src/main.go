package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct {
  Port string
}

type TelegramConfig struct {
	BotToken string
	ChannelID string
}


type Payload struct {
	Message string `json:"message"`
}

type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text string `json:"text"`
}

func (t *TelegramConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// read the request body and parse the json
	var p Payload
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// convert the body to string
	fmt.Println("Body: ", string(body))
	err = json.Unmarshal([]byte(body), &p)
	if err != nil {
		http.Error(w, "Please provide message in body", http.StatusBadRequest)
		return
	}
	// send the message to telegram channel
	// https://api.telegram.org/bot{bottoken}/sendMessage

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.BotToken)
	tm := TelegramMessage{
		ChatID: t.ChannelID,
		Text: p.Message,
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	defer r.Body.Close()
	body, err = json.Marshal(tm)
	if err != nil {
		http.Error(w, "Error parsing request", http.StatusInternalServerError)
		return
	}
	res, err1 := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err1 != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	fmt.Println("Client url: ", url)
	fmt.Println("Client body: ", body)
	fmt.Println("Client res: ", res)
	fmt.Println("Client err: ", err1)

	if res.StatusCode != http.StatusOK {
		http.Error(w, "Error sending message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message sent"))
}


func (a *App) Start() {
  addr := fmt.Sprintf(":%s", a.Port)
	telegramConfig := TelegramConfig{
		BotToken: env("BOT_TOKEN", ""),
		ChannelID: env("CHANNEL_ID", ""),
	}
	log.Printf("Telegram Config: %+v", telegramConfig)
	http.Handle("/", &telegramConfig)
  log.Printf("Starting app on %s", addr)
  log.Fatal(http.ListenAndServe(addr, nil))
}


func env(key, defaultValue string) string {
  val, ok := os.LookupEnv(key)
  if !ok {
    return defaultValue
  }
  return val
}

func main() {
  server := App{
    Port: env("PORT", "9999"),
  }
  server.Start()
}

