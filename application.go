package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape(url string, client *http.Client, req *http.Request) bool {
	// Request the HTML page.
	result := false
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".available-tickets").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		countStr := strings.TrimSpace(s.Text())
		count, err := strconv.Atoi(countStr)
		if err != nil {
			fmt.Errorf("%s", err)
		}
		if count > 0 {
			fmt.Printf("Tickets available! Count: %d", count)
			resp, _ := client.Do(req)
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				var data map[string]interface{}
				decoder := json.NewDecoder(resp.Body)
				err := decoder.Decode(&data)
				if err == nil {
					fmt.Println(data["sid"])
				}
			} else {
				fmt.Println(resp.Status)
			}
			result = true
		} else {
			fmt.Println("Tickets not available")
			result = false
		}

	})
	return result
}

func main() {
	accountSid := "AC09b36b662199a7b5c847adb42397621e"
	authToken := "50b5e9f4f024d825a501f5afd7f50b47"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	msg := "Lozynka, tickets available!"

	msgData := url.Values{}
	msgData.Set("To", "+380636341941")
	msgData.Set("From", "+19203102085")
	msgData.Set("Body", msg)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	argsWithoutProg := os.Args[1:]
	url := "https://ticketclub.com.ua/event/4965"
	if len(argsWithoutProg) == 0 {
		fmt.Printf("URL must be provided")
		return
	} else {
		url = argsWithoutProg[0]
	}
	for {

		if ExampleScrape(url, client, req) {
			return
		}
		fmt.Println("Sleeping...")
		time.Sleep(10 * time.Second)
	}

}
