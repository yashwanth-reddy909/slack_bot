package services

import(
	"log"
	"net/http"
	"os"
	"io"
)



func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	url := "https://slack.com/api/users.list"
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("Authorization", "Bearer " + os.Getenv("BOT_TOKEN"))
    http_client := &http.Client{}
    resp, err := http_client.Do(req)
    if err != nil {
        log.Println("Error on response.\n[ERROR] -", err)
    }

    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error while reading the response bytes:", err)
    }
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func GetChannelList(w http.ResponseWriter, r *http.Request) {

	url := "https://slack.com/api/conversations.list"
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("Authorization","Bearer "+ os.Getenv("BOT_TOKEN"))
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error on response.\n[ERROR] -", err)
    }
	defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error while reading the response bytes:", err)
    }
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}