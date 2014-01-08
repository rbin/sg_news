package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"log"
	"net/http"
)

type Item struct {
	Title     string
	URL       string
	LinkScore int `json:"score"`
}

type Response struct {
	Data struct {
		Children []struct {
			Data Item
		}
	}
}

func Get(reddit string) ([]Item, error) {

	url := fmt.Sprintf("http://reddit.com/r/%s.json", reddit)
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	resp := new(Response)
	err = json.NewDecoder(r.Body).Decode(resp)
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(resp.Data.Children))
	for i, child := range resp.Data.Children {
		items[i] = child.Data
	}
	return items, nil

}

func (i Item) String() string {

	com := ""
	switch i.LinkScore {
	case 0:
		// nothing
	case 1:
		com = " Score: 1"
	default:
		com = fmt.Sprintf(" (Score: %d)", i.LinkScore)
	}
	return fmt.Sprintf("<p>%s<b>%s</b><br/> <a href=\"%s\">%s</a></p>", i.Title, com, i.URL, i.URL)

}

func Email() string {

	var buffer bytes.Buffer

	items, err := Get("golang")
	if err != nil {
		log.Fatal(err)
	}
	// Need to build strings from items
	for _, item := range items {
		buffer.WriteString(item.String())
	}

	return buffer.String()
}

func main() {

	sg := sendgrid.NewSendGridClient("sendgrid_user", "sendgrid_key")
	message := sendgrid.NewMail()

	message.AddTo("myemail@me.com")
	message.AddToName("Robin Johnson")
	message.AddSubject("Your Daily Golang Breakfast News!")
	message.AddFrom("rbin@sendgrid.com")

	message.AddHTML(Email())

	if rep := sg.Send(message); rep == nil {
		fmt.Println("Email sent!")
		fmt.Println("Closing...")
	} else {
		fmt.Println(rep)
	}

}
