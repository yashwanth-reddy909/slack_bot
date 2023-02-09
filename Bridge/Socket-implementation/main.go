package main

import (
	"Bridge/connect"
	"Bridge/functions"
	"Bridge/services"

	"fmt"
	"net/http"
	"github.com/slack-go/slack"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appToken := os.Getenv("APP_TOKEN")
	botToken := os.Getenv("BOT_TOKEN")
	botId := os.Getenv("BOT_ID")
	client := connect.Connect(appToken, botToken)
	questions := []string{"how are you?", "what did you do?", "what will you do today"}

	Answers := map[string][]string{}
	go functions.UserQA(client, questions, Answers, botId)
	client.Client.PostMessage("U04N13TA8R4", slack.MsgOptionText(questions[0], false))
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	fmt.Println(string(colorRed), "In main:", string(colorReset))
	go client.Run()

	// server api
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/getUsers", services.GetAllUsers)
	r.Get("/getChannelList", services.GetChannelList)
	http.ListenAndServe(":3000", r)
}
