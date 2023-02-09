package main

import (
	"log"
	"io"
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
)

func PrettyString(str string) (string, error) {
    var prettyJSON bytes.Buffer
    if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
        return "", err
    }
    return prettyJSON.String(), nil
}


var bearer = "Bearer " + "xoxb-4734239656646-4734345418358-pLe3taBr8uLiM3R5L5oQhgTi"

func getChannelList() {

    url := "https://slack.com/api/conversations.list"
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("Authorization", bearer)

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
	res, err := PrettyString(string(body))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(res)
}

func getAllUsers() {
	url := "https://slack.com/api/users.list"
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("Authorization", bearer)

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


	res, err := PrettyString(string(body))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(res)
}


func main() {
	getChannelList()
	getAllUsers()
}
