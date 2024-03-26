package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/shopspring/decimal"
)

var BTCPortion = decimal.NewFromFloat(0.7)
var ETHPortion = decimal.NewFromFloat(0.3)

// our json looks like this:
//
//	{
//		 "data"=> {
//			 "currency"=>"USD",
//			 "rates"=> {
//				 "BTC"=>"0.000015409648952",
//				 "ETH"=>"0.0002930428692413"
//			 }
//	}
//
// It contains the rates for many more currencies, but we can ignore them for this problem.
// We also ignore the currency field because we are always requesting USD.
type CoinbaseResponse struct {
	Data CurrencyInfo
}

// Also returns a Currency field that is always USD, which we can ignore
type CurrencyInfo struct {
	Rates CryptoRates
}

// The API returns rates of many currencies, and Golang is nice enough
// to parse out only the data we define in our structs.
// We only care about BTC and ETH for this problem
type CryptoRates struct {
	BTC string
	ETH string
}

// Return type
type CryptoHoldings struct {
	BTCHoldings decimal.Decimal
	ETHHoldings decimal.Decimal
}

// I'll be honest.... I don't like passing in the URL here,
// it doesn't really make sense. But I couldn't figure out
// how to test the code with fake responses from the endpoint
// without passing in the URL generated from the httptest server.
// I'm sure there is a way to do it right, but this is what
// I came up with :)
func getRates(url string) CryptoRates {
	// http request
	resp, err := http.Get(url + "/v2/exchange-rates?currency=USD")
	if err != nil {
		log.Fatalln(err)
	}

	// retrieve body from successful request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// extract and parse data from the JSON byte string
	var responseData CoinbaseResponse
	err = json.Unmarshal([]byte(body), &responseData)
	if err != nil {
		log.Fatalln(err)
	}

	return responseData.Data.Rates
}

func computeCryptoHoldings(spendingMoney decimal.Decimal, rates CryptoRates) CryptoHoldings {
	btcRate, err := decimal.NewFromString(rates.BTC)
	if err != nil {
		log.Fatalln(err)
	}

	ethRate, err := decimal.NewFromString(rates.ETH)
	if err != nil {
		log.Fatalln(err)
	}

	// compute amount of USD to spend on each kind of asset
	btcSpendingMoney := BTCPortion.Mul(spendingMoney)
	ethSpendingMoney := ETHPortion.Mul(spendingMoney)

	// spending money USD * (btc / usd) = btc
	btcPurchases := btcSpendingMoney.Mul(btcRate)
	ethPurchases := ethSpendingMoney.Mul(ethRate)

	return CryptoHoldings{BTCHoldings: btcPurchases, ETHHoldings: ethPurchases}
}

func main() {
	// parse spending money from first argument
	spendingMoneyInput := os.Args[1]
	spendingMoney, err := decimal.NewFromString(spendingMoneyInput)
	if err != nil {
		log.Fatalln(err)
	}

	// http request
	rates := getRates("https://api.coinbase.com")
	if rates.BTC == "" || rates.ETH == "" {
		log.Fatalln("BTC or ETH rates not found")
	}

	totalHoldings := computeCryptoHoldings(spendingMoney, rates)
	resultJson, err := json.Marshal(totalHoldings)

	// TODO: write tests
	fmt.Println(string(resultJson[:]))
}
