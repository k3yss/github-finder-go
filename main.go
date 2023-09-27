package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
)

type Owner struct {
	Login string
}

type Item struct {
	ID              int
	Name            string
	FullName        string `json:"full_name"`
	Owner           Owner
	Description     string
	CreatedAt       string `json:"created_at"`
	StargazersCount int    `json:"stargazers_count"`
}

type JSONData struct {
	Count int `json:"total_count"`
	Items []Item
}

func printData(data JSONData) {
	log.Printf("Repositories Found: %d", data.Count)
	const format = "%v\t%v\t%v\t%v\t\n"
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Repository", "Stars", "Created at", "Description")
	fmt.Fprintf(tw, format, "----------", "-----", "----------", "----------")

	for _, i := range data.Items {
		desc := i.Description

		if len(desc) > 50 {
			desc = string(desc[:50]) + "..."
		}
		t, err := time.Parse(time.RFC3339, i.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tw, format, i.FullName, i.StargazersCount, t.Year(), desc)

	}

	tw.Flush()
}

func main() {

	fmt.Print("---------------------------------Github Finder in go---------------------------------\n")

	fmt.Print("Enter the language you want to search for\n >>")

	var inputLanguage string

	_, err := fmt.Scan(&inputLanguage)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	apiUrl := "https://api.github.com/search/repositories?q=stars:>=10000+language:" + inputLanguage + "&sort=stars&order=desc"

	resp, err := http.Get(apiUrl)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Unexpected status code", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	data := JSONData{}

	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Fatal(err)
	}
	printData(data)
}
