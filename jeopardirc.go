package main

import (
	"encoding/json"
	"fmt"
	irc "github.com/thoj/go-ircevent" //IDE plugin needs this https://github.com/go-lang-plugin-org/go-lang-idea-plugin/issues/720
	"io/ioutil"
	"net/http"
	"strings"
)

const EVENT_JOINED_SERVER = "001"
const EVENT_JOINED_CHANNEL = "366"
const EVENT_MESSAGE_RECEIVED = "PRIVMSG"

var settings = map[string]string{
	"server":  "irc.freenode.net:6667",
	"channel": "#jeotest",
	"nick":    "JeopardIRC",
	"name":    "I.AM.ERROR.",
}

type question struct {
	Answer   string
	Question string
}

var currentQuestion question
var currentAnswer string

func main() {
	fmt.Printf("Hello world!")

	irccon1 := irc.IRC(settings["nick"], settings["name"])
	//irccon1.VerboseCallbackHandler = true
	irccon1.Debug = true
	err := irccon1.Connect(settings["server"])
	if err != nil {
		fmt.Printf("Can't connect to server.", err.Error())
	}
	irccon1.AddCallback(EVENT_JOINED_SERVER, func(e *irc.Event) {
		fmt.Println("DOO")
		irccon1.Join(settings["channel"])
	})

	irccon1.AddCallback(EVENT_JOINED_CHANNEL, func(e *irc.Event) {
		fmt.Println("Yay")
		irccon1.Privmsg(settings["channel"], "Hi guys! JeopardIRC is in da place! Wanna play? Type start :)")
	})

	irccon1.AddCallback("PRIVMSG", func(e *irc.Event) {
		fmt.Printf("Msg recvd", e.Message())
		switch strings.ToLower(e.Message()) {
		case "start":
			currentQuestion = fetchQuestion()
			currentAnswer = strings.ToLower(currentQuestion.Answer)
			irccon1.Privmsgf(settings["channel"], "QUESTION: %s", currentQuestion.Question)
			fmt.Printf("Answer is %s", currentQuestion.Answer)

		case currentAnswer:
			irccon1.Privmsgf(settings["channel"], "Correct! The answer was %s!", currentQuestion.Answer)
		}
	})
	irccon1.Loop()

	//irccon2.Quit()

}

func fetchQuestion() question {
	resp, err := http.Get("http://jservice.io/api/random")
	if err != nil {
		fmt.Println("Something went wrong")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Body: %s\n", body)
	fmt.Printf("Error: %v\n", err)

	var newQuestions []question
	err = json.Unmarshal(body, &newQuestions)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	//fmt.Printf("%v", newQuestions)
	return newQuestions[0]
}
