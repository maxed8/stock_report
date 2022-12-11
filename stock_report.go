package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/piquette/finance-go/quote"
	"github.com/slack-go/slack"
)

// need to do webhook
// schedule it at certain time on every weekday/day market is open
// then do CLI/CLT with user credentials
// then comment the code
// write the read me
// do the style/coverage tests?
// turn in

// use main package
// command line tool - take in credential, send slack webhook daily to that slack webhook account

// pull stock symbols from SNP500
// https://github.com/datasets/s-and-p-500-companies/tree/master/data
// for each symbol, figure out which check to run
// add it to the list of callouts - both undervalued and overvalued
// then format and send to webhook
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

// if current price is average change outside of two hundred day average
// if current price < average price - 5 * abs(price change)
// 5 SDs away from average
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

func SendWebhook(url string, outputString string) {
	var message slack.WebhookMessage
	message.Text = outputString
	slack.PostWebhook("https://hooks.slack.com/services/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", &message)
}

func main() {
	fmt.Println("Enter your Slack Webhook URL: ")
	var webhookURL string
	fmt.Scanln(&webhookURL)
	symbols := ReadSymbols()
	stocks := GetQuotes(symbols)
	outputString := FormatOutput(stocks)
	fmt.Println(outputString)
	SendWebhook(webhookURL, outputString)
}
