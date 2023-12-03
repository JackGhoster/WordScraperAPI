package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type ScrapedWord struct {
	Text string
}

type WordsStruct struct {
	Words []string `json:"words"`
}

var words = WordsStruct{}

func main() {
	var router = gin.Default()
	router.GET("/words", getWordsJson)
	router.Run("localhost:8080")
}

// api functions
func getWordsJson(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, getWords())
}

// local functions
func scrape() {
	var words []ScrapedWord
	scraper := colly.NewCollector()
	scraper.OnHTML(".wordlist-item", func(e *colly.HTMLElement) {
		newWord := ScrapedWord{}
		newWord.Text = e.Text

		words = append(words, newWord)
	})

	scraper.Visit("https://www.enchantedlearning.com/wordlist/food.shtml")

	file, err := os.Create("words.csv")
	if err != nil {
		log.Fatalln("Failed to create CSV file", err)

	}
	defer file.Close()

	fileWriter := csv.NewWriter(file)

	for _, currentWord := range words {
		record := []string{
			currentWord.Text,
		}

		fileWriter.Write(record)
	}

	defer fileWriter.Flush()
}

func sourceScrapedWords() []string {
	var words []string
	//open the file
	var file, err1 = os.Open("words.csv")
	if err1 != nil {
		words = append(words, err1.Error())
		return words
	}
	defer file.Close()
	//initialize the csv reader
	var reader = csv.NewReader(file)
	var records, err2 = reader.ReadAll()
	if err2 != nil {
		words = append(words, err2.Error())
		return words
	}
	//populate the array

	for _, row := range records {
		word := row[0]
		words = append(words, word)
	}

	return words
}

func getWords() WordsStruct {
	var listOfWords = sourceScrapedWords()

	for _, wordFromList := range listOfWords {
		words.Words = append(words.Words, wordFromList)
	}

	return words
}
