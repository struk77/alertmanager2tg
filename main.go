package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/prometheus/alertmanager/template"
)

type TgBot struct {
	API    string `json:"api"`
	ChatID string `json:"chat_id"`
}

func sendAlarm(text string, tgbot *TgBot) {
	apiURL := "https://api.telegram.org/bot" + tgbot.API + "/sendMessage"
	client := &http.Client{}
	form := url.Values{}
	form.Add("chat_id", tgbot.ChatID)
	form.Add("text", text)
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	log.Printf("%s: %s\n", text, resp.Status)
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Println("Incoming Request:", r.Method)

	data := template.Data{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v\n", data)
	log.Printf("Alerts: GroupLabels=%v, CommonLabels=%v", data.GroupLabels, data.CommonLabels)
	tgBotFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully opened tgbot.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer tgBotFile.Close()

	byteValue, _ := ioutil.ReadAll(tgBotFile)

	var tgbot TgBot

	err = json.Unmarshal(byteValue, &tgbot)
	if err != nil {
		log.Fatalln(err)
	}

	for _, alert := range data.Alerts {
		//log.Printf("Alert: status=%s,Labels=%v,Annotations=%v", alert.Status, alert.Labels, alert.Annotations)
		greenHeart := '\U0001F49A'
		brokenHeart := '\U0001F494'
		text := ""
		switch alert.Status {
		case "firing":
			text += fmt.Sprintf("%s %s\n", string(brokenHeart), alert.Annotations["description"])
		case "resolved":
			text += fmt.Sprintf("%s OK!\n%s\n", string(greenHeart), alert.Annotations["description"])
		default:
			text = "unknown"
		}
		sendAlarm(text, &tgbot)
	}
}

func main() {
	port := ":" + os.Args[2]
	log.Printf("starting HTTP Server on localhost%s.", port)

	http.HandleFunc("/", HandleRequest)

	var err = http.ListenAndServe(port, nil)

	if err != nil {
		log.Panic("Server failed starting. Error: ", err)
	}
}
