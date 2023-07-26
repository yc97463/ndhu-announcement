package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/anaskhan96/soup"
)

type Link struct {
	Subject    string `json:"subject"`
	Link       string `json:"link"`
	Date       string `json:"date"`
	Department string `json:"department"`
	User       string `json:"user"`
	Detail     string `json:"detail"`
}

func createFile(file string){

	if err := os.MkdirAll("dist", 0755); err != nil {
		fmt.Println("Error creating 'dist' directory:", err)
		return
	}
	// Create the file
	if file, err := os.Create(file); err != nil {
		fmt.Println("Error creating file:", err)
		return
	} else {
		// Close the file when done with it
		defer file.Close()
	}

	// Write an empty JSON array to the file
	if err := os.WriteFile(file, []byte("[]"), 0644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	return 
}

func addLinks(file string, subject string, link string, date string, department string, user string, detail string) {
	// Create a new Link instance with the provided data
	newLink := Link{
		Subject:    subject,
		Link:       link,
		Date:       date,
		Department: department,
		User:       user,
		Detail:     detail,
	}

	// Read existing JSON data from the file, if any
	
	data, err := os.ReadFile(file)
	if err != nil {
		createFile(file)
		data, err = os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
	}
	
	// if err != nil {
	// 	createFile(file)
	// }


	var links []Link

	// Unmarshal the existing JSON data into the links slice
	if err := json.Unmarshal(data, &links); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Append the newLink to the links slice
	links = append(links, newLink)

	// Marshal the updated links slice back to JSON
	updatedData, err := json.MarshalIndent(links, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Write the updated data back to the file
	if err := os.WriteFile(file, updatedData, 0644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Link added successfully | ", subject)
}

func announce_detail(host string, link string) (result string) {
	resp, err := soup.Get(host + link)
	if err != nil {
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)

	content := doc.Find("div", "class", "column1-unit").HTML()

	// fmt.Print(content.HTML(), "\n", "\n", "\n", "\n")

	return content
}

func main() {
	host := "https://announce.ndhu.edu.tw/"
	resp, err := soup.Get(host + "mail_page.php?sort=0")
	if err != nil {
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)
	// table := doc.Find("div", "class", "column1-unit").FindAll("table")
	table := doc.Find("div", "class", "column1-unit").Find("table").Find("tbody")
	items := table.FindAll("tr")

	fmt.Printf("Found %d items:\n", len(items))

	for _, item := range items {
		// fmt.Println(link.Text(), "| Link :", link.Attrs()["href"])
		// addLinks(link.Text(), link.Attrs()["href"])

		subject := item.Find("td", "class", "subject").FindAll("a")[0].Text()
		link := item.Find("td", "class", "subject").FindAll("a")[0].Attrs()["href"]
		date := item.Find("td", "class", "date").Text()
		department := item.Find("td", "class", "department").Text()
		user := item.Find("td", "class", "user").Text()
		announce_detail(host, link)
		detail := announce_detail(host, link)

		// fmt.Print(subject, link, date, department, user, detail, "\n")

		addLinks("dist/latest.json", subject, link, date, department, user, detail)
	}
}
