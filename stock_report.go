package main

import (
	"bufio"
	"fmt"
	"math"
	"os"

	"github.com/piquette/finance-go/quote"
)

// use main package
// command line tool - take in credential, send slack webhook daily to that slack webhook account
// use cobra library

// pull stock symbols from SNP500
// https://github.com/datasets/s-and-p-500-companies/tree/master/data
// for each symbol, figure out which check to run
// add it to the list of callouts - both undervalued and overvalued
// then format and send to webhook
// use env variable or config file
func ReadSymbols() {
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
	GetQuotes(symbols)
}

// if current price is average change outside of two hundred day average
// if current price < average price - 5 * abs(price change)
func GetQuotes(symbols []string) {
	var undervalued = map[string]map[string]float64{}
	var overvalued = map[string]map[string]float64{}
	for _, ticker := range symbols {
		q, err := quote.Get(ticker)
		if err == nil && q != nil {
			var stock = map[string]float64{
				"price":         q.RegularMarketOpen,
				"average":       q.TwoHundredDayAverage,
				"averageChange": math.Abs(q.TwoHundredDayAverageChange),
			}
			if stock["price"] < stock["average"]-5*stock["averageChange"] {
				undervalued[ticker] = stock
			}

			if stock["price"] > stock["average"]+5*stock["averageChange"] {
				overvalued[ticker] = stock
			}
		}
	}
	fmt.Println(undervalued)
	fmt.Println(overvalued)
}

func FormatOutput() {
	return
}

func SendToWebhook() {
	return
}

func main() {
	ReadSymbols()
}
