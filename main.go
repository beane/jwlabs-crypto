package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

const BTCPercentage = 70
const ETHPercentage = 30

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

type CryptoHoldings struct {
	BTCHoldings float64
	ETHHoldings float64
}

func getRates() CryptoRates {
	// http request
	resp, err := http.Get("https://api.coinbase.com/v2/exchange-rates?currency=USD")
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

func main() {
	// parse spending money from first argument
	spendingMoneyInput := os.Args[1]
	spendingMoney, err := strconv.ParseFloat(spendingMoneyInput, 64)
	if err != nil {
		log.Fatalln(err)
	}

	// http request
	rates := getRates()

	// parse rates from strings
	if rates.BTC == "" || rates.ETH == "" {
		log.Fatalln("BTC or ETH rates not found")
	}

	// convert rates to int64, multiplying by 10^18 to avoid floating point errors
	btcRateFloat, err := strconv.ParseFloat(rates.BTC, 64)
	if err != nil {
		log.Fatalln(err)
	}
	btcRate := int64(btcRateFloat * math.Pow(10, 18))

	ethRateFloat, err := strconv.ParseFloat(rates.ETH, 64)
	if err != nil {
		log.Fatalln(err)
	}
	ethRate := int64(ethRateFloat * math.Pow(10, 18))

	// compute amount of USD to spend on each kind of asset
	btcSpendingMoney := BTCPercentage * spendingMoney / 100
	ethSpendingMoney := spendingMoney - btcSpendingMoney

	// Round to 2 digits of precision
	btcSpendingMoney = math.Round(btcSpendingMoney*100) / 100
	ethSpendingMoney = math.Round(ethSpendingMoney*100) / 100

	// spending money usd * (btc / usd) = btc
	// divide by 10^18 to avoid floating point errors
	btcPurchases := btcSpendingMoney * float64(btcRate) / math.Pow(10, 18)
	ethPurchases := ethSpendingMoney * float64(ethRate) / math.Pow(10, 18)

	totalHoldings := CryptoHoldings{BTCHoldings: btcPurchases, ETHHoldings: ethPurchases}
	resultJson, err := json.Marshal(totalHoldings)

	// TODO: write tests
	fmt.Println(string(resultJson[:]))
}
