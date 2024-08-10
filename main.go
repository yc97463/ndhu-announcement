package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/anaskhan96/soup"
	"golang.org/x/net/html"
)

type Author struct {
	Department string `json:"department"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

type Link struct {
	Title      string `json:"title"`
	Timestamp  string `json:"timestamp"`
	Url        string `json:"url"`
	Date       string `json:"date"`
	Department string `json:"department"`
	Author     Author `json:"author"`
}

type Attachment struct {
	FileName string `json:"fileName"`
	FileSize string `json:"fileSize"`
	FileURL  string `json:"fileURL"`
}

type Detail struct {
	Title       string       `json:"title"`
	Timestamp   string       `json:"timestamp"`
	Url         string       `json:"url"`
	Date        string       `json:"date"`
	Department  string       `json:"department"`
	Author      Author       `json:"author"`
	Content     string       `json:"content"`
	Attachments []Attachment `json:"attachments"`
}

func createDir(dir string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating 'dist' directory:", err)
		return
	}
}

func createFile(file string, _ string) {
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

func insertSummary(file string, timestamp string, title string, url string, date string, department string, author Author) {
	// Create a new Link instance with the provided data
	newLink := Link{
		Title:      title,
		Timestamp:  timestamp,
		Url:        url,
		Date:       date,
		Department: department,
		Author:     author,
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

func createArticle(file string, timestamp string, title string, url string, date string, department string, author Author, content string, attachments []Attachment) {
	newDetail := Detail{
		Title:       title,
		Timestamp:   timestamp,
		Url:         url,
		Date:        date,
		Department:  department,
		Author:      author,
		Content:     content,
		Attachments: attachments,
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
	jsonData, err := json.MarshalIndent(newDetail, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Write the JSON data to the file
	err = os.WriteFile(file, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Article created successfully:", title)
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

func catchAuthorFromContent(content string) []string {
	// Find the line containing the author information
	lines := strings.Split(content, "\n")
	var authorLine string
	for _, line := range lines {
		if strings.Contains(line, "來源：") {
			authorLine = line
			break
		}
	}

	if authorLine == "" {
		return []string{"", "", "", ""}
	}

	// Extract the author information
	parts := strings.Split(authorLine, "來源： ")
	if len(parts) < 2 {
		return []string{"", "", "", ""}
	}

	authorInfo := strings.Split(parts[1], " - ")
	if len(authorInfo) < 4 {
		// Pad the slice with empty strings if there are missing elements
		for len(authorInfo) < 4 {
			authorInfo = append(authorInfo, "")
		}
	}

	// Trim spaces from each element
	for i := range authorInfo {
		authorInfo[i] = strings.TrimSpace(authorInfo[i])
	}

	// Extract phone number from the last element
	phoneIndex := strings.Index(authorInfo[3], "電話")
	phone := ""
	if phoneIndex != -1 {
		phone = strings.TrimSpace(authorInfo[3][phoneIndex+len("電話"):])
		authorInfo[3] = strings.TrimSpace(authorInfo[3][:phoneIndex])
	}

	return []string{authorInfo[0], authorInfo[1], authorInfo[2], phone}
}

func extractAttachments(content string) []Attachment {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return nil
	}

	var attachments []Attachment
	doc.Find("pre a").Each(func(i int, s *goquery.Selection) {
		fileName := s.Text()
		fileURL, _ := s.Attr("href")
		fileSize := strings.TrimSpace(s.Parent().Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
			return s.Nodes[0].Type == html.TextNode
		}).Last().Text())
		attachments = append(attachments, Attachment{
			FileName: fileName,
			FileSize: fileSize,
			FileURL:  fileURL,
		})
	})

	return attachments
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
		createDir("dist/" + value)
	}

	// get data from announce.ndhu.edu.tw
	for key, value := range category {
		fmt.Print("\n=====================================\n")
		fmt.Println("Category: ", value)

		for i := 0; i < 5; i++ {
			page := fmt.Sprintf("%d", i+1)
			fmt.Println("Page: ", page)

			resp, err := soup.Get(baseUrl + "mail_page.php?sort=" + key + "&page=" + page)
			if err != nil {
				os.Exit(1)
			}

			doc := soup.HTMLParse(resp)
			table := doc.Find("div", "class", "column1-unit").Find("table").Find("tbody")
			items := table.FindAll("tr")

			fmt.Printf("Found %d items:\n", len(items))

			for _, item := range items {
				title := item.Find("td", "class", "subject").FindAll("a")[0].Text()
				url := item.Find("td", "class", "subject").FindAll("a")[0].Attrs()["href"]
				timestamp := strings.Split(url, "?timestamp=")[1]
				date := item.Find("td", "class", "date").Text()
				department := item.Find("td", "class", "department").Text()

				content := announce_detail(baseUrl, url)

				authorInfo := catchAuthorFromContent(content)
				author := Author{
					Department: authorInfo[0],
					Name:       authorInfo[1],
					Email:      authorInfo[2],
					Phone:      authorInfo[3],
				}

				attachments := extractAttachments(content)

				insertSummary("dist/"+value+"/"+page+".json", timestamp, title, url, date, department, author)
				createArticle("dist/article/"+timestamp+".json", timestamp, title, url, date, department, author, content, attachments)
			}
		}
	}
}
