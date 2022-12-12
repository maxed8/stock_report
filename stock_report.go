package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/piquette/finance-go/quote"
	"github.com/robfig/cron/v3"
	"github.com/slack-go/slack"
)

/*
The ReadSymbols function reads in an external file and returns a list of the
stock symbols in the SNP500
*/
func ReadSymbols() []string {
	file, err := os.Open("snp500_symbols.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var symbols []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		symbols = append(symbols, scanner.Text())
	}
	return symbols
}

/*
The GetQuotes function takes in a list of stock symbols from the SNP500
It determines which stocks are overvalued and undervalued by finding the
stocks that are more than 5 standard deviations away from their 200 day average
It then returns a map containing the target stocks
*/
func GetQuotes(symbols []string) map[string]map[string]map[string]float64 {
	var undervalued = map[string]map[string]float64{}
	var overvalued = map[string]map[string]float64{}
	for _, ticker := range symbols {
		q, err := quote.Get(ticker)
		if err == nil && q != nil {
			var stock = map[string]float64{
				"Price":   q.RegularMarketOpen,
				"Average": q.TwoHundredDayAverage,
				"SD":      math.Abs(q.TwoHundredDayAverageChange),
			}
			if stock["Price"] < stock["Average"]-5*stock["SD"] {
				undervalued[ticker] = stock
			}

			if stock["Price"] > stock["Average"]+5*stock["SD"] {
				overvalued[ticker] = stock
			}
		}
	}
	var stocks = map[string]map[string]map[string]float64{}
	stocks["undervalued"] = undervalued
	stocks["overvalued"] = overvalued
	return stocks
}

/*
The FormatOutput function takes in a map of overvalued and undervalued stocks
and formats it for the webhook in a user-readable way
This output string is returned
*/
func FormatOutput(stocks map[string]map[string]map[string]float64) string {
	outputString := "Good morning! \n"
	for value, stocks := range stocks {
		outputString += "\nHere are the stocks that are " + value + "(5 or more SDs outside their 200 day average): \n"
		for ticker, stockInfo := range stocks {
			outputString += "\n---------------------\n"
			outputString += ticker
			for k, v := range stockInfo {
				outputString += ("\n" + k + ": $" + strconv.FormatFloat(v, 'f', 2, 64))
			}
			outputString += "\n---------------------"
		}
	}

	return outputString
}

// The SendWebhook function takes in the webhook url and the message and sends the message to a slack webhook
func SendWebhook(url string, outputString string) {
	var message slack.WebhookMessage
	message.Text = outputString
	slack.PostWebhook(url, &message)
}

/*
The RunAnalysis function calls all the helper functions to read in the stock symbols,
get quotes for each ticker and determine which ones are overvalued and undervalued,
formats the output, and sends it to the slack webhook
*/
func RunAnalysis(webhookURL string) {

	symbols := ReadSymbols()
	stocks := GetQuotes(symbols)
	outputString := FormatOutput(stocks)
	fmt.Println(outputString)
	SendWebhook(webhookURL, outputString)
}

/*
The main function prompts the user for their unique slack webhook url
It then calls the RunAnalysis function with the webhook url
It calls this function every weekday at 9am using a cron function
*/
func main() {
	fmt.Println("Enter your Slack Webhook URL: ")
	var webhookURL string
	fmt.Scanln(&webhookURL)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	c := cron.New()
	c.AddFunc("* * * * *", func() { RunAnalysis(webhookURL) })
	c.Start()
	wg.Wait()
}
