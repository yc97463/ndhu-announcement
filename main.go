package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/anaskhan96/soup"
)

type Link struct {
	Title    string `json:"title"`
	Timestamp	 string `json:"timestamp"`
	Url       string `json:"url"`
	Date       string `json:"date"`
	Department string `json:"department"`
	Author       string `json:"author"`
	Content     string `json:"content"`
}

type Detail struct {
	Title    string `json:"title"`
	Timestamp	 string `json:"timestamp"`
	Url       string `json:"url"`
	Date       string `json:"date"`
	Department string `json:"department"`
	Author       string `json:"author"`
	Content     string `json:"content"`
}

func createDir(dir string){
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating 'dist' directory:", err)
		return
	}
}

func createFile(file string, jsonType string){
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
}

func insertSummary(file string, timestamp string, title string, url string, date string, department string, author string, content string) {
	// Create a new Link instance with the provided data
	newLink := Link{
		Title:    title,
		Timestamp:  timestamp,
		Url:       url,
		Date:       date,
		Department: department,
		Author:       author,
		Content:     content,
	}

	// Read existing JSON data from the file, if any
	
	data, err := os.ReadFile(file)
	if err != nil {
		createFile(file, "summary")
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

	fmt.Println("Link added successfully   | ", title)
}

func createArticle(file string, timestamp string, title string, url string, date string, department string, author string, content string) {
	newDetail := Detail{
		Title: 	title,
		Timestamp:  timestamp,
		Url:       url,
		Date:       date,
		Department: department,
		Author:       author,
		Content:     content,
	}

	data, err := os.ReadFile(file)
	if err != nil {
		createFile(file, "detail")
		data, err = os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
	}

	var detail []Detail

	// Unmarshal the existing JSON data into the links slice
	if err := json.Unmarshal(data, &detail); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// write newDetail to detail
	detail = []Detail{newDetail}

	// Marshal the updated links slice back to JSON
	updatedData, err := json.MarshalIndent(detail, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Write the updated data back to the file
	if err := os.WriteFile(file, updatedData, 0644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Detail added successfully | ", title)
}

func announce_detail(baseUrl string, link string) (result string) {
	resp, err := soup.Get(baseUrl + link)
	if err != nil {
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)
	content := doc.Find("div", "class", "column1-unit").HTML()

	return content
}

func main() {
	// baseUrl
	baseUrl := "https://announce.ndhu.edu.tw/"
	category := map[string]string{
		"0": "latest",
		"1": "administration",
		"2": "events",
		"3": "course",
		"4": "admission",
		"5": "conference",
		"6": "pt-scholarship",
		"7": "carreer",
		"8": "other",
	}

	// create directories
	createDir("dist")
	createDir("dist/article")
	// create dir from category
	for _, value := range category {
		createDir("dist/"+value)
	}

	// get data from announce.ndhu.edu.tw
	for key, value := range category {
		fmt.Print("\n=====================================\n")
		fmt.Println("Category: ", value)

		// loop for 5 times
		for i := 0; i < 5; i++ {
			page := fmt.Sprintf("%d", i+1)

			// print page
			fmt.Println("Page: ", page)

			// get data in category
			resp, err := soup.Get(baseUrl + "mail_page.php?sort=" + key + "&page=" + page)
			if err != nil {
				os.Exit(1)
			}

			// parse data to HTML
			doc := soup.HTMLParse(resp)

			// get tr in table
			table := doc.Find("div", "class", "column1-unit").Find("table").Find("tbody")
			items := table.FindAll("tr")

			// print found data
			fmt.Printf("Found %d items:\n", len(items))



				// get data in tr
			for _, item := range items {

				// fields
				title := item.Find("td", "class", "subject").FindAll("a")[0].Text()
				url := item.Find("td", "class", "subject").FindAll("a")[0].Attrs()["href"]
				timestamp := strings.Split(url, "?timestamp=")[1]
				date := item.Find("td", "class", "date").Text()
				department := item.Find("td", "class", "department").Text()
				author := item.Find("td", "class", "user").Text()
				content := announce_detail(baseUrl, url)

				insertSummary("dist/"+value+"/"+page+".json", timestamp, title, url, date, department, author, content)
				createArticle("dist/article/"+timestamp+".json", timestamp, title, url, date, department, author, content)
			}
		}
	}
}
