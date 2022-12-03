package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/piquette/finance-go/quote"
)

func AppleQuote() {
	q, err := quote.Get("AAPL")
	if err != nil {
		panic(err)
	}
	fmt.Println(q.RegularMarketOpen)
	fmt.Println(q.TwoHundredDayAverage)
	fmt.Println(q.TwoHundredDayAverageChange)
	fmt.Println(q.TwoHundredDayAverageChangePercent)
}

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
	counter := 1
	for _, ticker := range symbols {
		fmt.Println(counter)
		GetQuote(ticker)
		counter++
	}
}

func GetQuote(ticker string) {
	q, err := quote.Get(ticker)
	if err != nil {
		return
	} else {
		fmt.Println(q.Symbol)
		fmt.Println(q.RegularMarketOpen)
		fmt.Println(q.TwoHundredDayAverage)
		fmt.Println(q.TwoHundredDayAverageChange)
		fmt.Println(q.TwoHundredDayAverageChangePercent)
		fmt.Println("---------------------")
	}
}

func CalculateValue() {
	return
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
